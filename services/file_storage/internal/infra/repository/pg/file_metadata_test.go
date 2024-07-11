package pg

import (
  "context"
  "fmt"
  sharedConf "github.com/arcorium/nexa/shared/config"
  "github.com/arcorium/nexa/shared/database"
  "github.com/arcorium/nexa/shared/optional"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  "github.com/brianvoe/gofakeit/v7"
  "github.com/stretchr/testify/suite"
  "github.com/testcontainers/testcontainers-go"
  "github.com/testcontainers/testcontainers-go/modules/postgres"
  "github.com/testcontainers/testcontainers-go/wait"
  "github.com/uptrace/bun"
  "go.opentelemetry.io/otel/trace"
  "go.opentelemetry.io/otel/trace/noop"
  rand2 "math/rand/v2"
  entity "nexa/services/file_storage/internal/domain/entity"
  "nexa/services/file_storage/internal/infra/repository/model"
  "reflect"
  "testing"
  "time"
)

const (
  USERNAME = "user"
  PASSWORD = "password"
  DATABASE = "nexa"
)

var dataSeed []entity.FileMetadata

const METADATA_SEED_SIZE = 5

type fileMetadataTestSuite struct {
  suite.Suite
  container *postgres.PostgresContainer
  db        bun.IDB
  tracer    trace.Tracer // Mock
}

func (f *fileMetadataTestSuite) SetupSuite() {
  ctx := context.Background()

  container, err := postgres.RunContainer(ctx,
    postgres.WithUsername(USERNAME),
    postgres.WithPassword(PASSWORD),
    postgres.WithDatabase(DATABASE),
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
    Username: USERNAME,
    Password: PASSWORD,
    Name:     DATABASE,
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
  err = model.SeedFromDomain(db, dataSeed...)
  f.Require().NoError(err)
}

func (f *fileMetadataTestSuite) TearDownSuite() {
  err := f.container.Terminate(context.Background())
  f.Require().NoError(err)
}

func (f *fileMetadataTestSuite) Test_metadataRepository_Create() {
  ctx := context.Background()
  type args struct {
    ctx      context.Context
    metadata *entity.FileMetadata
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "Normal",
      args: args{
        ctx:      ctx,
        metadata: generateRandomFileMetadataP(optional.NullId),
      },
      wantErr: false,
    },
    {
      name: "Duplicate Id",
      args: args{
        ctx:      ctx,
        metadata: generateRandomFileMetadataP(optional.New(&dataSeed[0].Id)),
      },
      wantErr: true,
    },
    {
      name: "Duplicate Username",
      args: args{
        ctx: ctx,
        metadata: sharedUtil.CopyWithP(dataSeed[0], func(d *entity.FileMetadata) {
          d.Id = types.MustCreateId()
        }),
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {
    f.Run(tt.name, func() {
      t := f.T()
      tx, err := f.db.BeginTx(tt.args.ctx, nil)
      f.Require().NoError(err)
      defer tx.Rollback()

      repo := metadataRepository{
        db:     tx,
        tracer: f.tracer,
      }

      err = repo.Create(tt.args.ctx, tt.args.metadata)
      if res := err != nil; res {
        if res != tt.wantErr {
          t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      got, err := repo.FindByIds(tt.args.ctx, tt.args.metadata.Id)
      f.Require().Nil(err)
      f.Require().Len(got, 1)

      // Ignore time fields
      ignoreFileMetadataFields(got...)
      ignoreFileMetadata(tt.args.metadata)

      if !reflect.DeepEqual(got[0], *tt.args.metadata) != tt.wantErr {
        t.Errorf("FindByIds() got = %v, want %v", got[0], tt.args.metadata)
      }
    })
  }
}

func (f *fileMetadataTestSuite) Test_metadataRepository_DeleteById() {
  ctx := context.Background()

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
      name: "Normal",
      args: args{
        ctx: ctx,
        id:  dataSeed[0].Id,
      },
      wantErr: false,
    },
    {
      name: "Bad Id",
      args: args{
        ctx: ctx,
        id:  types.NullId(),
      },
      wantErr: true,
    },
    {
      name: "Data not found",
      args: args{
        ctx: ctx,
        id:  types.MustCreateId(),
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {
    f.Run(tt.name, func() {
      t := f.T()
      tx, err := f.db.BeginTx(tt.args.ctx, nil)
      f.Require().NoError(err)
      defer tx.Rollback()

      repo := metadataRepository{
        db:     tx,
        tracer: f.tracer,
      }

      err = repo.DeleteById(tt.args.ctx, tt.args.id)
      if res := err != nil; res {
        if res != tt.wantErr {
          t.Errorf("DeleteById() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      got, err := repo.FindByIds(tt.args.ctx, tt.args.id)
      f.Require().Error(err)
      f.Require().Nil(got)
    })
  }
}

func (f *fileMetadataTestSuite) Test_metadataRepository_FindByIds() {
  ctx := context.Background()

  type args struct {
    ctx context.Context
    ids []types.Id
  }
  tests := []struct {
    name    string
    args    args
    want    []entity.FileMetadata
    wantErr bool
  }{
    {
      name: "Single Data",
      args: args{
        ctx: ctx,
        ids: []types.Id{dataSeed[0].Id},
      },
      want:    dataSeed[:1],
      wantErr: false,
    },
    {
      name: "Multiple Data",
      args: args{
        ctx: ctx,
        ids: []types.Id{dataSeed[0].Id, dataSeed[1].Id},
      },
      want:    dataSeed[:2],
      wantErr: false,
    },
    {
      name: "Some Id Not Found",
      args: args{
        ctx: ctx,
        ids: []types.Id{types.MustCreateId(), dataSeed[0].Id, dataSeed[1].Id, types.MustCreateId()},
      },
      want:    dataSeed[:2],
      wantErr: false,
    },
    {
      name: "File Not Found",
      args: args{
        ctx: context.Background(),
        ids: []types.Id{types.MustCreateId()},
      },
      want:    nil,
      wantErr: true,
    },
    {
      name: "Nil parameter",
      args: args{
        ctx: ctx,
        ids: nil,
      },
      want:    nil,
      wantErr: true,
    },
    {
      name: "Empty Ids",
      args: args{
        ctx: ctx,
        ids: []types.Id{},
      },
      want:    nil,
      wantErr: true,
    },
  }
  for _, tt := range tests {
    f.Run(tt.name, func() {
      t := f.T()
      tx, err := f.db.BeginTx(tt.args.ctx, nil)
      f.Require().NoError(err)
      defer tx.Rollback()

      repo := metadataRepository{
        db:     tx,
        tracer: f.tracer,
      }

      got, err := repo.FindByIds(tt.args.ctx, tt.args.ids...)
      if res := err != nil; res {
        if res != tt.wantErr {
          t.Errorf("FindByIds() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      // Ignore time fields
      ignoreFileMetadataFields(got...)
      ignoreFileMetadataFields(tt.want...)

      comparatorFunc := func(e *entity.FileMetadata, e2 *entity.FileMetadata) bool {
        return e.Id == e2.Id
      }

      if !sharedUtil.ArbitraryCheck(got, tt.want, comparatorFunc) != tt.wantErr {
        t.Errorf("FindByIds() got = %v, want %v", got, tt.want)
      }
    })
  }
}

func (f *fileMetadataTestSuite) Test_metadataRepository_FindByNames() {
  ctx := context.Background()

  type args struct {
    ctx   context.Context
    names []string
  }
  tests := []struct {
    name    string
    args    args
    want    []entity.FileMetadata
    wantErr bool
  }{
    {
      name: "Single Username",
      args: args{
        ctx:   ctx,
        names: []string{dataSeed[0].Name},
      },
      want:    dataSeed[:1],
      wantErr: false,
    },
    {
      name: "Multiple Names",
      args: args{
        ctx:   ctx,
        names: []string{dataSeed[0].Name, dataSeed[1].Name},
      },
      want:    dataSeed[:2],
      wantErr: false,
    },
    {
      name: "Some Username Not Found",
      args: args{
        ctx:   ctx,
        names: []string{gofakeit.Name(), dataSeed[0].Name, dataSeed[1].Name, gofakeit.Name()},
      },
      want:    dataSeed[:2],
      wantErr: false,
    },
    {
      name: "Username Not Found",
      args: args{
        ctx:   ctx,
        names: []string{gofakeit.Name()},
      },
      want:    nil,
      wantErr: true,
    },
    {
      name: "Multiple Names",
      args: args{
        ctx:   ctx,
        names: []string{dataSeed[0].Name, dataSeed[1].Name},
      },
      want:    dataSeed[:2],
      wantErr: false,
    },
    {
      name: "Nil parameter",
      args: args{
        ctx:   ctx,
        names: nil,
      },
      want:    nil,
      wantErr: true,
    },
    {
      name: "Empty Names",
      args: args{
        ctx:   ctx,
        names: []string{},
      },
      want:    nil,
      wantErr: true,
    },
  }
  for _, tt := range tests {
    f.Run(tt.name, func() {
      t := f.T()
      tx, err := f.db.BeginTx(tt.args.ctx, nil)
      f.Require().NoError(err)
      defer tx.Rollback()

      repo := metadataRepository{
        db:     tx,
        tracer: f.tracer,
      }

      got, err := repo.FindByNames(tt.args.ctx, tt.args.names...)
      if res := err != nil; res {
        if res != tt.wantErr {
          t.Errorf("FindByNames() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      // Ignore time fields
      ignoreFileMetadataFields(got...)
      ignoreFileMetadataFields(tt.want...)

      comparatorFunc := func(e *entity.FileMetadata, e2 *entity.FileMetadata) bool {
        return e.Id == e2.Id
      }

      if !sharedUtil.ArbitraryCheck(got, tt.want, comparatorFunc) != tt.wantErr {
        t.Errorf("FindByNames() got = %v, want %v", got, tt.want)
      }
    })
  }
}

func (f *fileMetadataTestSuite) Test_metadataRepository_Update() {
  ctx := context.Background()

  type args struct {
    ctx      context.Context
    metadata *entity.PatchedFileMetadata
    baseId   int
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "Update all fields",
      args: args{
        ctx: ctx,
        metadata: &entity.PatchedFileMetadata{
          Id:           dataSeed[0].Id,
          IsPublic:     types.SomeNullable(!dataSeed[0].IsPublic),
          Provider:     types.SomeNullable(entity.StorageProviderMinIO),
          ProviderPath: "provider/path",
          FullPath:     types.SomeNullable("other"),
        },
        baseId: 0,
      },
      wantErr: false,
    },
    {
      name: "File metadata not found",
      args: args{
        ctx: ctx,
        metadata: &entity.PatchedFileMetadata{
          Id:           types.MustCreateId(),
          IsPublic:     types.SomeNullable(!dataSeed[0].IsPublic),
          Provider:     types.SomeNullable(entity.StorageProviderMinIO),
          ProviderPath: "provider/path",
          FullPath:     types.SomeNullable("other"),
        },
        baseId: -1,
      },
      wantErr: true,
    },
    {
      name: "Update only the visibility",
      args: args{
        ctx: ctx,
        metadata: &entity.PatchedFileMetadata{
          Id:       dataSeed[0].Id,
          IsPublic: types.SomeNullable(!dataSeed[0].IsPublic),
        },
        baseId: 0,
      },
      wantErr: false,
    },
    {
      name: "Update provider",
      args: args{
        ctx: ctx,
        metadata: &entity.PatchedFileMetadata{
          Id:           dataSeed[0].Id,
          Provider:     types.SomeNullable(entity.StorageProviderMinIO),
          ProviderPath: "provider/path",
        },
        baseId: 0,
      },
      wantErr: false,
    },
    {
      name: "Update fullpath",
      args: args{
        ctx: ctx,
        metadata: &entity.PatchedFileMetadata{
          Id:       dataSeed[0].Id,
          FullPath: types.SomeNullable(gofakeit.URL()),
        },
        baseId: 0,
      },
      wantErr: false,
    },
  }
  for _, tt := range tests {
    f.Run(tt.name, func() {
      t := f.T()
      tx, err := f.db.BeginTx(tt.args.ctx, nil)
      f.Require().NoError(err)
      defer tx.Rollback()

      repo := metadataRepository{
        db:     tx,
        tracer: f.tracer,
      }

      err = repo.Patch(tt.args.ctx, tt.args.metadata)
      if res := err != nil; res {
        if res != tt.wantErr {
          t.Errorf("Patch() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      // Check result
      result, err := repo.FindByIds(tt.args.ctx, tt.args.metadata.Id)
      f.Require().NoError(err)
      f.Require().Len(result, 1)

      comparator := dataSeed[tt.args.baseId]
      if tt.args.metadata.IsPublic.HasValue() {
        comparator.IsPublic = tt.args.metadata.IsPublic.RawValue()
      }
      if tt.args.metadata.Provider.HasValue() {
        comparator.Provider = tt.args.metadata.Provider.RawValue()
      }
      if tt.args.metadata.ProviderPath != "" {
        comparator.ProviderPath = tt.args.metadata.ProviderPath
      }
      if tt.args.metadata.FullPath.HasValue() {
        comparator.FullPath = tt.args.metadata.FullPath.RawValue()
      }

      ignoreFileMetadata(&comparator)
      ignoreFileMetadata(&result[0])

      if !reflect.DeepEqual(comparator, result[0]) != tt.wantErr {
        t.Errorf("FindByIds() got = %v, want %v", result[0], comparator)
      }
    })
  }
}

func TestFileMetadata(t *testing.T) {
  for i := 0; i < METADATA_SEED_SIZE; i += 1 {
    dataSeed = append(dataSeed, generateRandomFileMetadata(optional.NullId))
  }

  suite.Run(t, &fileMetadataTestSuite{})
}

func generateRandomFileMetadata(id optional.Object[types.Id]) entity.FileMetadata {
  return entity.FileMetadata{
    Id:           id.ValueOr(types.MustCreateId()),
    Name:         gofakeit.AppName(),
    MimeType:     gofakeit.FileMimeType(),
    Size:         gofakeit.Uint64(),
    IsPublic:     gofakeit.Bool(),
    Provider:     entity.StorageProvider(rand2.UintN(uint(entity.StorageProviderAWSS3.Underlying() + 1))),
    ProviderPath: gofakeit.URL(),
    FullPath:     gofakeit.URL(),
  }
}

func generateRandomFileMetadataP(id optional.Object[types.Id]) *entity.FileMetadata {
  obj := generateRandomFileMetadata(id)
  return &obj
}

func ignoreFileMetadata(metadata *entity.FileMetadata) {
  metadata.CreatedAt = time.Time{}
  metadata.LastModified = time.Time{}
  metadata.FullPath = ""
}

func ignoreFileMetadataFields(datas ...entity.FileMetadata) {
  for i := 0; i < len(datas); i += 1 {
    ignoreFileMetadata(&datas[i])
  }
}
