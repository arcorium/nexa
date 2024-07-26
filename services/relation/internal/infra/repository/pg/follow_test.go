package pg

import (
  "context"
  "database/sql"
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
  FOLLOW_DB_USERNAME = "user"
  FOLLOW_DB_PASSWORD = "password"
  FOLLOW_DB          = "nexa"

  SEED_FOLLOW_DATA_SIZE = 12
)

type followTestSuite struct {
  suite.Suite
  container *postgres.PostgresContainer
  db        bun.IDB
  tracer    trace.Tracer // Mock

  followSeed []entity.Follow
}

func (b *followTestSuite) SetupSuite() {
  // Index 0 - 4 is same user with different followed user
  // Index 5 - 9 is different follower, but the same followed user
  userId := types.MustCreateId()
  for i := 0; i < 5; i += 1 {
    ent := generateFollow(optional.Some(userId))
    b.followSeed = append(b.followSeed, ent)
  }
  for i := 5; i < 10; i += 1 {
    ent := generateFollow(optional.NullId)
    ent.FolloweeId = userId
    b.followSeed = append(b.followSeed, ent)
  }

  // Index 10 - 11 is mutual follows
  blockerId := types.MustCreateId()
  blockedId := types.MustCreateId()

  ent := generateFollow(optional.NullId)
  ent.FollowerId = blockerId
  ent.FolloweeId = blockedId
  b.followSeed = append(b.followSeed, ent)

  ent = generateFollow(optional.NullId)
  ent.FollowerId = blockedId
  ent.FolloweeId = blockerId
  b.followSeed = append(b.followSeed, ent)

  ctx := context.Background()

  container, err := postgres.Run(ctx,
    "docker.io/postgres:16-alpine",
    postgres.WithUsername(FOLLOW_DB_USERNAME),
    postgres.WithPassword(FOLLOW_DB_PASSWORD),
    postgres.WithDatabase(FOLLOW_DB),
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
    Username: FOLLOW_DB_USERNAME,
    Password: FOLLOW_DB_PASSWORD,
    Name:     FOLLOW_DB,
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
  follows := util.CastSliceP(b.followSeed, func(from *entity.Follow) model.Follow {
    counts++
    return model.FromFollowDomain(from, func(token *entity.Follow, models *model.Follow) {
      models.CreatedAt = time.Now().Add(time.Duration(counts) * time.Hour).UTC()
    })
  })

  err = database.Seed(b.db, follows...)
  b.Require().NoError(err)
}

func (b *followTestSuite) TearDownSuite() {
  err := b.container.Terminate(context.Background())
  b.Require().NoError(err)
}

func (b *followTestSuite) Test_followRepository_Create() {
  type args struct {
    ctx    context.Context
    follow *entity.Follow
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "User not yet followed the user",
      args: args{
        ctx: context.Background(),
        follow: &entity.Follow{
          FollowerId: b.followSeed[0].FollowerId,
          FolloweeId: types.MustCreateId(),
          CreatedAt:  time.Now().UTC(),
        },
      },
      wantErr: false,
    },
    {
      name: "User already followed the user",
      args: args{
        ctx:    context.Background(),
        follow: &b.followSeed[0],
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {
    b.Run(tt.name, func() {
      tx, err := b.db.BeginTx(tt.args.ctx, nil)
      b.Require().NoError(err)
      defer tx.Rollback()

      f := &followRepository{
        db:     tx,
        tracer: b.tracer,
      }
      t := b.T()

      err = f.Create(tt.args.ctx, tt.args.follow)
      if res := (err != nil); res {
        if res != tt.wantErr {
          t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      err = f.Delete(tt.args.ctx, *tt.args.follow)
      b.Require().NoError(err)
    })
  }
}

func (b *followTestSuite) Test_followRepository_Delete() {
  type args struct {
    ctx    context.Context
    follow []entity.Follow
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "User already followed the user",
      args: args{
        ctx:    context.Background(),
        follow: b.followSeed[:1],
      },
      wantErr: false,
    },
    {
      name: "User not yet follows the user",
      args: args{
        ctx: context.Background(),
        follow: []entity.Follow{
          {
            FollowerId: b.followSeed[0].FollowerId,
            FolloweeId: types.MustCreateId(),
            CreatedAt:  time.Now().UTC(),
          },
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

      f := &followRepository{
        db:     tx,
        tracer: b.tracer,
      }
      t := b.T()

      err = f.Delete(tt.args.ctx, tt.args.follow...)
      if res := err != nil; res {
        if res != tt.wantErr {
          t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      err = f.Delete(tt.args.ctx, tt.args.follow...)
      require.Error(t, err)
      require.ErrorIs(t, err, sql.ErrNoRows)
    })
  }
}

func (b *followTestSuite) Test_followRepository_DeleteByUserId() {
  type args struct {
    ctx            context.Context
    deleteFollower bool
    userId         types.Id
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "Valid user without deleting followers",
      args: args{
        ctx:            context.Background(),
        deleteFollower: false,
        userId:         b.followSeed[0].FollowerId,
      },
      wantErr: false,
    },
    {
      name: "Valid user",
      args: args{
        ctx:            context.Background(),
        deleteFollower: true,
        userId:         b.followSeed[0].FollowerId,
      },
      wantErr: false,
    },
    {
      name: "Invalid user",
      args: args{
        ctx:            context.Background(),
        deleteFollower: true,
        userId:         types.MustCreateId(),
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {
    b.Run(tt.name, func() {
      tx, err := b.db.BeginTx(tt.args.ctx, nil)
      b.Require().NoError(err)
      defer tx.Rollback()

      f := &followRepository{
        db:     tx,
        tracer: b.tracer,
      }
      t := b.T()

      err = f.DeleteByUserId(tt.args.ctx, tt.args.deleteFollower, tt.args.userId)
      if res := err != nil; res {
        if res != tt.wantErr {
          t.Errorf("DeleteByUserId() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      counts, err := f.GetCounts(tt.args.ctx, tt.args.userId)
      b.Require().NoError(err)
      b.Require().Len(counts, 1)
      b.Require().Equal(counts[0].TotalFollowings, uint64(0))
      if tt.args.deleteFollower {
        b.Require().Equal(counts[0].TotalFollowers, uint64(0))
      }
    })
  }
}

func (b *followTestSuite) Test_followRepository_Delsert() {
  type args struct {
    ctx    context.Context
    follow *entity.Follow
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

      f := &followRepository{
        db:     tx,
        tracer: b.tracer,
      }
      t := b.T()

      if err := f.Delsert(tt.args.ctx, tt.args.follow); (err != nil) != tt.wantErr {
        t.Errorf("Delsert() error = %v, wantErr %v", err, tt.wantErr)
      }
    })
  }
}

func (b *followTestSuite) Test_followRepository_GetCounts() {
  dummyId := types.MustCreateId()
  dummyId2 := types.MustCreateId()

  type args struct {
    ctx     context.Context
    userIds []types.Id
  }
  tests := []struct {
    name    string
    args    args
    want    []entity.FollowCount
    wantErr bool
  }{
    {
      name: "Single valid user",
      args: args{
        ctx:     context.Background(),
        userIds: []types.Id{b.followSeed[0].FollowerId},
      },
      want: []entity.FollowCount{
        {
          UserId:          b.followSeed[0].FollowerId,
          TotalFollowers:  5,
          TotalFollowings: 5,
        },
      },
      wantErr: false,
    },
    {
      name: "Multiple valid user",
      args: args{
        ctx:     context.Background(),
        userIds: []types.Id{b.followSeed[0].FolloweeId, b.followSeed[1].FolloweeId, b.followSeed[5].FolloweeId},
      },
      want: []entity.FollowCount{
        {
          UserId:          b.followSeed[0].FolloweeId,
          TotalFollowers:  1,
          TotalFollowings: 0,
        },
        {
          UserId:          b.followSeed[1].FolloweeId,
          TotalFollowers:  1,
          TotalFollowings: 0,
        },
        {
          UserId:          b.followSeed[5].FolloweeId,
          TotalFollowers:  5,
          TotalFollowings: 5,
        },
      },
      wantErr: false,
    },
    {
      name: "Multiple combination user",
      args: args{
        ctx:     context.Background(),
        userIds: []types.Id{dummyId, b.followSeed[0].FolloweeId, dummyId2, b.followSeed[1].FolloweeId, b.followSeed[5].FolloweeId},
      },
      want: []entity.FollowCount{
        {
          UserId:          dummyId,
          TotalFollowers:  0,
          TotalFollowings: 0,
        },
        {
          UserId:          b.followSeed[0].FolloweeId,
          TotalFollowers:  1,
          TotalFollowings: 0,
        },
        {
          UserId:          dummyId2,
          TotalFollowers:  0,
          TotalFollowings: 0,
        },
        {
          UserId:          b.followSeed[1].FolloweeId,
          TotalFollowers:  1,
          TotalFollowings: 0,
        },
        {
          UserId:          b.followSeed[5].FolloweeId,
          TotalFollowers:  5,
          TotalFollowings: 5,
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

      f := &followRepository{
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

func (b *followTestSuite) Test_followRepository_GetFollowers() {
  type args struct {
    ctx    context.Context
    userId types.Id
    param  repo.QueryParameter
  }
  tests := []struct {
    name    string
    args    args
    want    repo.PaginatedResult[entity.Follow]
    wantErr bool
  }{
    {
      name: "Get user with multiple followers",
      args: args{
        ctx:    context.Background(),
        userId: b.followSeed[5].FolloweeId,
        param: repo.QueryParameter{
          Offset: 0,
          Limit:  0,
        },
      },
      want: repo.PaginatedResult[entity.Follow]{
        Data: []entity.Follow{
          {
            FollowerId: b.followSeed[9].FollowerId,
            FolloweeId: b.followSeed[5].FolloweeId,
          },
          {
            FollowerId: b.followSeed[8].FollowerId,
            FolloweeId: b.followSeed[5].FolloweeId,
          },
          {
            FollowerId: b.followSeed[7].FollowerId,
            FolloweeId: b.followSeed[5].FolloweeId,
          },
          {
            FollowerId: b.followSeed[6].FollowerId,
            FolloweeId: b.followSeed[5].FolloweeId,
          },
          {
            FollowerId: b.followSeed[5].FollowerId,
            FolloweeId: b.followSeed[5].FolloweeId,
          },
        },
        Total:   5,
        Element: 5,
      },
      wantErr: false,
    },
    {
      name: "Get user with multiple followers with limit and offset",
      args: args{
        ctx:    context.Background(),
        userId: b.followSeed[5].FolloweeId,
        param: repo.QueryParameter{
          Offset: 1,
          Limit:  3,
        },
      },
      want: repo.PaginatedResult[entity.Follow]{
        Data: []entity.Follow{
          {
            FollowerId: b.followSeed[8].FollowerId,
            FolloweeId: b.followSeed[5].FolloweeId,
          },
          {
            FollowerId: b.followSeed[7].FollowerId,
            FolloweeId: b.followSeed[5].FolloweeId,
          },
          {
            FollowerId: b.followSeed[6].FollowerId,
            FolloweeId: b.followSeed[5].FolloweeId,
          },
        },
        Total:   5,
        Element: 3,
      },
      wantErr: false,
    },
    {
      name: "Get user with single followers",
      args: args{
        ctx:    context.Background(),
        userId: b.followSeed[0].FolloweeId,
      },
      want: repo.PaginatedResult[entity.Follow]{
        Data: []entity.Follow{
          {
            FollowerId: b.followSeed[0].FollowerId,
            FolloweeId: b.followSeed[0].FolloweeId,
          },
        },
        Total:   1,
        Element: 1,
      },
      wantErr: false,
    },
    {
      name: "Get user with no following",
      args: args{
        ctx:    context.Background(),
        userId: types.MustCreateId(),
      },
      want: repo.PaginatedResult[entity.Follow]{
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

      f := &followRepository{
        db:     tx,
        tracer: b.tracer,
      }
      t := b.T()

      got, err := f.GetFollowers(tt.args.ctx, tt.args.userId, tt.args.param)
      if (err != nil) != tt.wantErr {
        t.Errorf("GetFollowers() error = %v, wantErr %v", err, tt.wantErr)
        return
      }

      ignoreFollowFields(got.Data...)
      ignoreFollowFields(tt.want.Data...)

      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("GetFollowers() got = %v, want %v", got, tt.want)
      }
    })
  }
}

func (b *followTestSuite) Test_followRepository_GetFollowings() {
  type args struct {
    ctx    context.Context
    userId types.Id
    param  repo.QueryParameter
  }
  tests := []struct {
    name    string
    args    args
    want    repo.PaginatedResult[entity.Follow]
    wantErr bool
  }{
    {
      name: "Get user with multiple following",
      args: args{
        ctx:    context.Background(),
        userId: b.followSeed[0].FollowerId,
        param: repo.QueryParameter{
          Offset: 0,
          Limit:  0,
        },
      },
      want: repo.PaginatedResult[entity.Follow]{
        Data: []entity.Follow{
          {
            FollowerId: b.followSeed[0].FollowerId,
            FolloweeId: b.followSeed[4].FolloweeId,
          },
          {
            FollowerId: b.followSeed[0].FollowerId,
            FolloweeId: b.followSeed[3].FolloweeId,
          },
          {
            FollowerId: b.followSeed[0].FollowerId,
            FolloweeId: b.followSeed[2].FolloweeId,
          },
          {
            FollowerId: b.followSeed[0].FollowerId,
            FolloweeId: b.followSeed[1].FolloweeId,
          },
          {
            FollowerId: b.followSeed[0].FollowerId,
            FolloweeId: b.followSeed[0].FolloweeId,
          },
        },
        Total:   5,
        Element: 5,
      },
      wantErr: false,
    },
    {
      name: "Get user with multiple following with limit and offset",
      args: args{
        ctx:    context.Background(),
        userId: b.followSeed[0].FollowerId,
        param: repo.QueryParameter{
          Offset: 1,
          Limit:  3,
        },
      },
      want: repo.PaginatedResult[entity.Follow]{
        Data: []entity.Follow{
          {
            FollowerId: b.followSeed[0].FollowerId,
            FolloweeId: b.followSeed[3].FolloweeId,
          },
          {
            FollowerId: b.followSeed[0].FollowerId,
            FolloweeId: b.followSeed[2].FolloweeId,
          },
          {
            FollowerId: b.followSeed[0].FollowerId,
            FolloweeId: b.followSeed[1].FolloweeId,
          },
        },
        Total:   5,
        Element: 3,
      },
      wantErr: false,
    },
    {
      name: "Get user with single following",
      args: args{
        ctx:    context.Background(),
        userId: b.followSeed[5].FollowerId,
      },
      want: repo.PaginatedResult[entity.Follow]{
        Data: []entity.Follow{
          {
            FollowerId: b.followSeed[5].FollowerId,
            FolloweeId: b.followSeed[5].FolloweeId,
          },
        },
        Total:   1,
        Element: 1,
      },
      wantErr: false,
    },
    {
      name: "Get user with no following",
      args: args{
        ctx:    context.Background(),
        userId: types.MustCreateId(),
      },
      want: repo.PaginatedResult[entity.Follow]{
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

      f := &followRepository{
        db:     tx,
        tracer: b.tracer,
      }
      t := b.T()

      got, err := f.GetFollowings(tt.args.ctx, tt.args.userId, tt.args.param)
      if (err != nil) != tt.wantErr {
        t.Errorf("GetFollowings() error = %v, wantErr %v", err, tt.wantErr)
        return
      }

      ignoreFollowFields(got.Data...)
      ignoreFollowFields(tt.want.Data...)

      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("GetFollowings() got = %v, want %v", got, tt.want)
      }
    })
  }
}

func (b *followTestSuite) Test_followRepository_IsFollowing() {
  type args struct {
    ctx         context.Context
    userId      types.Id
    followeeIds []types.Id
  }
  tests := []struct {
    name    string
    args    args
    want    []bool
    wantErr bool
  }{
    {
      name: "User followed single users",
      args: args{
        ctx:         context.Background(),
        userId:      b.followSeed[0].FollowerId,
        followeeIds: []types.Id{b.followSeed[0].FolloweeId},
      },
      want:    []bool{true},
      wantErr: false,
    },
    {
      name: "User doesn't follow the users",
      args: args{
        ctx:         context.Background(),
        userId:      types.MustCreateId(),
        followeeIds: []types.Id{b.followSeed[0].FolloweeId},
      },
      want:    []bool{false},
      wantErr: false,
    },
    {
      name: "User followed multiple users",
      args: args{
        ctx:    context.Background(),
        userId: b.followSeed[0].FollowerId,
        followeeIds: []types.Id{
          b.followSeed[0].FolloweeId,
          b.followSeed[1].FolloweeId,
          b.followSeed[2].FolloweeId,
          b.followSeed[3].FolloweeId,
          b.followSeed[4].FolloweeId,
        },
      },
      want: []bool{
        true,
        true,
        true,
        true,
        true,
      },
      wantErr: false,
    },
    {
      name: "User followed combination of multiple users",
      args: args{
        ctx:    context.Background(),
        userId: b.followSeed[0].FollowerId,
        followeeIds: []types.Id{
          b.followSeed[0].FolloweeId,
          types.MustCreateId(),
          b.followSeed[9].FolloweeId,
          b.followSeed[3].FolloweeId,
          types.MustCreateId(),
        },
      },
      want: []bool{
        true,
        false,
        false,
        true,
        false,
      },
      wantErr: false,
    },
  }
  for _, tt := range tests {
    b.Run(tt.name, func() {
      tx, err := b.db.BeginTx(tt.args.ctx, nil)
      b.Require().NoError(err)
      defer tx.Rollback()

      f := &followRepository{
        db:     tx,
        tracer: b.tracer,
      }
      t := b.T()

      got, err := f.IsFollowing(tt.args.ctx, tt.args.userId, tt.args.followeeIds...)
      if (err != nil) != tt.wantErr {
        t.Errorf("IsFollowing() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("IsFollowing() got = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestFollow(t *testing.T) {
  suite.Run(t, &followTestSuite{})
}

func generateFollow(followerId optional.Object[types.Id]) entity.Follow {
  return entity.Follow{
    FollowerId: followerId.ValueOr(types.MustCreateId()),
    FolloweeId: types.MustCreateId(),
    CreatedAt:  time.Now().UTC(),
  }
}

func generateFollowP(followerId optional.Object[types.Id]) *entity.Follow {
  temp := generateFollow(followerId)
  return &temp
}

func ignoreFollowField(follow *entity.Follow) {
  follow.CreatedAt = time.Time{}
}

func ignoreFollowFields(follows ...entity.Follow) {
  for i := 0; i < len(follows); i++ {
    ignoreFollowField(&follows[i])
  }
}
