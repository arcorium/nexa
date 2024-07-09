package pg

import (
  "context"
  "github.com/brianvoe/gofakeit/v7"
  "github.com/stretchr/testify/suite"
  "github.com/testcontainers/testcontainers-go"
  "github.com/testcontainers/testcontainers-go/modules/postgres"
  "github.com/testcontainers/testcontainers-go/wait"
  "github.com/uptrace/bun"
  "go.opentelemetry.io/otel/trace"
  "go.opentelemetry.io/otel/trace/noop"
  "nexa/services/mailer/internal/domain/entity"
  "nexa/services/mailer/internal/infra/repository/model"
  sharedConf "nexa/shared/config"
  "nexa/shared/database"
  "nexa/shared/types"
  "nexa/shared/util"
  "nexa/shared/util/repo"
  "reflect"
  "strconv"
  "testing"
  "time"
)

const (
  TAG_DB_USERNAME = "user"
  TAG_DB_PASSWORD = "password"
  TAG_DB          = "nexa"

  SEED_TAG_DATA_SIZE = 3
)

type tagTestSuite struct {
  suite.Suite
  container *postgres.PostgresContainer
  db        bun.IDB
  tracer    trace.Tracer // Mock

  tagSeed []entity.Tag
}

func (f *tagTestSuite) SetupSuite() {
  // Create data
  for i := 0; i < SEED_TAG_DATA_SIZE; i += 1 {
    f.tagSeed = append(f.tagSeed, generateTag())
  }

  ctx := context.Background()

  container, err := postgres.RunContainer(ctx,
    postgres.WithUsername(TAG_DB_USERNAME),
    postgres.WithPassword(TAG_DB_PASSWORD),
    postgres.WithDatabase(TAG_DB),
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

  db, err := database.OpenPostgres(&sharedConf.Database{
    Protocol: "postgres",
    Host:     types.Must(container.Host(ctx)),
    Port:     uint16(types.Must(strconv.Atoi(mapped[0].HostPort))),
    Username: TAG_DB_USERNAME,
    Password: TAG_DB_PASSWORD,
    Name:     TAG_DB,
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
  // Tag
  counter := 0
  tags := util.CastSliceP(f.tagSeed, func(from *entity.Tag) model.Tag {
    counter += 1
    return model.FromTagDomain(from, func(ent *entity.Tag, tag *model.Tag) {
      tag.CreatedAt = time.Now().Add(time.Hour * time.Duration(counter*-1))
    })
  })

  err = database.Seed(f.db, tags...)
  f.Require().NoError(err)
}

func (f *tagTestSuite) TearDownSuite() {
  err := f.container.Terminate(context.Background())
  f.Require().NoError(err)
}

func (f *tagTestSuite) Test_tagRepository_Create() {
  type args struct {
    ctx context.Context
    tag *entity.Tag
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "Create valid tag",
      args: args{
        ctx: context.Background(),
        tag: &entity.Tag{
          Id:          types.MustCreateId(),
          Name:        gofakeit.AnimalType(),
          Description: gofakeit.LoremIpsumSentence(10),
        },
      },
      wantErr: false,
    },
    {
      name: "Create tag with empty name",
      args: args{
        ctx: context.Background(),
        tag: &entity.Tag{
          Id:          types.MustCreateId(),
          Description: gofakeit.LoremIpsumSentence(10),
        },
      },
      wantErr: true,
    },
    {
      name: "Create duplicated id tag",
      args: args{
        ctx: context.Background(),
        tag: &entity.Tag{
          Id:          f.tagSeed[0].Id,
          Name:        gofakeit.AnimalType(),
          Description: gofakeit.LoremIpsumSentence(10),
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

      t := tagRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t1 := f.T()

      err = t.Create(tt.args.ctx, tt.args.tag)
      if res := err != nil; res {
        if res != tt.wantErr {
          t1.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      gotTag, err := t.FindByIds(tt.args.ctx, tt.args.tag.Id)
      f.Require().NoError(err)
      f.Require().Len(gotTag, 1)

      if !reflect.DeepEqual(gotTag[0], *tt.args.tag) != tt.wantErr {
        t1.Errorf("Get() got = %v, want %v", gotTag[0], *tt.args.tag)
      }
    })
  }
}

func (f *tagTestSuite) Test_tagRepository_Get() {
  type args struct {
    ctx   context.Context
    query repo.QueryParameter
  }
  tests := []struct {
    name    string
    args    args
    want    repo.PaginatedResult[entity.Tag]
    wantErr bool
  }{
    {
      name: "Get all",
      args: args{
        ctx: context.Background(),
        query: repo.QueryParameter{
          Offset: 0,
          Limit:  0,
        },
      },
      want: repo.PaginatedResult[entity.Tag]{
        Data:    f.tagSeed,
        Total:   uint64(len(f.tagSeed)),
        Element: uint64(len(f.tagSeed)),
      },
      wantErr: false,
    },
    {
      name: "Get all with offset and limit",
      args: args{
        ctx: context.Background(),
        query: repo.QueryParameter{
          Offset: 1,
          Limit:  1,
        },
      },
      want: repo.PaginatedResult[entity.Tag]{
        Data:    f.tagSeed[1:2],
        Total:   uint64(len(f.tagSeed)),
        Element: 1,
      },
      wantErr: false,
    },
    {
      name: "Get all with offset",
      args: args{
        ctx: context.Background(),
        query: repo.QueryParameter{
          Offset: 2,
          Limit:  0,
        },
      },
      want: repo.PaginatedResult[entity.Tag]{
        Data:    f.tagSeed[2:],
        Total:   uint64(len(f.tagSeed)),
        Element: uint64(len(f.tagSeed)) - 2,
      },
      wantErr: false,
    },
    {
      name: "Get all with limit",
      args: args{
        ctx: context.Background(),
        query: repo.QueryParameter{
          Offset: 0,
          Limit:  2,
        },
      },
      want: repo.PaginatedResult[entity.Tag]{
        Data:    f.tagSeed[:2],
        Total:   uint64(len(f.tagSeed)),
        Element: 2,
      },
      wantErr: false,
    },
    {
      name: "Get out of bound offset",
      args: args{
        ctx: context.Background(),
        query: repo.QueryParameter{
          Offset: uint64(len(f.tagSeed)),
          Limit:  3,
        },
      },
      want: repo.PaginatedResult[entity.Tag]{
        Data:    nil,
        Total:   uint64(len(f.tagSeed)),
        Element: 0,
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {
    f.Run(tt.name, func() {
      tx, err := f.db.BeginTx(tt.args.ctx, nil)
      f.Require().NoError(err)
      defer tx.Rollback()

      t := tagRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t1 := f.T()

      got, err := t.Get(tt.args.ctx, tt.args.query)
      if (err != nil) != tt.wantErr {
        t1.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if !reflect.DeepEqual(got, tt.want) {
        t1.Errorf("Get() got = %v, want %v", got, tt.want)
      }
    })
  }
}

func (f *tagTestSuite) Test_tagRepository_FindByIds() {
  type args struct {
    ctx context.Context
    ids []types.Id
  }
  tests := []struct {
    name    string
    args    args
    want    []entity.Tag
    wantErr bool
  }{
    {
      name: "Get single tag",
      args: args{
        ctx: context.Background(),
        ids: []types.Id{f.tagSeed[0].Id},
      },
      want:    f.tagSeed[:1],
      wantErr: false,
    },
    {
      name: "Get multiple mails",
      args: args{
        ctx: context.Background(),
        ids: []types.Id{f.tagSeed[0].Id, f.tagSeed[1].Id},
      },
      want:    f.tagSeed[:2],
      wantErr: false,
    },
    {
      name: "Some mail is not valid",
      args: args{
        ctx: context.Background(),
        ids: []types.Id{f.tagSeed[2].Id, types.MustCreateId(), f.tagSeed[1].Id},
      },
      want:    []entity.Tag{f.tagSeed[1], f.tagSeed[2]},
      wantErr: false,
    },
    {
      name: "Mail not found",
      args: args{
        ctx: context.Background(),
        ids: []types.Id{types.MustCreateId()},
      },
      want:    nil,
      wantErr: true,
    },
    {
      name: "All mail not found",
      args: args{
        ctx: context.Background(),
        ids: []types.Id{types.MustCreateId(), types.MustCreateId(), types.MustCreateId()},
      },
      want:    nil,
      wantErr: true,
    },
  }
  for _, tt := range tests {
    f.Run(tt.name, func() {
      tx, err := f.db.BeginTx(tt.args.ctx, nil)
      f.Require().NoError(err)
      defer tx.Rollback()

      t := tagRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t1 := f.T()

      got, err := t.FindByIds(tt.args.ctx, tt.args.ids...)
      if (err != nil) != tt.wantErr {
        t1.Errorf("FindByIds() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if !reflect.DeepEqual(got, tt.want) {
        t1.Errorf("FindByIds() got = %v, want %v", got, tt.want)
      }
    })
  }
}

func (f *tagTestSuite) Test_tagRepository_FindByName() {
  type args struct {
    ctx  context.Context
    name string
  }
  tests := []struct {
    name    string
    args    args
    want    *entity.Tag
    wantErr bool
  }{
    {
      name: "Valid name",
      args: args{
        ctx:  context.Background(),
        name: f.tagSeed[0].Name,
      },
      want:    &f.tagSeed[0],
      wantErr: false,
    },
    {
      name: "Invalid name",
      args: args{
        ctx:  context.Background(),
        name: gofakeit.AppName(),
      },
      want:    nil,
      wantErr: true,
    },
  }
  for _, tt := range tests {
    f.Run(tt.name, func() {
      tx, err := f.db.BeginTx(tt.args.ctx, nil)
      f.Require().NoError(err)
      defer tx.Rollback()

      t := tagRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t1 := f.T()

      got, err := t.FindByName(tt.args.ctx, tt.args.name)
      if (err != nil) != tt.wantErr {
        t1.Errorf("FindByName() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if !reflect.DeepEqual(got, tt.want) {
        t1.Errorf("FindByName() got = %v, want %v", got, tt.want)
      }
    })
  }
}

func (f *tagTestSuite) Test_tagRepository_Patch() {
  type args struct {
    ctx    context.Context
    tag    *entity.PatchedTag
    baseId int
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "Update all tag fields",
      args: args{
        ctx: context.Background(),
        tag: &entity.PatchedTag{
          Id:          f.tagSeed[0].Id,
          Name:        "another",
          Description: types.SomeNullable(gofakeit.LoremIpsumSentence(20)),
        },
        baseId: 0,
      },
      wantErr: false,
    },
    {
      name: "Update name tag fields",
      args: args{
        ctx: context.Background(),
        tag: &entity.PatchedTag{
          Id:   f.tagSeed[0].Id,
          Name: "another",
        },
        baseId: 0,
      },
      wantErr: false,
    },
    {
      name: "Set description to empty",
      args: args{
        ctx: context.Background(),
        tag: &entity.PatchedTag{
          Id:          f.tagSeed[0].Id,
          Description: types.SomeNullable(""),
        },
        baseId: 0,
      },
      wantErr: false,
    },
    {
      name: "Set name to empty",
      args: args{
        ctx: context.Background(),
        tag: &entity.PatchedTag{
          Id: f.tagSeed[0].Id,
        },
        baseId: 0,
      },
      wantErr: false, // Name is not modified, but the updated_at still does
    },
    {
      name: "Tag not found",
      args: args{
        ctx: context.Background(),
        tag: &entity.PatchedTag{
          Id:   types.MustCreateId(),
          Name: "other",
        },
        baseId: -1,
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {
    f.Run(tt.name, func() {
      tx, err := f.db.BeginTx(tt.args.ctx, nil)
      f.Require().NoError(err)
      defer tx.Rollback()

      t := tagRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t1 := f.T()

      err = t.Patch(tt.args.ctx, tt.args.tag)
      if res := err != nil; res {
        if res != tt.wantErr {
          t1.Errorf("Patch() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      got, err := t.FindByIds(tt.args.ctx, tt.args.tag.Id)
      f.Require().NoError(err)
      f.Require().Len(got, 1)

      comparator := f.tagSeed[tt.args.baseId]
      if tt.args.tag.Name != "" {
        comparator.Name = tt.args.tag.Name
      }
      if tt.args.tag.Description.HasValue() {
        comparator.Description = tt.args.tag.Description.RawValue()
      }

      if !reflect.DeepEqual(comparator, got[0]) != tt.wantErr {
        t1.Errorf("FindByIds() got = %v, want %v", got[0], comparator)
      }
    })
  }
}

func (f *tagTestSuite) Test_tagRepository_Remove() {
  type args struct {
    ctx context.Context
    id  types.Id
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "Remove valid tag",
      args: args{
        ctx: context.Background(),
        id:  f.tagSeed[0].Id,
      },
      wantErr: false,
    },
    {
      name: "Remove invalid tag",
      args: args{
        ctx: context.Background(),
        id:  types.MustCreateId(),
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {
    f.Run(tt.name, func() {
      tx, err := f.db.BeginTx(tt.args.ctx, nil)
      f.Require().NoError(err)
      defer tx.Rollback()

      t := tagRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t1 := f.T()

      err = t.Remove(tt.args.ctx, tt.args.id)
      if res := err != nil; res {
        if res != tt.wantErr {
          t1.Errorf("Remove() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      got, err := t.FindByIds(tt.args.ctx, tt.args.id)
      f.Require().Error(err)
      f.Require().Nil(got)
    })
  }
}

func TestTag(t *testing.T) {
  suite.Run(t, &tagTestSuite{})
}

func generateTag() entity.Tag {
  return entity.Tag{
    Id:          types.MustCreateId(),
    Name:        gofakeit.EmojiTag(),
    Description: gofakeit.LoremIpsumSentence(3),
  }
}
