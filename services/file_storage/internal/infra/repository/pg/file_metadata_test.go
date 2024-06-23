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
  rand2 "math/rand/v2"
  domain "nexa/services/file_storage/internal/domain/entity"
  "nexa/services/file_storage/internal/infra/repository/model"
  sharedConf "nexa/shared/config"
  "nexa/shared/database"
  "nexa/shared/optional"
  "nexa/shared/types"
  sharedUtil "nexa/shared/util"
  "nexa/shared/wrapper"
  "reflect"
  "strconv"
  "testing"
  "time"
)

const (
  USERNAME = "user"
  PASSWORD = "password"
  DATABASE = "nexa"
)

var dataSeed []domain.FileMetadata

const METADATA_SEED_SIZE = 5

type noopLogger struct{}

func (n noopLogger) Printf(format string, v ...interface{}) {

}

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
    testcontainers.WithLogger(&noopLogger{}),
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

  sharedUtil.DoNothing(ports, mapped)

  db, err := database.OpenPostgres(&sharedConf.Database{
    Protocol: "postgres",
    Host:     wrapper.PanicDropError(container.Host(ctx)),
    Port:     uint16(wrapper.PanicDropError(strconv.Atoi(mapped[0].HostPort))),
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
  err = model.SeedDatabase(db, dataSeed...)
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
    metadata *domain.FileMetadata
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
      name: "Duplicate Name",
      args: args{
        ctx: ctx,
        metadata: sharedUtil.CopyWithP(dataSeed[0], func(d *domain.FileMetadata) {
          d.Id = types.NewId2()
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

      if err := repo.Create(tt.args.ctx, tt.args.metadata); (err != nil) != tt.wantErr {
        t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
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
        id:  types.NewId2(),
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

      if err := repo.DeleteById(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
        t.Errorf("DeleteById() error = %v, wantErr %v", err, tt.wantErr)
      }
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
    want    []domain.FileMetadata
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
        ids: []types.Id{types.NewId2(), dataSeed[0].Id, dataSeed[1].Id, types.NewId2()},
      },
      want:    dataSeed[:2],
      wantErr: false,
    },
    {
      name: "File Not Found",
      args: args{
        ctx: ctx,
        ids: []types.Id{types.NewId2()},
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
      if (err != nil) != tt.wantErr {
        t.Errorf("FindByIds() error = %v, wantErr %v", err, tt.wantErr)
        return
      }

      // Ignore time fields
      ignoreUnimportantFields(got...)
      ignoreUnimportantFields(tt.want...)

      if !reflect.DeepEqual(got, tt.want) {
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
    want    []domain.FileMetadata
    wantErr bool
  }{
    {
      name: "Single Name",
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
      name: "Some Name Not Found",
      args: args{
        ctx:   ctx,
        names: []string{gofakeit.Name(), dataSeed[0].Name, dataSeed[1].Name, gofakeit.Name()},
      },
      want:    dataSeed[:2],
      wantErr: false,
    },
    {
      name: "Name Not Found",
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
      if (err != nil) != tt.wantErr {
        t.Errorf("FindByNames() error = %v, wantErr %v", err, tt.wantErr)
        return
      }

      // Ignore time fields
      ignoreUnimportantFields(got...)
      ignoreUnimportantFields(tt.want...)

      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("FindByNames() got = %v, want %v", got, tt.want)
      }
    })
  }
}

func (f *fileMetadataTestSuite) Test_metadataRepository_Update() {
  ctx := context.Background()

  type args struct {
    ctx      context.Context
    metadata *domain.FileMetadata
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
        metadata: &domain.FileMetadata{
          Id:       dataSeed[0].Id,
          Name:     "something.jpg",
          Type:     domain.FileTypeImage,
          Provider: domain.StorageProviderMinIO,
          IsPublic: !dataSeed[0].IsPublic,
        },
        //metadata: sharedUtil.CopyWithP(dataSeed[0], func(d *domain.FileMetadata) {
        //  d.Name = "something.jpg"
        //  d.IsPublic = !d.IsPublic
        //}),
      },
      wantErr: false,
    },
    {
      name: "Data Not Found",
      args: args{
        ctx: ctx,
        metadata: &domain.FileMetadata{
          Id:       types.NewId2(),
          Name:     "something.jpg",
          IsPublic: !dataSeed[0].IsPublic,
        },
      },
      wantErr: true,
    },
    {
      name: "Change Name Into Duplicate",
      args: args{
        ctx: ctx,
        metadata: &domain.FileMetadata{
          Id:       dataSeed[0].Id,
          Name:     dataSeed[1].Name,
          IsPublic: !dataSeed[0].IsPublic,
        },
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

      err = repo.Update(tt.args.ctx, tt.args.metadata)

      if err != nil {
        if (err != nil) != tt.wantErr {
          t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      // Check result
      result, err := repo.FindByIds(tt.args.ctx, tt.args.metadata.Id)
      f.Require().NoError(err)

      if result[0].Name != tt.args.metadata.Name ||
          result[0].IsPublic != tt.args.metadata.IsPublic ||
          result[0].Type != tt.args.metadata.Type ||
          result[0].Provider != tt.args.metadata.Provider {

        t.Errorf("Update() got = %v, want %v", result[0], tt.args.metadata)
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

func generateRandomFileMetadata(id optional.Object[types.Id]) domain.FileMetadata {
  return domain.FileMetadata{
    Id:           id.ValueOr(types.NewId2()),
    Name:         gofakeit.AppName(),
    Type:         domain.FileType(rand2.UintN(uint(domain.FileTypeOther.Underlying() + 1))),
    Size:         gofakeit.Uint64(),
    IsPublic:     gofakeit.Bool(),
    Provider:     domain.StorageProvider(rand2.UintN(uint(domain.StorageProviderAWSS3.Underlying() + 1))),
    ProviderPath: gofakeit.URL(),
    FullPath:     gofakeit.URL(),
  }
}

func generateRandomFileMetadataP(id optional.Object[types.Id]) *domain.FileMetadata {
  obj := generateRandomFileMetadata(id)
  return &obj
}

func ignoreUnimportantFields(datas ...domain.FileMetadata) {
  for i := 0; i < len(datas); i += 1 {
    datas[i].CreatedAt = time.Time{}
    datas[i].LastModified = time.Time{}
    datas[i].FullPath = ""
  }
}
