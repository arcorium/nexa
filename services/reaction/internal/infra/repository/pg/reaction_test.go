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
  "github.com/brianvoe/gofakeit/v7"
  "github.com/stretchr/testify/suite"
  "github.com/testcontainers/testcontainers-go"
  "github.com/testcontainers/testcontainers-go/modules/postgres"
  "github.com/testcontainers/testcontainers-go/wait"
  "github.com/uptrace/bun"
  "go.opentelemetry.io/otel/trace"
  "go.opentelemetry.io/otel/trace/noop"
  "nexa/services/reaction/internal/domain/entity"
  "nexa/services/reaction/internal/infra/repository/model"
  "reflect"
  "testing"
  "time"
)

const (
  LIKE_DB_USERNAME = "user"
  LIKE_DB_PASSWORD = "password"
  LIKE_DB          = "nexa"

  SEED_LIKE_DATA_SIZE = 20
)

type likeTestSuite struct {
  suite.Suite
  container *postgres.PostgresContainer
  db        bun.IDB
  tracer    trace.Tracer // Mock

  reactionSeed []entity.Reaction
}

func (f *likeTestSuite) SetupSuite() {
  for i := 0; i < SEED_LIKE_DATA_SIZE; i += 1 {
    f.reactionSeed = append(f.reactionSeed, generateReaction())
  }
  // Index 0 - 5 will have same post id (post with have multiple likes)
  for i := 0; i < 5; i += 1 {
    f.reactionSeed[i].ItemType = entity.ItemPost
    f.reactionSeed[i].ItemId = f.reactionSeed[0].ItemId
  }
  // Index 0 - 2 is Like reaction
  // Index 3 - 4 is Dislike reaction
  f.reactionSeed[0].ReactionType = entity.ReactionLike
  f.reactionSeed[1].ReactionType = entity.ReactionLike
  f.reactionSeed[2].ReactionType = entity.ReactionLike
  f.reactionSeed[3].ReactionType = entity.ReactionDislike
  f.reactionSeed[4].ReactionType = entity.ReactionDislike
  // Index 6 - 10 will have same comment id (post with have multiple likes)
  for i := 5; i < 10; i += 1 {
    f.reactionSeed[i].ItemType = entity.ItemComment
    f.reactionSeed[i].ItemId = f.reactionSeed[5].ItemId
  }
  // Index 5 - 6 is Like reaction
  // Index 7 - 9 is Dislike reaction
  f.reactionSeed[5].ReactionType = entity.ReactionLike
  f.reactionSeed[6].ReactionType = entity.ReactionLike
  f.reactionSeed[7].ReactionType = entity.ReactionDislike
  f.reactionSeed[8].ReactionType = entity.ReactionDislike
  f.reactionSeed[9].ReactionType = entity.ReactionDislike
  // Index 10 - 14 will have same user id
  for i := 10; i < 15; i += 1 {
    f.reactionSeed[i].UserId = f.reactionSeed[10].UserId
  }
  f.reactionSeed[10].ReactionType = entity.ReactionLike
  f.reactionSeed[10].ItemType = entity.ItemPost
  f.reactionSeed[11].ReactionType = entity.ReactionDislike
  f.reactionSeed[11].ItemType = entity.ItemPost

  // Index 15 - 19 will have same post id (post with have multiple likes)
  for i := 15; i < 20; i += 1 {
    f.reactionSeed[i].ItemType = entity.ItemPost
    f.reactionSeed[i].ItemId = f.reactionSeed[15].ItemId
  }
  f.reactionSeed[15].ReactionType = entity.ReactionLike
  f.reactionSeed[16].ReactionType = entity.ReactionLike
  f.reactionSeed[17].ReactionType = entity.ReactionDislike
  f.reactionSeed[18].ReactionType = entity.ReactionDislike
  f.reactionSeed[19].ReactionType = entity.ReactionDislike

  ctx := context.Background()

  container, err := postgres.Run(ctx,
    "docker.io/postgres:16-alpine",
    postgres.WithUsername(LIKE_DB_USERNAME),
    postgres.WithPassword(LIKE_DB_PASSWORD),
    postgres.WithDatabase(LIKE_DB),
    testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").
      WithOccurrence(2).
      WithStartupTimeout(5*time.Second)),
  )
  f.Require().NoError(err)
  f.container = container

  inspect, err := container.Inspect(ctx)
  f.Require().NoError(err)
  ports := inspect.NetworkSettings.Ports
  mapped := ports["5432/tcp"]

  db, err := database.OpenPostgresWithConfig(&sharedConf.PostgresDatabase{
    Address:  fmt.Sprintf("%s:%s", types.Must(container.Host(ctx)), mapped[0].HostPort),
    Username: LIKE_DB_USERNAME,
    Password: LIKE_DB_PASSWORD,
    Name:     LIKE_DB,
    IsSecure: false,
    Timeout:  time.Second * 10,
  }, true)
  f.Require().NoError(err)
  f.db = db

  // Tracer
  provider := noop.NewTracerProvider()
  f.tracer = provider.Tracer("MOCK")

  // Seed
  model.RegisterBunModels(db)
  err = model.CreateTables(db)
  f.Require().NoError(err)

  // Seeding
  tokenCount := 0
  token := util.CastSliceP(f.reactionSeed, func(from *entity.Reaction) model.Reaction {
    tokenCount++
    return model.FromReactionDomain(from, func(token *entity.Reaction, models *model.Reaction) {
      models.CreatedAt = time.Now().Add(time.Duration(tokenCount) * time.Hour).UTC()
    })
  })

  err = database.Seed(f.db, token...)
  f.Require().NoError(err)
}

func (f *likeTestSuite) TearDownSuite() {
  err := f.container.Terminate(context.Background())
  f.Require().NoError(err)
}

func (f *likeTestSuite) Test_reactionRepository_DeleteByItemId() {
  type args struct {
    ctx      context.Context
    itemType entity.ItemType
    itemIds  []types.Id
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    f.Run(tt.name, func() {
      tx, err := f.db.BeginTx(tt.args.ctx, nil)
      f.Require().NoError(err)
      defer tx.Rollback()

      r := &reactionRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      if err := r.DeleteByItemId(tt.args.ctx, tt.args.itemType, tt.args.itemIds...); (err != nil) != tt.wantErr {
        t.Errorf("DeleteByItemId() error = %v, wantErr %v", err, tt.wantErr)
      }
    })
  }
}

func (f *likeTestSuite) Test_reactionRepository_DeleteByUserId() {
  type args struct {
    ctx    context.Context
    userId types.Id
    opt    optional.Object[entity.Item]
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    f.Run(tt.name, func() {
      tx, err := f.db.BeginTx(tt.args.ctx, nil)
      f.Require().NoError(err)
      defer tx.Rollback()

      r := &reactionRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      if err := r.DeleteByUserId(tt.args.ctx, tt.args.userId, tt.args.opt); (err != nil) != tt.wantErr {
        t.Errorf("DeleteByUserId() error = %v, wantErr %v", err, tt.wantErr)
      }
    })
  }
}

func (f *likeTestSuite) Test_reactionRepository_Delsert() {
  type args struct {
    ctx          context.Context
    reaction     *entity.Reaction
    expectInsert bool
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "Should upsert",
      args: args{
        ctx:          context.Background(),
        reaction:     generateReactionP(),
        expectInsert: true,
      },
      wantErr: false,
    },
    {
      name: "Should delete",
      args: args{
        ctx:          context.Background(),
        reaction:     &f.reactionSeed[11],
        expectInsert: false,
      },
      wantErr: false,
    },
  }
  for _, tt := range tests {
    f.Run(tt.name, func() {
      tx, err := f.db.BeginTx(tt.args.ctx, nil)
      f.Require().NoError(err)
      defer tx.Rollback()

      r := &reactionRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      err = r.Delsert(tt.args.ctx, tt.args.reaction)
      if res := err != nil; res {
        if res != tt.wantErr {
          t.Errorf("Delsert() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      param := repo.QueryParameter{
        Offset: 0,
        Limit:  0,
      }
      res, err := r.FindByItemId(tt.args.ctx, tt.args.reaction.ItemType, tt.args.reaction.ItemId, param)
      f.Require().Equal(err == nil, tt.args.expectInsert)
      f.Require().Equal(res.Data != nil, tt.args.expectInsert)
    })
  }
}

func (f *likeTestSuite) Test_reactionRepository_FindByItemId() {
  type args struct {
    ctx       context.Context
    itemType  entity.ItemType
    itemId    types.Id
    parameter repo.QueryParameter
  }
  tests := []struct {
    name    string
    args    args
    want    repo.PaginatedResult[entity.Reaction]
    wantErr bool
  }{
    {
      name: "Find valid item",
      args: args{
        ctx:      context.Background(),
        itemType: f.reactionSeed[0].ItemType,
        itemId:   f.reactionSeed[0].ItemId,
        parameter: repo.QueryParameter{
          Offset: 0,
          Limit:  0,
        },
      },
      want: repo.PaginatedResult[entity.Reaction]{
        Data:    nil,
        Total:   0,
        Element: 0,
      },
      wantErr: false,
    },
  }
  for _, tt := range tests {
    f.Run(tt.name, func() {
      tx, err := f.db.BeginTx(tt.args.ctx, nil)
      f.Require().NoError(err)
      defer tx.Rollback()

      r := &reactionRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      got, err := r.FindByItemId(tt.args.ctx, tt.args.itemType, tt.args.itemId, tt.args.parameter)
      if (err != nil) != tt.wantErr {
        t.Errorf("FindByItemId() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("FindByItemId() got = %v, want %v", got, tt.want)
      }
    })
  }
}

func (f *likeTestSuite) Test_reactionRepository_FindByUserId() {
  type args struct {
    ctx       context.Context
    userId    types.Id
    parameter repo.QueryParameter
  }
  tests := []struct {
    name    string
    args    args
    want    repo.PaginatedResult[entity.Reaction]
    wantErr bool
  }{
    {
      name: "User id has like data",
      args: args{
        ctx:    context.Background(),
        userId: f.reactionSeed[0].UserId,
        parameter: repo.QueryParameter{
          Offset: 0,
          Limit:  0,
        },
      },
      want: repo.PaginatedResult[entity.Reaction]{
        Data:    nil,
        Total:   0,
        Element: 0,
      },
      wantErr: false,
    },
  }
  for _, tt := range tests {
    f.Run(tt.name, func() {
      tx, err := f.db.BeginTx(tt.args.ctx, nil)
      f.Require().NoError(err)
      defer tx.Rollback()

      r := &reactionRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      got, err := r.FindByUserId(tt.args.ctx, tt.args.userId, tt.args.parameter)
      if (err != nil) != tt.wantErr {
        t.Errorf("FindByUserId() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("FindByUserId() got = %v, want %v", got, tt.want)
      }
    })
  }
}

func (f *likeTestSuite) Test_reactionRepository_GetCounts() {
  dummyId := types.MustCreateId()
  dummyId2 := types.MustCreateId()

  type args struct {
    ctx      context.Context
    itemType entity.ItemType
    itemIds  []types.Id
  }
  tests := []struct {
    name    string
    args    args
    want    []entity.Count
    wantErr bool
  }{
    {
      name: "Valid single id",
      args: args{
        ctx:      context.Background(),
        itemType: f.reactionSeed[0].ItemType,
        itemIds:  []types.Id{f.reactionSeed[0].ItemId},
      },
      want: []entity.Count{{
        ItemId:   f.reactionSeed[0].ItemId,
        ItemType: f.reactionSeed[0].ItemType,
        Like:     3,
        Dislike:  2,
      }},
      wantErr: false,
    },
    {
      name: "Invalid single id",
      args: args{
        ctx:      context.Background(),
        itemType: f.reactionSeed[0].ItemType,
        itemIds:  []types.Id{dummyId},
      },
      want: []entity.Count{{
        ItemId:   dummyId,
        ItemType: f.reactionSeed[0].ItemType,
        Like:     0,
        Dislike:  0,
      }},
      wantErr: false,
    },
    {
      name: "Valid multiple ids",
      args: args{
        ctx:      context.Background(),
        itemType: f.reactionSeed[0].ItemType,
        itemIds:  []types.Id{f.reactionSeed[11].ItemId, f.reactionSeed[15].ItemId, f.reactionSeed[10].ItemId, f.reactionSeed[0].ItemId},
      },
      want: []entity.Count{
        {
          ItemId:   f.reactionSeed[11].ItemId,
          ItemType: f.reactionSeed[0].ItemType,
          Like:     0,
          Dislike:  1,
        },
        {
          ItemId:   f.reactionSeed[15].ItemId,
          ItemType: f.reactionSeed[0].ItemType,
          Like:     2,
          Dislike:  3,
        },
        {
          ItemId:   f.reactionSeed[10].ItemId,
          ItemType: f.reactionSeed[0].ItemType,
          Like:     1,
          Dislike:  0,
        },
        {
          ItemId:   f.reactionSeed[0].ItemId,
          ItemType: f.reactionSeed[0].ItemType,
          Like:     3,
          Dislike:  2,
        },
      },
      wantErr: false,
    },
    {
      name: "Combination ids",
      args: args{
        ctx:      context.Background(),
        itemType: f.reactionSeed[0].ItemType,
        itemIds:  []types.Id{f.reactionSeed[0].ItemId, dummyId, f.reactionSeed[10].ItemId, dummyId2},
      },
      want: []entity.Count{
        {
          ItemId:   f.reactionSeed[0].ItemId,
          ItemType: f.reactionSeed[0].ItemType,
          Like:     3,
          Dislike:  2,
        },
        {
          ItemId:   dummyId,
          ItemType: f.reactionSeed[0].ItemType,
          Like:     0,
          Dislike:  0,
        },
        {
          ItemId:   f.reactionSeed[10].ItemId,
          ItemType: f.reactionSeed[0].ItemType,
          Like:     1,
          Dislike:  0,
        },
        {
          ItemId:   dummyId2,
          ItemType: f.reactionSeed[0].ItemType,
          Like:     0,
          Dislike:  0,
        },
      },
      wantErr: false,
    },
  }
  for _, tt := range tests {
    f.Run(tt.name, func() {
      tx, err := f.db.BeginTx(tt.args.ctx, nil)
      f.Require().NoError(err)
      defer tx.Rollback()

      r := &reactionRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      got, err := r.GetCounts(tt.args.ctx, tt.args.itemType, tt.args.itemIds...)
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

func (f *likeTestSuite) Test_reactionRepository_delete() {
  type args struct {
    ctx      context.Context
    reaction *model.Reaction
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "New data",
      args: args{
        ctx: context.Background(),
        reaction: &model.Reaction{
          UserId:    types.MustCreateId().String(),
          Reaction:  entity.ReactionLike.Underlying(),
          ItemType:  entity.ItemPost.Underlying(),
          ItemId:    types.MustCreateId().String(),
          CreatedAt: time.Now().UTC(),
        },
      },
      wantErr: true,
    },
    {
      name: "Duplicated data",
      args: args{
        ctx: context.Background(),
        reaction: &model.Reaction{
          UserId:    f.reactionSeed[0].UserId.String(),
          Reaction:  0,
          ItemType:  0,
          ItemId:    f.reactionSeed[0].ItemId.String(),
          CreatedAt: time.Now().UTC(),
        },
      },
      wantErr: false,
    },
  }
  for _, tt := range tests {
    f.Run(tt.name, func() {
      tx, err := f.db.BeginTx(tt.args.ctx, nil)
      f.Require().NoError(err)
      defer tx.Rollback()

      r := &reactionRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      if err := r.delete(tt.args.ctx, tt.args.reaction); (err != nil) != tt.wantErr {
        t.Errorf("delete() error = %v, wantErr %v", err, tt.wantErr)
      }
    })
  }
}

func (f *likeTestSuite) Test_reactionRepository_insert() {
  type args struct {
    ctx      context.Context
    reaction *model.Reaction
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "Valid data",
      args: args{
        ctx: context.Background(),
        reaction: &model.Reaction{
          UserId:    types.MustCreateId().String(),
          Reaction:  entity.ReactionLike.Underlying(),
          ItemType:  entity.ItemPost.Underlying(),
          ItemId:    types.MustCreateId().String(),
          CreatedAt: time.Now().UTC(),
        },
      },
      wantErr: false,
    },
    {
      name: "Duplicated data",
      args: args{
        ctx: context.Background(),
        reaction: &model.Reaction{
          UserId:    f.reactionSeed[0].UserId.String(),
          Reaction:  0,
          ItemType:  0,
          ItemId:    f.reactionSeed[0].ItemId.String(),
          CreatedAt: time.Now().UTC(),
        },
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {
    f.Run(tt.name, func() {
      tx, err := f.db.BeginTx(tt.args.ctx, nil)
      f.Require().NoError(err)
      defer tx.Rollback()

      r := &reactionRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      if err := r.upsert(tt.args.ctx, tt.args.reaction); (err != nil) != tt.wantErr {
        t.Errorf("upsert() error = %v, wantErr %v", err, tt.wantErr)
      }
    })
  }
}

func TestLike(t *testing.T) {
  suite.Run(t, &likeTestSuite{})
}

func generateReaction() entity.Reaction {
  return entity.Reaction{
    UserId:       types.MustCreateId(),
    ReactionType: entity.ReactionType(gofakeit.UintN(uint(entity.ReactionDislike.Underlying()))),
    ItemType:     entity.ItemType(gofakeit.UintN(uint(entity.ItemComment.Underlying()))),
    ItemId:       types.MustCreateId(),
    CreatedAt:    time.Now().UTC(),
  }
}

func generateReactionP() *entity.Reaction {
  temp := generateReaction()
  return &temp
}

func ignoreReactionField(reaction *entity.Reaction) {
  reaction.CreatedAt = time.Time{}
}

func ignoreReactionsFields(reactions ...entity.Reaction) {
  for i := 0; i < len(reactions); i++ {
    ignoreReactionField(&reactions[i])
  }
}
