package pg

import (
  "context"
  "fmt"
  sharedConf "github.com/arcorium/nexa/shared/config"
  "github.com/arcorium/nexa/shared/database"
  "github.com/arcorium/nexa/shared/optional"
  "github.com/arcorium/nexa/shared/types"
  "github.com/arcorium/nexa/shared/util"
  "github.com/arcorium/nexa/shared/util/repo"
  "github.com/stretchr/testify/require"
  "github.com/stretchr/testify/suite"
  "github.com/testcontainers/testcontainers-go"
  "github.com/testcontainers/testcontainers-go/modules/postgres"
  "github.com/testcontainers/testcontainers-go/wait"
  "github.com/uptrace/bun"
  "go.opentelemetry.io/otel/trace"
  "go.opentelemetry.io/otel/trace/noop"
  "nexa/services/relation/internal/domain/entity"
  "nexa/services/relation/internal/infra/repository/model"
  "reflect"
  "testing"
  "time"
)

const (
  BLOCK_DB_USERNAME = "user"
  BLOCK_DB_PASSWORD = "password"
  BLOCK_DB          = "nexa"

  SEED_BLOCK_DATA_SIZE = 12
)

type blockTestSuite struct {
  suite.Suite
  container *postgres.PostgresContainer
  db        bun.IDB
  tracer    trace.Tracer // Mock

  blockSeed []entity.Block
}

func (b *blockTestSuite) SetupSuite() {
  // Index 0 - 4 is same user with different blocked user
  // Index 5 - 9 is different blocker, but the same blocked user
  userId := types.MustCreateId()
  for i := 0; i < 5; i += 1 {
    ent := generateBlock(optional.Some(userId))
    b.blockSeed = append(b.blockSeed, ent)
  }
  for i := 5; i < 10; i += 1 {
    ent := generateBlock(optional.NullId)
    ent.BlockedId = userId
    b.blockSeed = append(b.blockSeed, ent)
  }

  // Index 10 - 11 is mutual block
  blockerId := types.MustCreateId()
  blockedId := types.MustCreateId()

  ent := generateBlock(optional.NullId)
  ent.BlockerId = blockerId
  ent.BlockedId = blockedId
  b.blockSeed = append(b.blockSeed, ent)

  ent = generateBlock(optional.NullId)
  ent.BlockerId = blockedId
  ent.BlockedId = blockerId
  b.blockSeed = append(b.blockSeed, ent)

  ctx := context.Background()

  container, err := postgres.Run(ctx,
    "docker.io/postgres:16-alpine",
    postgres.WithUsername(BLOCK_DB_USERNAME),
    postgres.WithPassword(BLOCK_DB_PASSWORD),
    postgres.WithDatabase(BLOCK_DB),
    testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").
      WithOccurrence(2).
      WithStartupTimeout(5*time.Second)),
  )
  b.Require().NoError(err)
  b.container = container

  inspect, err := container.Inspect(ctx)
  b.Require().NoError(err)
  ports := inspect.NetworkSettings.Ports
  mapped := ports["5432/tcp"]

  db, err := database.OpenPostgresWithConfig(&sharedConf.PostgresDatabase{
    Address:  fmt.Sprintf("%s:%s", types.Must(container.Host(ctx)), mapped[0].HostPort),
    Username: BLOCK_DB_USERNAME,
    Password: BLOCK_DB_PASSWORD,
    Name:     BLOCK_DB,
    IsSecure: false,
    Timeout:  time.Second * 10,
  }, true)
  b.Require().NoError(err)
  b.db = db

  // Tracer
  provider := noop.NewTracerProvider()
  b.tracer = provider.Tracer("MOCK")

  // Seed
  model.RegisterBunModels(db)
  err = model.CreateTables(db)
  b.Require().NoError(err)

  // Seeding
  counts := 0
  blocks := util.CastSliceP(b.blockSeed, func(from *entity.Block) model.Block {
    counts++
    return model.FromBlockDomain(from, func(token *entity.Block, models *model.Block) {
      models.CreatedAt = time.Now().Add(time.Duration(counts) * time.Hour).UTC()
    })
  })

  err = database.Seed(b.db, blocks...)
  b.Require().NoError(err)
}

func (b *blockTestSuite) TearDownSuite() {
  err := b.container.Terminate(context.Background())
  b.Require().NoError(err)
}

func (b *blockTestSuite) Test_blockRepository_Create() {
  type args struct {
    ctx   context.Context
    block *entity.Block
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "User not yet block the user",
      args: args{
        ctx: context.Background(),
        block: &entity.Block{
          BlockerId: b.blockSeed[0].BlockerId,
          BlockedId: types.MustCreateId(),
          CreatedAt: time.Now().UTC(),
        },
      },
      wantErr: false,
    },
    {
      name: "User already block the user",
      args: args{
        ctx:   context.Background(),
        block: &b.blockSeed[0],
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {
    b.Run(tt.name, func() {
      tx, err := b.db.BeginTx(tt.args.ctx, nil)
      b.Require().NoError(err)
      defer tx.Rollback()

      f := &blockRepository{
        db:     tx,
        tracer: b.tracer,
      }
      t := b.T()

      err = f.Create(tt.args.ctx, tt.args.block)
      if res := err != nil; res {
        if res != tt.wantErr {
          t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      blocked, err := f.GetCounts(tt.args.ctx, tt.args.block.BlockerId)
      require.NoError(t, err)
      require.Len(t, blocked, 1)
      require.GreaterOrEqual(t, blocked[0].TotalBlocked, uint64(1))

    })
  }
}

func (b *blockTestSuite) Test_blockRepository_Delete() {
  type fields struct {
    db     bun.IDB
    tracer trace.Tracer
  }
  type args struct {
    ctx   context.Context
    block *entity.Block
  }
  tests := []struct {
    name    string
    fields  fields
    args    args
    wantErr bool
  }{
    {
      name: "User already block the user",
      args: args{
        ctx:   context.Background(),
        block: &b.blockSeed[0],
      },
      wantErr: false,
    },
    {
      name: "User not yet block the user",
      args: args{
        ctx: context.Background(),
        block: &entity.Block{
          BlockerId: b.blockSeed[0].BlockerId,
          BlockedId: types.MustCreateId(),
          CreatedAt: time.Now().UTC(),
        },
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {
    b.Run(tt.name, func() {
      tx, err := b.db.BeginTx(tt.args.ctx, nil)
      b.Require().NoError(err)
      defer tx.Rollback()

      f := &blockRepository{
        db:     tx,
        tracer: b.tracer,
      }
      t := b.T()

      if err := f.Delete(tt.args.ctx, tt.args.block); (err != nil) != tt.wantErr {
        t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
      }
    })
  }
}

func (b *blockTestSuite) Test_blockRepository_DeleteByUserId() {
  type args struct {
    ctx           context.Context
    deleteBlocker bool
    userId        types.Id
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "User has blocked user list",
      args: args{
        ctx:           context.Background(),
        deleteBlocker: true,
        userId:        b.blockSeed[0].BlockedId,
      },
      wantErr: false,
    },
    {
      name: "User has no blocked user list",
      args: args{
        ctx:           context.Background(),
        deleteBlocker: true,
        userId:        types.MustCreateId(),
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {
    b.Run(tt.name, func() {
      tx, err := b.db.BeginTx(tt.args.ctx, nil)
      b.Require().NoError(err)
      defer tx.Rollback()

      f := &blockRepository{
        db:     tx,
        tracer: b.tracer,
      }
      t := b.T()

      err = f.DeleteByUserId(tt.args.ctx, tt.args.deleteBlocker, tt.args.userId)
      if res := err != nil; res {
        if res != tt.wantErr {
          t.Errorf("DeleteByUserId() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      counts, err := f.GetCounts(tt.args.ctx, tt.args.userId)
      require.NoError(t, err)
      require.Len(t, counts, 1)
      require.Equal(t, counts[0].TotalBlocked, uint64(0))
    })
  }
}

func (b *blockTestSuite) Test_blockRepository_Delsert() {
  type args struct {
    ctx   context.Context
    block *entity.Block
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    b.Run(tt.name, func() {
      tx, err := b.db.BeginTx(tt.args.ctx, nil)
      b.Require().NoError(err)
      defer tx.Rollback()

      f := &blockRepository{
        db:     tx,
        tracer: b.tracer,
      }
      t := b.T()

      if err := f.Delsert(tt.args.ctx, tt.args.block); (err != nil) != tt.wantErr {
        t.Errorf("Delsert() error = %v, wantErr %v", err, tt.wantErr)
      }
    })
  }
}

func (b *blockTestSuite) Test_blockRepository_GetBlocked() {
  type args struct {
    ctx    context.Context
    userId types.Id
    param  repo.QueryParameter
  }
  tests := []struct {
    name    string
    args    args
    want    repo.PaginatedResult[entity.Block]
    wantErr bool
  }{
    {
      name: "User has single blocked user list",
      args: args{
        ctx:    context.Background(),
        userId: b.blockSeed[5].BlockerId,
        param: repo.QueryParameter{
          Offset: 0,
          Limit:  0,
        },
      },
      want: repo.PaginatedResult[entity.Block]{
        Data: []entity.Block{
          {
            BlockerId: b.blockSeed[5].BlockerId,
            BlockedId: b.blockSeed[6].BlockedId,
            CreatedAt: b.blockSeed[5].CreatedAt,
          },
        },
        Total:   1,
        Element: 1,
      },
      wantErr: false,
    },
    {
      name: "User has multiple blocked user list",
      args: args{
        ctx:    context.Background(),
        userId: b.blockSeed[0].BlockerId,
      },
      want: repo.PaginatedResult[entity.Block]{
        Data: []entity.Block{
          {
            BlockerId: b.blockSeed[0].BlockerId,
            BlockedId: b.blockSeed[4].BlockedId,
          },
          {
            BlockerId: b.blockSeed[0].BlockerId,
            BlockedId: b.blockSeed[3].BlockedId,
          },
          {
            BlockerId: b.blockSeed[0].BlockerId,
            BlockedId: b.blockSeed[2].BlockedId,
          },
          {
            BlockerId: b.blockSeed[0].BlockerId,
            BlockedId: b.blockSeed[1].BlockedId,
          },
          {
            BlockerId: b.blockSeed[0].BlockerId,
            BlockedId: b.blockSeed[0].BlockedId,
          },
        },
        Total:   5,
        Element: 5,
      },
      wantErr: false,
    },
    {
      name: "User has multiple blocked user list with limit and offset",
      args: args{
        ctx:    context.Background(),
        userId: b.blockSeed[0].BlockerId,
        param: repo.QueryParameter{
          Offset: 1,
          Limit:  3,
        },
      },
      want: repo.PaginatedResult[entity.Block]{
        Data: []entity.Block{
          {
            BlockerId: b.blockSeed[0].BlockerId,
            BlockedId: b.blockSeed[3].BlockedId,
          },
          {
            BlockerId: b.blockSeed[0].BlockerId,
            BlockedId: b.blockSeed[2].BlockedId,
          },
          {
            BlockerId: b.blockSeed[0].BlockerId,
            BlockedId: b.blockSeed[1].BlockedId,
          },
        },
        Total:   5,
        Element: 3,
      },
      wantErr: false,
    },
    {
      name: "User has no blocked user list",
      args: args{
        ctx:    context.Background(),
        userId: types.MustCreateId(),
      },
      want: repo.PaginatedResult[entity.Block]{
        Data:    nil,
        Total:   0,
        Element: 0,
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {
    b.Run(tt.name, func() {
      tx, err := b.db.BeginTx(tt.args.ctx, nil)
      b.Require().NoError(err)
      defer tx.Rollback()

      f := &blockRepository{
        db:     tx,
        tracer: b.tracer,
      }
      t := b.T()

      got, err := f.GetBlocked(tt.args.ctx, tt.args.userId, tt.args.param)
      if (err != nil) != tt.wantErr {
        t.Errorf("GetBlocked() error = %v, wantErr %v", err, tt.wantErr)
        return
      }

      ignoreBlockFields(got.Data...)
      ignoreBlockFields(tt.want.Data...)

      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("GetBlocked() got = %v, want %v", got, tt.want)
      }
    })
  }
}

func (b *blockTestSuite) Test_blockRepository_GetCounts() {
  dummyId := types.MustCreateId()
  dummyId2 := types.MustCreateId()
  dummyId3 := types.MustCreateId()

  type args struct {
    ctx     context.Context
    userIds []types.Id
  }
  tests := []struct {
    name    string
    args    args
    want    []entity.BlockCount
    wantErr bool
  }{
    {
      name: "Single user block some users",
      args: args{
        ctx:     context.Background(),
        userIds: []types.Id{b.blockSeed[0].BlockerId},
      },
      want: []entity.BlockCount{
        {
          UserId:       b.blockSeed[0].BlockerId,
          TotalBlocked: 5,
        },
      },
      wantErr: false,
    },
    {
      name: "Single user without block list",
      args: args{
        ctx:     context.Background(),
        userIds: []types.Id{dummyId},
      },
      want: []entity.BlockCount{
        {
          UserId:       dummyId,
          TotalBlocked: 0,
        },
      },
      wantErr: false,
    },
    {
      name: "Multiple user block some users",
      args: args{
        ctx:     context.Background(),
        userIds: []types.Id{b.blockSeed[0].BlockerId, b.blockSeed[5].BlockerId, b.blockSeed[6].BlockerId},
      },
      want: []entity.BlockCount{
        {
          UserId:       b.blockSeed[0].BlockerId,
          TotalBlocked: 5,
        },
        {
          UserId:       b.blockSeed[5].BlockerId,
          TotalBlocked: 1,
        },
        {
          UserId:       b.blockSeed[6].BlockerId,
          TotalBlocked: 1,
        },
      },
      wantErr: false,
    },
    {
      name: "Multiple combination user block some users",
      args: args{
        ctx:     context.Background(),
        userIds: []types.Id{dummyId, b.blockSeed[0].BlockerId, dummyId2, b.blockSeed[5].BlockerId, b.blockSeed[6].BlockerId, dummyId3},
      },
      want: []entity.BlockCount{
        {
          UserId:       dummyId,
          TotalBlocked: 0,
        },
        {
          UserId:       b.blockSeed[0].BlockerId,
          TotalBlocked: 5,
        },
        {
          UserId:       dummyId2,
          TotalBlocked: 0,
        },
        {
          UserId:       b.blockSeed[5].BlockerId,
          TotalBlocked: 1,
        },
        {
          UserId:       b.blockSeed[6].BlockerId,
          TotalBlocked: 1,
        },
        {
          UserId:       dummyId3,
          TotalBlocked: 0,
        },
      },
      wantErr: false,
    },
  }
  for _, tt := range tests {
    b.Run(tt.name, func() {
      tx, err := b.db.BeginTx(tt.args.ctx, nil)
      b.Require().NoError(err)
      defer tx.Rollback()

      f := &blockRepository{
        db:     tx,
        tracer: b.tracer,
      }
      t := b.T()

      got, err := f.GetCounts(tt.args.ctx, tt.args.userIds...)
      if (err != nil) != tt.wantErr {
        t.Errorf("GetCounts() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("GetCounts() got = %v, want %v", got, tt.want)
      }
    })
  }
}

func (b *blockTestSuite) Test_blockRepository_IsBlocked() {
  type args struct {
    ctx       context.Context
    blockerId types.Id
    targetId  types.Id
  }
  tests := []struct {
    name    string
    args    args
    want    bool
    wantErr bool
  }{
    {
      name: "User is blocked",
      args: args{
        ctx:       context.Background(),
        blockerId: b.blockSeed[0].BlockerId,
        targetId:  b.blockSeed[0].BlockedId,
      },
      want:    true,
      wantErr: false,
    },
    {
      name: "User is not blocked",
      args: args{
        ctx:       context.Background(),
        blockerId: b.blockSeed[0].BlockerId,
        targetId:  types.MustCreateId(),
      },
      want:    false,
      wantErr: false,
    },
  }
  for _, tt := range tests {
    b.Run(tt.name, func() {
      tx, err := b.db.BeginTx(tt.args.ctx, nil)
      b.Require().NoError(err)
      defer tx.Rollback()

      f := &blockRepository{
        db:     tx,
        tracer: b.tracer,
      }
      t := b.T()

      got, err := f.IsBlocked(tt.args.ctx, tt.args.blockerId, tt.args.targetId)
      if (err != nil) != tt.wantErr {
        t.Errorf("IsBlocked() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if got != tt.want {
        t.Errorf("IsBlocked() got = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestBlock(t *testing.T) {
  suite.Run(t, &blockTestSuite{})
}

func generateBlock(blockerId optional.Object[types.Id]) entity.Block {
  return entity.Block{
    BlockerId: blockerId.ValueOr(types.MustCreateId()),
    BlockedId: types.MustCreateId(),
    CreatedAt: time.Now().UTC(),
  }
}

func generateBlockP(blockerId optional.Object[types.Id]) *entity.Block {
  temp := generateBlock(blockerId)
  return &temp
}

func ignoreBlockField(comment *entity.Block) {
  comment.CreatedAt = time.Time{}
}

func ignoreBlockFields(comments ...entity.Block) {
  for i := 0; i < len(comments); i++ {
    ignoreBlockField(&comments[i])
  }
}
