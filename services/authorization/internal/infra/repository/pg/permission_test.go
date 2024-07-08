package pg

import (
  "context"
  "fmt"
  "github.com/brianvoe/gofakeit/v7"
  "github.com/stretchr/testify/suite"
  "github.com/testcontainers/testcontainers-go"
  "github.com/testcontainers/testcontainers-go/modules/postgres"
  "github.com/testcontainers/testcontainers-go/wait"
  "github.com/uptrace/bun"
  "go.opentelemetry.io/otel/trace"
  "go.opentelemetry.io/otel/trace/noop"
  "nexa/services/authorization/constant"
  "nexa/services/authorization/internal/domain/entity"
  "nexa/services/authorization/internal/infra/repository/model"
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
  PERM_DB_USERNAME = "user"
  PERM_DB_PASSWORD = "password"
  PERM_DB          = "nexa"

  SEED_PERM_DATA_SIZE = 5
)

var permSeed []entity.Permission

func ignorePermsField(got *entity.Permission) {
  got.CreatedAt = time.Time{}
}

func ignorePermsFields(got ...entity.Permission) {
  // Ignore time fields
  for i := 0; i < len(got); i += 1 {
    ignorePermsField(&got[i])
  }
}

type permTestSuite struct {
  suite.Suite
  container *postgres.PostgresContainer
  db        bun.IDB
  tracer    trace.Tracer // Mock
}

func (f *permTestSuite) SetupSuite() {
  ctx := context.Background()

  container, err := postgres.RunContainer(ctx,
    postgres.WithUsername(PERM_DB_USERNAME),
    postgres.WithPassword(PERM_DB_PASSWORD),
    postgres.WithDatabase(PERM_DB),
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
    Username: PERM_DB_USERNAME,
    Password: PERM_DB_PASSWORD,
    Name:     PERM_DB,
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
  rolePerms := []model.RolePermission{
    {
      RoleId:       roleSeed[0].Id.String(),
      PermissionId: permSeed[0].Id.String(),
      CreatedAt:    time.Now().Add(time.Hour * 1),
    },
    {
      RoleId:       roleSeed[0].Id.String(),
      PermissionId: permSeed[1].Id.String(),
      CreatedAt:    time.Now().Add(time.Hour * 2),
    },
    {
      RoleId:       roleSeed[1].Id.String(),
      PermissionId: permSeed[0].Id.String(),
      CreatedAt:    time.Now().Add(time.Hour * 3),
    },
    {
      RoleId:       roleSeed[2].Id.String(),
      PermissionId: permSeed[0].Id.String(),
      CreatedAt:    time.Now().Add(time.Hour * 4),
    },
  }
  perms := util.CastSliceP(permSeed, func(from *entity.Permission) model.Permission {
    return model.FromPermissionDomain(from)
  })
  roleCount := 0
  roles := util.CastSliceP(roleSeed, func(from *entity.Role) model.Role {
    roleCount++
    return model.FromRoleDomain(from, func(ent *entity.Role, role *model.Role) {
      role.CreatedAt = time.Now().Add(time.Duration(roleCount) * time.Hour).UTC()
    })
  })

  err = database.Seed(f.db, perms...)
  f.Require().NoError(err)

  err = database.Seed(f.db, roles...)
  f.Require().NoError(err)

  err = database.Seed(f.db, rolePerms...)
  f.Require().NoError(err)
}

func (f *permTestSuite) TearDownSuite() {
  err := f.container.Terminate(context.Background())
  f.Require().NoError(err)
}

func (f *permTestSuite) Test_permissionRepository_Create() {
  type args struct {
    ctx        context.Context
    permission *entity.Permission
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "Create valid permission",
      args: args{
        ctx: context.Background(),
        permission: &entity.Permission{
          Id:        types.MustCreateId(),
          Resource:  constant.SERVICE_RESOURCE,
          Action:    gofakeit.AnimalType(),
          CreatedAt: time.Now(),
        },
      },
      wantErr: false,
    },
    {
      name: "Create duplicate permission",
      args: args{
        ctx: context.Background(),
        permission: &entity.Permission{
          Id:        types.MustCreateId(),
          Resource:  permSeed[0].Resource,
          Action:    permSeed[0].Action,
          CreatedAt: time.Now(),
        },
      },
      wantErr: true,
    },
    {
      name: "Create permission with null on nullable fields",
      args: args{
        ctx: context.Background(),
        permission: &entity.Permission{
          Id:        types.MustCreateId(),
          CreatedAt: time.Now(),
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

      p := &permissionRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      err = p.Create(tt.args.ctx, tt.args.permission)
      if res := err != nil; res {
        if res != tt.wantErr {
          t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      got, err := p.FindByIds(tt.args.ctx, tt.args.permission.Id)
      f.Require().NoError(err)
      f.Require().Len(got, 1)

      ignorePermsFields(got...)
      ignorePermsField(tt.args.permission)

      if !reflect.DeepEqual(got[0], *tt.args.permission) != tt.wantErr {
        t.Errorf("Get() got = %v, want %v", got[0], *tt.args.permission)
      }
    })
  }
}

func (f *permTestSuite) Test_permissionRepository_Delete() {
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
      name: "Delete valid permission",
      args: args{
        ctx: context.Background(),
        id:  permSeed[0].Id,
      },
      wantErr: false,
    },
    {
      name: "Delete invalid permission",
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

      p := &permissionRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      err = p.Delete(tt.args.ctx, tt.args.id)
      if res := err != nil; res {
        if res != tt.wantErr {
          t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      got, err := p.FindByIds(tt.args.ctx, tt.args.id)
      f.Require().Error(err)
      f.Require().Nil(got)
    })
  }
}

func (f *permTestSuite) Test_permissionRepository_Get() {
  type args struct {
    ctx   context.Context
    query repo.QueryParameter
  }
  tests := []struct {
    name    string
    args    args
    want    repo.PaginatedResult[entity.Permission]
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
      want: repo.PaginatedResult[entity.Permission]{
        Data:    permSeed,
        Total:   uint64(len(permSeed)),
        Element: uint64(len(permSeed)),
      },
      wantErr: false,
    },
    {
      name: "Get all with offset and limit",
      args: args{
        ctx: context.Background(),
        query: repo.QueryParameter{
          Offset: 1,
          Limit:  2,
        },
      },
      want: repo.PaginatedResult[entity.Permission]{
        Data:    permSeed[1:3],
        Total:   uint64(len(permSeed)),
        Element: 2,
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
      want: repo.PaginatedResult[entity.Permission]{
        Data:    permSeed[2:],
        Total:   uint64(len(permSeed)),
        Element: uint64(len(permSeed)) - 2,
      },
      wantErr: false,
    },
    {
      name: "Get all with limit",
      args: args{
        ctx: context.Background(),
        query: repo.QueryParameter{
          Offset: 0,
          Limit:  3,
        },
      },
      want: repo.PaginatedResult[entity.Permission]{
        Data:    permSeed[:3],
        Total:   uint64(len(permSeed)),
        Element: 3,
      },
      wantErr: false,
    },
    {
      name: "Get out of bound offset",
      args: args{
        ctx: context.Background(),
        query: repo.QueryParameter{
          Offset: uint64(len(permSeed)),
          Limit:  3,
        },
      },
      want: repo.PaginatedResult[entity.Permission]{
        Data:    nil,
        Total:   uint64(len(permSeed)),
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

      p := &permissionRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      got, err := p.Get(tt.args.ctx, tt.args.query)
      if (err != nil) != tt.wantErr {
        t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
        return
      }

      ignorePermsFields(tt.want.Data...)
      ignorePermsFields(got.Data...)

      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("Get() got = %v, want %v", got, tt.want)
      }
    })
  }
}

func (f *permTestSuite) Test_permissionRepository_FindByIds() {
  type args struct {
    ctx context.Context
    ids []types.Id
  }
  tests := []struct {
    name    string
    args    args
    want    []entity.Permission
    wantErr bool
  }{
    {
      name: "Get single perm",
      args: args{
        ctx: context.Background(),
        ids: []types.Id{permSeed[0].Id},
      },
      want:    permSeed[:1],
      wantErr: false,
    },
    {
      name: "Get multiple perms",
      args: args{
        ctx: context.Background(),
        ids: []types.Id{permSeed[0].Id, permSeed[1].Id},
      },
      want:    permSeed[:2],
      wantErr: false,
    },
    {
      name: "Some perm is not valid",
      args: args{
        ctx: context.Background(),
        ids: []types.Id{permSeed[2].Id, types.MustCreateId(), permSeed[1].Id},
      },
      want:    []entity.Permission{permSeed[1], permSeed[2]},
      wantErr: false,
    },
    {
      name: "perm not found",
      args: args{
        ctx: context.Background(),
        ids: []types.Id{types.MustCreateId()},
      },
      want:    nil,
      wantErr: true,
    },
    {
      name: "All perm not found",
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

      p := &permissionRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      got, err := p.FindByIds(tt.args.ctx, tt.args.ids...)
      if res := err != nil; res {
        if res != tt.wantErr {
          t.Errorf("FindByIds() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      ignorePermsFields(tt.want...)
      ignorePermsFields(got...)

      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("FindByIds() got = %v, want %v", got, tt.want)
      }
    })
  }
}

func (f *permTestSuite) Test_permissionRepository_FindByRoleIds() {
  type args struct {
    ctx     context.Context
    roleIds []types.Id
  }
  tests := []struct {
    name    string
    args    args
    want    []entity.Permission
    wantErr bool
  }{
    {
      name: "Get single role permissions",
      args: args{
        ctx:     context.Background(),
        roleIds: []types.Id{roleSeed[0].Id},
      },
      want:    []entity.Permission{permSeed[0], permSeed[1]},
      wantErr: false,
    },
    {
      name: "Get multiple role permissions",
      args: args{
        ctx:     context.Background(),
        roleIds: []types.Id{roleSeed[0].Id, roleSeed[1].Id},
      },
      want:    []entity.Permission{permSeed[0], permSeed[1]},
      wantErr: false,
    },
    {
      name: "Role doesn't have permission",
      args: args{
        ctx:     context.Background(),
        roleIds: []types.Id{roleSeed[3].Id, types.MustCreateId()},
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

      p := &permissionRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      got, err := p.FindByRoleIds(tt.args.ctx, tt.args.roleIds...)
      if res := err != nil; res {
        if res != tt.wantErr {
          t.Errorf("FindByRoleIds() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      ignorePermsFields(tt.want...)
      ignorePermsFields(got...)

      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("FindByRoleIds() got = %v, want %v", got, tt.want)
      }
    })
  }
}

func TestPerms(t *testing.T) {
  seedPermsData()
  seedRolesData()

  suite.Run(t, &permTestSuite{})
}

func seedPermsData() {
  for i := 0; i < SEED_PERM_DATA_SIZE; i += 1 {
    permSeed = append(permSeed, generatePerms())
  }
}

var generatePermsCount = 0

func generatePerms() entity.Permission {
  generatePermsCount++
  return entity.Permission{
    Id:        types.MustCreateId(),
    Resource:  gofakeit.AnimalType(),
    Action:    fmt.Sprintf("%s:%s", gofakeit.Animal(), gofakeit.Username()),
    CreatedAt: time.Now().Add(time.Duration(generatePermsCount) * time.Hour).UTC(),
  }
}

func generatePermsP() *entity.Permission {
  user := generatePerms()
  return &user
}
