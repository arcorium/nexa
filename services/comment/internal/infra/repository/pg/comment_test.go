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
  "nexa/services/comment/internal/domain/entity"
  "nexa/services/comment/internal/infra/repository/model"
  "reflect"
  "testing"
  "time"
)

const (
  COMMENT_DB_USERNAME = "user"
  COMMENT_DB_PASSWORD = "password"
  COMMENT_DB          = "nexa"

  SEED_COMMENT_DATA_SIZE = 20
)

type commentTestSuite struct {
  suite.Suite
  container *postgres.PostgresContainer
  db        bun.IDB
  tracer    trace.Tracer // Mock

  commentSeed []entity.Comment
}

func (f *commentTestSuite) SetupSuite() {
  // Index 0 - 9 will have same post id
  postId := types.MustCreateId()
  for i := 0; i < 10; i += 1 {
    f.commentSeed = append(f.commentSeed, generateComment(optional.Some(postId)))

  }
  for i := 10; i < 15; i += 1 {
    f.commentSeed = append(f.commentSeed, generateComment(optional.Some(postId)))
  }
  for i := 15; i < SEED_COMMENT_DATA_SIZE; i += 1 {
    f.commentSeed = append(f.commentSeed, generateComment(optional.NullId))
  }
  // Index 1-4 will have parent index 0
  f.commentSeed[1].Parent = &entity.Comment{Id: f.commentSeed[0].Id}
  f.commentSeed[2].Parent = &entity.Comment{Id: f.commentSeed[0].Id}
  f.commentSeed[3].Parent = &entity.Comment{Id: f.commentSeed[0].Id}
  f.commentSeed[4].Parent = &entity.Comment{Id: f.commentSeed[0].Id}
  // Index 6 will have parent index 5
  f.commentSeed[6].Parent = &entity.Comment{Id: f.commentSeed[5].Id}
  // Index 7 will have parent index 6
  f.commentSeed[7].Parent = &entity.Comment{Id: f.commentSeed[6].Id}
  // Index 8 will have parent index 7
  f.commentSeed[8].Parent = &entity.Comment{Id: f.commentSeed[7].Id}
  // Index 9 will have parent index 5
  f.commentSeed[9].Parent = &entity.Comment{Id: f.commentSeed[5].Id}

  ctx := context.Background()

  container, err := postgres.Run(ctx,
    "docker.io/postgres:16-alpine",
    postgres.WithUsername(COMMENT_DB_USERNAME),
    postgres.WithPassword(COMMENT_DB_PASSWORD),
    postgres.WithDatabase(COMMENT_DB),
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
    Username: COMMENT_DB_USERNAME,
    Password: COMMENT_DB_PASSWORD,
    Name:     COMMENT_DB,
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
  commentCount := 0
  comments := util.CastSliceP(f.commentSeed, func(from *entity.Comment) model.Comment {
    commentCount++
    return model.FromCommentDomain(from, func(token *entity.Comment, models *model.Comment) {
      models.CreatedAt = time.Now().Add(time.Duration(commentCount) * time.Hour).UTC()
    })
  })

  err = database.Seed(f.db, comments...)
  f.Require().NoError(err)
}

func (f *commentTestSuite) TearDownSuite() {
  err := f.container.Terminate(context.Background())
  f.Require().NoError(err)
}

func (f *commentTestSuite) Test_commentRepository_Create() {
  type args struct {
    ctx     context.Context
    comment *entity.Comment
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

      c := &commentRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      if err := c.Create(tt.args.ctx, tt.args.comment); (err != nil) != tt.wantErr {
        t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
      }
    })
  }
}

func (f *commentTestSuite) Test_commentRepository_DeleteByIds() {
  type args struct {
    ctx        context.Context
    commentIds []types.Id
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

      c := &commentRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      if err := c.DeleteByIds(tt.args.ctx, tt.args.commentIds...); (err != nil) != tt.wantErr {
        t.Errorf("DeleteByIds() error = %v, wantErr %v", err, tt.wantErr)
      }
    })
  }
}

func (f *commentTestSuite) Test_commentRepository_DeleteByPostIds() {
  type args struct {
    ctx     context.Context
    postIds []types.Id
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "Valid post id",
      args: args{
        ctx:     context.Background(),
        postIds: []types.Id{f.commentSeed[0].PostId},
      },
      wantErr: false,
    },
  }
  for _, tt := range tests {
    f.Run(tt.name, func() {
      tx, err := f.db.BeginTx(tt.args.ctx, nil)
      f.Require().NoError(err)
      defer tx.Rollback()

      c := &commentRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      if _, err := c.DeleteByPostIds(tt.args.ctx, tt.args.postIds...); (err != nil) != tt.wantErr {
        t.Errorf("DeleteByPostIds() error = %v, wantErr %v", err, tt.wantErr)
      }
    })
  }
}

func (f *commentTestSuite) Test_commentRepository_DeleteUsers() {
  type args struct {
    ctx        context.Context
    userId     types.Id
    commentIds []types.Id
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

      c := &commentRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      if _, err := c.DeleteUsers(tt.args.ctx, tt.args.userId, tt.args.commentIds...); (err != nil) != tt.wantErr {
        t.Errorf("DeleteUsers() error = %v, wantErr %v", err, tt.wantErr)
      }
    })
  }
}

func (f *commentTestSuite) Test_commentRepository_FindByPostId() {
  type args struct {
    ctx       context.Context
    postId    types.Id
    showReply bool
    parameter repo.QueryParameter
  }
  tests := []struct {
    name    string
    args    args
    want    repo.PaginatedResult[entity.Comment]
    wantErr bool
  }{
    {
      name: "Valid post",
      args: args{
        ctx:       context.Background(),
        showReply: false,
        postId:    f.commentSeed[0].PostId,
        parameter: repo.QueryParameter{
          Offset: 0,
          Limit:  0,
        },
      },
      want: repo.PaginatedResult[entity.Comment]{
        Data:    f.commentSeed[:10],
        Total:   7,
        Element: 7,
      },
      wantErr: false,
    },
    {
      name: "Valid post with nested comments",
      args: args{
        ctx:       context.Background(),
        showReply: true,
        postId:    f.commentSeed[0].PostId,
        parameter: repo.QueryParameter{
          Offset: 0,
          Limit:  0,
        },
      },
      want: repo.PaginatedResult[entity.Comment]{
        Data:    f.commentSeed[:10],
        Total:   7,
        Element: 7,
      },
      wantErr: false,
    },
    {
      name: "Valid post with nested comments with offset and limit",
      args: args{
        ctx:       context.Background(),
        showReply: true,
        postId:    f.commentSeed[0].PostId,
        parameter: repo.QueryParameter{
          Offset: 3,
          Limit:  2,
        },
      },
      want: repo.PaginatedResult[entity.Comment]{
        Data:    f.commentSeed[:10],
        Total:   10,
        Element: 2,
      },
      wantErr: false,
    },
  }
  for _, tt := range tests {
    f.Run(tt.name, func() {
      tx, err := f.db.BeginTx(tt.args.ctx, nil)
      f.Require().NoError(err)
      defer tx.Rollback()

      c := &commentRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      got, err := c.FindByPostId(tt.args.ctx, tt.args.showReply, tt.args.postId, tt.args.parameter)
      if (err != nil) != tt.wantErr {
        t.Errorf("FindByPostId() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("FindByPostId() got = %v, want %v", got, tt.want)
      }
    })
  }
}

func (f *commentTestSuite) Test_commentRepository_GetReplyCounts() {
  type args struct {
    ctx        context.Context
    commentIds []types.Id
  }
  tests := []struct {
    name    string
    args    args
    want    []entity.Count
    wantErr bool
  }{
    {
      name: "Invalid single post id",
      args: args{
        ctx:        context.Background(),
        commentIds: []types.Id{types.MustCreateId()},
      },
      want: []entity.Count{
        {
          TotalComments: 0,
        },
      },
      wantErr: false,
    },
    {
      name: "Valid single post id",
      args: args{
        ctx:        context.Background(),
        commentIds: []types.Id{f.commentSeed[0].Id},
      },
      want: []entity.Count{
        {
          TotalComments: 4,
        },
      },
      wantErr: false,
    },
    {
      name: "Valid multiple post ids",
      args: args{
        ctx:        context.Background(),
        commentIds: []types.Id{f.commentSeed[0].Id, f.commentSeed[5].Id, f.commentSeed[6].Id, f.commentSeed[5].Id},
      },
      want: []entity.Count{
        {
          TotalComments: 4,
        },
        {
          TotalComments: 2,
        },
        {
          TotalComments: 1,
        },
        {
          TotalComments: 2,
        },
      },
      wantErr: false,
    },
    {
      name: "Combination multiple post ids",
      args: args{
        ctx:        context.Background(),
        commentIds: []types.Id{types.MustCreateId(), f.commentSeed[0].Id, f.commentSeed[11].Id, f.commentSeed[5].Id},
      },
      want: []entity.Count{
        {
          TotalComments: 0,
        },
        {
          TotalComments: 4,
        },
        {
          TotalComments: 0,
        },
        {
          TotalComments: 2,
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

      c := &commentRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      got, err := c.GetReplyCounts(tt.args.ctx, tt.args.commentIds...)
      if (err != nil) != tt.wantErr {
        t.Errorf("GetReplyCounts() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("GetReplyCounts() got = %v, want %v", got, tt.want)
      }
    })
  }
}

func (f *commentTestSuite) Test_commentRepository_GetPostCounts() {
  type args struct {
    ctx     context.Context
    postIds []types.Id
  }
  tests := []struct {
    name    string
    args    args
    want    []entity.Count
    wantErr bool
  }{
    {
      name: "Invalid single post id",
      args: args{
        ctx:     context.Background(),
        postIds: []types.Id{types.MustCreateId()},
      },
      want: []entity.Count{
        {
          TotalComments: 0,
        },
      },
      wantErr: false,
    },
    {
      name: "Valid single post id",
      args: args{
        ctx:     context.Background(),
        postIds: []types.Id{f.commentSeed[0].PostId},
      },
      want: []entity.Count{
        {
          TotalComments: 10,
        },
      },
      wantErr: false,
    },
    {
      name: "Valid multiple post ids",
      args: args{
        ctx:     context.Background(),
        postIds: []types.Id{f.commentSeed[10].PostId, f.commentSeed[0].PostId, f.commentSeed[11].PostId, f.commentSeed[5].PostId},
      },
      want: []entity.Count{
        {
          TotalComments: 1,
        },
        {
          TotalComments: 10,
        },
        {
          TotalComments: 1,
        },
        {
          TotalComments: 10,
        },
      },
      wantErr: false,
    },
    {
      name: "Combination multiple post ids",
      args: args{
        ctx:     context.Background(),
        postIds: []types.Id{types.MustCreateId(), f.commentSeed[0].PostId, f.commentSeed[11].PostId, types.MustCreateId()},
      },
      want: []entity.Count{
        {
          TotalComments: 0,
        },
        {
          TotalComments: 10,
        },
        {
          TotalComments: 1,
        },
        {
          TotalComments: 0,
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

      c := &commentRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      got, err := c.GetPostCounts(tt.args.ctx, tt.args.postIds...)
      if (err != nil) != tt.wantErr {
        t.Errorf("GetPostCounts() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("GetPostCounts() got = %v, want %v", got, tt.want)
      }
    })
  }
}

func (f *commentTestSuite) Test_commentRepository_GetReplies() {
  type args struct {
    ctx       context.Context
    commentId types.Id
    showReply bool
    parameter repo.QueryParameter
  }
  tests := []struct {
    name    string
    args    args
    want    repo.PaginatedResult[entity.Comment]
    wantErr bool
  }{
    {
      name: "Valid single nested comments",
      args: args{
        ctx:       context.Background(),
        commentId: f.commentSeed[0].Id,
        showReply: true,
        parameter: repo.QueryParameter{
          Offset: 0,
          Limit:  0,
        },
      },
      want: repo.PaginatedResult[entity.Comment]{
        Data:    nil,
        Total:   4,
        Element: 0,
      },
      wantErr: false,
    },
    {
      name: "Valid single nested comments with limit and offset",
      args: args{
        ctx:       context.Background(),
        commentId: f.commentSeed[0].Id,
        showReply: true,
        parameter: repo.QueryParameter{
          Offset: 1,
          Limit:  2,
        },
      },
      want: repo.PaginatedResult[entity.Comment]{
        Data:    f.commentSeed[1:3],
        Total:   4,
        Element: 2,
      },
      wantErr: false,
    },
    {
      name: "Valid deep nested comments",
      args: args{
        ctx:       context.Background(),
        showReply: true,
        commentId: f.commentSeed[5].Id,
        parameter: repo.QueryParameter{
          Offset: 0,
          Limit:  0,
        },
      },
      want: repo.PaginatedResult[entity.Comment]{
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

      c := &commentRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      got, err := c.GetReplies(tt.args.ctx, tt.args.showReply, tt.args.commentId, tt.args.parameter)
      if (err != nil) != tt.wantErr {
        t.Errorf("GetReplies() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("GetReplies() got = %v, want %v", got, tt.want)
      }
    })
  }
}

func (f *commentTestSuite) Test_commentRepository_UpdateContent() {
  type args struct {
    ctx       context.Context
    commentId types.Id
    content   string
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

      c := &commentRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      if err := c.UpdateContent(tt.args.ctx, tt.args.commentId, tt.args.content); (err != nil) != tt.wantErr {
        t.Errorf("UpdateContent() error = %v, wantErr %v", err, tt.wantErr)
      }
    })
  }
}

func TestComment(t *testing.T) {
  suite.Run(t, &commentTestSuite{})
}

func generateComment(postId optional.Object[types.Id]) entity.Comment {
  return entity.Comment{
    Id:        types.MustCreateId(),
    PostId:    postId.ValueOr(types.MustCreateId()),
    UserId:    types.MustCreateId(),
    Content:   gofakeit.LoremIpsumSentence(5),
    CreatedAt: time.Now().UTC(),
  }
}

func generateCommentP(postId optional.Object[types.Id]) *entity.Comment {
  temp := generateComment(postId)
  return &temp
}

func ignoreReactionField(comment *entity.Comment) {
  comment.UpdatedAt = time.Time{}
  comment.CreatedAt = time.Time{}
}

func ignoreReactionsFields(comments ...entity.Comment) {
  for i := 0; i < len(comments); i++ {
    ignoreReactionField(&comments[i])
  }
}
