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
  "nexa/services/authorization/internal/domain/entity"
  "nexa/services/authorization/internal/infra/repository/model"
  sharedConf "nexa/shared/config"
  "nexa/shared/database"
  "nexa/shared/types"
  "nexa/shared/util"
  "nexa/shared/util/repo"
  "reflect"
  "slices"
  "strconv"
  "testing"
  "time"
)

const (
  ROLE_DB_USERNAME = "user"
  ROLE_DB_PASSWORD = "password"
  ROLE_DB          = "nexa"

  SEED_ROLE_DATA_SIZE = 4
  SEED_USER_DATA_SIZE = 5
)

var roleSeed []entity.Role

var userSeed []types.Id

type roleTestSuite struct {
  suite.Suite
  container *postgres.PostgresContainer
  db        bun.IDB
  tracer    trace.Tracer // Mock
}

func (f *roleTestSuite) SetupSuite() {
  ctx := context.Background()

  container, err := postgres.RunContainer(ctx,
    postgres.WithUsername(ROLE_DB_USERNAME),
    postgres.WithPassword(ROLE_DB_PASSWORD),
    postgres.WithDatabase(ROLE_DB),
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
    Username: ROLE_DB_USERNAME,
    Password: ROLE_DB_PASSWORD,
    Name:     ROLE_DB,
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
      CreatedAt:    time.Now().Add(time.Hour * -1),
    },
    {
      RoleId:       roleSeed[0].Id.String(),
      PermissionId: permSeed[1].Id.String(),
      CreatedAt:    time.Now().Add(time.Hour * -2),
    },
    {
      RoleId:       roleSeed[1].Id.String(),
      PermissionId: permSeed[0].Id.String(),
      CreatedAt:    time.Now().Add(time.Hour * -3),
    },
    {
      RoleId:       roleSeed[2].Id.String(),
      PermissionId: permSeed[0].Id.String(),
      CreatedAt:    time.Now().Add(time.Hour * -4),
    },
  }
  userRoles := []model.UserRole{
    {
      UserId:    userSeed[0].String(),
      RoleId:    roleSeed[0].Id.String(),
      CreatedAt: time.Now().Add(time.Hour * -1),
    },
    {
      UserId:    userSeed[0].String(),
      RoleId:    roleSeed[1].Id.String(),
      CreatedAt: time.Now().Add(time.Hour * -2),
    },
    {
      UserId:    userSeed[1].String(),
      RoleId:    roleSeed[0].Id.String(),
      CreatedAt: time.Now().Add(time.Hour * -3),
    },
  }

  perms := util.CastSliceP(permSeed, func(from *entity.Permission) model.Permission {
    return model.FromPermissionDomain(from)
  })
  roleCount := 0
  roles := util.CastSliceP(roleSeed, func(from *entity.Role) model.Role {
    roleCount--
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

  err = database.Seed(f.db, userRoles...)
  f.Require().NoError(err)
}

func (f *roleTestSuite) TearDownSuite() {
  err := f.container.Terminate(context.Background())
  f.Require().NoError(err)
}

func (f *roleTestSuite) Test_roleRepository_AddPermissions() {
  type args struct {
    ctx           context.Context
    roleId        types.Id
    permissionIds []types.Id
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "Add multiple permissions to valid roles",
      args: args{
        ctx:           context.Background(),
        roleId:        roleSeed[3].Id,
        permissionIds: []types.Id{permSeed[0].Id, permSeed[1].Id},
      },
      wantErr: false,
    },
    {
      name: "Add multiple permissions to invalid roles",
      args: args{
        ctx:           context.Background(),
        roleId:        types.MustCreateId(),
        permissionIds: []types.Id{permSeed[0].Id, permSeed[1].Id},
      },
      wantErr: true,
    },
    {
      name: "Add permissions that roles already has",
      args: args{
        ctx:           context.Background(),
        roleId:        roleSeed[0].Id,
        permissionIds: []types.Id{permSeed[0].Id, permSeed[1].Id},
      },
      wantErr: true,
    },
    {
      name: "Add multiple invalid permissions to valid roles",
      args: args{
        ctx:           context.Background(),
        roleId:        roleSeed[3].Id,
        permissionIds: []types.Id{permSeed[0].Id, permSeed[1].Id, types.MustCreateId()},
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {
    f.Run(tt.name, func() {
      tx, err := f.db.BeginTx(tt.args.ctx, nil)
      f.Require().NoError(err)
      defer tx.Rollback()

      r := &roleRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      err = r.AddPermissions(tt.args.ctx, tt.args.roleId, tt.args.permissionIds...)
      if res := err != nil; res {
        if res != tt.wantErr {
          t.Errorf("AddPermissions() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      permRepo := &permissionRepository{
        db:     tx,
        tracer: f.tracer,
      }

      got, err := permRepo.FindByRoleIds(tt.args.ctx, tt.args.roleId)
      f.Require().NoError(err)

      for _, permId := range tt.args.permissionIds {
        if !slices.ContainsFunc(got, func(permission entity.Permission) bool {
          return permission.Id.Eq(permId)
        }) {
          t.Errorf("AddPermissions() error = permission %v doesn't appended", permId)
        }
      }
    })
  }
}

func (f *roleTestSuite) Test_roleRepository_AddUser() {
  type args struct {
    ctx     context.Context
    userId  types.Id
    roleIds []types.Id
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "Add single roles to user",
      args: args{
        ctx:     context.Background(),
        userId:  userSeed[0],
        roleIds: []types.Id{roleSeed[2].Id},
      },
      wantErr: false,
    },
    {
      name: "Add multiple roles to user",
      args: args{
        ctx:     context.Background(),
        userId:  userSeed[3],
        roleIds: []types.Id{roleSeed[0].Id, roleSeed[1].Id},
      },
      wantErr: false,
    },
    {
      name: "Add invalid roles to user",
      args: args{
        ctx:     context.Background(),
        userId:  userSeed[0],
        roleIds: []types.Id{types.MustCreateId()},
      },
      wantErr: true,
    },
    {
      name: "Add combination roles to user",
      args: args{
        ctx:     context.Background(),
        userId:  userSeed[3],
        roleIds: []types.Id{types.MustCreateId(), roleSeed[0].Id},
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {
    f.Run(tt.name, func() {
      tx, err := f.db.BeginTx(tt.args.ctx, nil)
      f.Require().NoError(err)
      defer tx.Rollback()

      r := &roleRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      err = r.AddUser(tt.args.ctx, tt.args.userId, tt.args.roleIds...)
      if res := err != nil; res {
        if res != tt.wantErr {
          t.Errorf("AddUser() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      got, err := r.FindByUserId(tt.args.ctx, tt.args.userId)
      f.Require().NoError(err)

      for _, roleId := range tt.args.roleIds {
        if !slices.ContainsFunc(got, func(role entity.Role) bool {
          return role.Id.Eq(roleId)
        }) {
          t.Errorf("AddUser() error = role %v doesn't appended", roleId)
        }
      }
    })
  }
}

func (f *roleTestSuite) Test_roleRepository_Create() {
  type args struct {
    ctx  context.Context
    role *entity.Role
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "Create new role",
      args: args{
        ctx:  context.Background(),
        role: generateRoleP(),
      },
      wantErr: false,
    },
    {
      name: "Create new role with duplicated name",
      args: args{
        ctx: context.Background(),
        role: util.CopyWithP(generateRole(), func(e *entity.Role) {
          e.Name = roleSeed[0].Name
        }),
      },
      wantErr: true,
    },
    {
      name: "Create new role with empty name",
      args: args{
        ctx: context.Background(),
        role: util.CopyWithP(generateRole(), func(e *entity.Role) {
          e.Name = ""
        }),
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {
    f.Run(tt.name, func() {
      tx, err := f.db.BeginTx(tt.args.ctx, nil)
      f.Require().NoError(err)
      defer tx.Rollback()

      r := &roleRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      err = r.Create(tt.args.ctx, tt.args.role)
      if res := err != nil; res {
        if res != tt.wantErr {
          t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      got, err := r.FindByIds(tt.args.ctx, tt.args.role.Id)
      f.Require().NoError(err)
      f.Require().Len(got, 1)
    })
  }
}

func (f *roleTestSuite) Test_roleRepository_Delete() {
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
      name: "Delete valid role",
      args: args{
        ctx: context.Background(),
        id:  roleSeed[0].Id,
      },
      wantErr: false,
    },
    {
      name: "Delete invalid role",
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

      r := &roleRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      err = r.Delete(tt.args.ctx, tt.args.id)
      if res := err != nil; res {
        if res != tt.wantErr {
          t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      got, err := r.FindByIds(tt.args.ctx, tt.args.id)
      f.Require().Error(err)
      f.Require().Nil(got)
    })
  }
}

func (f *roleTestSuite) Test_roleRepository_Get() {
  type args struct {
    ctx   context.Context
    query repo.QueryParameter
  }
  tests := []struct {
    name    string
    args    args
    want    repo.PaginatedResult[entity.Role]
    wantErr bool
  }{
    {
      name: "Get all roles",
      args: args{
        ctx: context.Background(),
        query: repo.QueryParameter{
          Offset: 0,
          Limit:  0,
        },
      },
      want: repo.PaginatedResult[entity.Role]{
        Data:    roleSeed,
        Total:   uint64(len(roleSeed)),
        Element: uint64(len(roleSeed)),
      },
      wantErr: false,
    },
    {
      name: "Use offset and limit",
      args: args{
        ctx: context.Background(),
        query: repo.QueryParameter{
          Offset: 1,
          Limit:  2,
        },
      },
      want: repo.PaginatedResult[entity.Role]{
        Data:    roleSeed[1:3],
        Total:   SEED_ROLE_DATA_SIZE,
        Element: 2,
      },
      wantErr: false,
    },
    {
      name: "Outside Roles count",
      args: args{
        ctx: context.Background(),
        query: repo.QueryParameter{
          Offset: 2,
          Limit:  5,
        },
      },
      want: repo.PaginatedResult[entity.Role]{
        Data:    roleSeed[2:],
        Total:   SEED_ROLE_DATA_SIZE,
        Element: 2,
      },
      wantErr: false,
    },
    {
      name: "Outside Roles offset",
      args: args{
        ctx: context.Background(),
        query: repo.QueryParameter{
          Offset: 5,
          Limit:  1,
        },
      },
      want: repo.PaginatedResult[entity.Role]{
        Data:    nil,
        Total:   SEED_ROLE_DATA_SIZE,
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

      r := &roleRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      got, err := r.Get(tt.args.ctx, tt.args.query)
      if res := err != nil; res {
        if res != tt.wantErr {
          t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      ignoreRolesFields(got.Data...)
      ignoreRolesFields(tt.want.Data...)

      if !reflect.DeepEqual(got, tt.want) != tt.wantErr {
        t.Errorf("Get() got = %v, want %v", got, tt.want)
      }
    })
  }
}

func (f *roleTestSuite) Test_roleRepository_FindByIds() {
  type args struct {
    ctx context.Context
    ids []types.Id
  }
  tests := []struct {
    name    string
    args    args
    want    []entity.Role
    wantErr bool
  }{
    //{
    //  name: "Find single valid role ids",
    //  args: args{
    //    ctx: context.Background(),
    //    ids: []types.Id{roleSeed[0].Id},
    //  },
    //  want:    roleSeed[:1],
    //  wantErr: false,
    //},
    {
      name: "Find multiple valid role ids",
      args: args{
        ctx: context.Background(),
        ids: []types.Id{roleSeed[1].Id, roleSeed[2].Id},
      },
      want:    roleSeed[1:3],
      wantErr: false,
    },
    {
      name: "Find combination multiple role ids",
      args: args{
        ctx: context.Background(),
        ids: []types.Id{roleSeed[1].Id, types.MustCreateId(), roleSeed[2].Id},
      },
      want:    roleSeed[1:3],
      wantErr: false,
    },
    {
      name: "Find all invalid role ids",
      args: args{
        ctx: context.Background(),
        ids: util.GenerateMultiple(3, types.MustCreateId),
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

      r := &roleRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      got, err := r.FindByIds(tt.args.ctx, tt.args.ids...)
      if res := err != nil; res {
        if res != tt.wantErr {
          t.Errorf("FindByIds() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      ignoreRolesFields(got...)
      ignoreRolesFields(tt.want...)

      comparatorFunc := func(e *entity.Role, e2 *entity.Role) bool { return e.Id == e2.Id }
      if !util.ArbitraryCheck(got, tt.want, comparatorFunc) != tt.wantErr {
        t.Errorf("FindByIds() got = %v, want %v", got, tt.want)
      }
    })
  }
}

func (f *roleTestSuite) Test_roleRepository_FindByName() {
  type args struct {
    ctx  context.Context
    name string
  }
  tests := []struct {
    name    string
    args    args
    want    entity.Role
    wantErr bool
  }{
    {
      name: "Find valid role by name",
      args: args{
        ctx:  context.Background(),
        name: roleSeed[0].Name,
      },
      want:    roleSeed[0],
      wantErr: false,
    },
    {
      name: "Find invalid role by name",
      args: args{
        ctx:  context.Background(),
        name: gofakeit.AppName(),
      },
      want:    entity.Role{},
      wantErr: true,
    },
  }
  for _, tt := range tests {
    f.Run(tt.name, func() {
      tx, err := f.db.BeginTx(tt.args.ctx, nil)
      f.Require().NoError(err)
      defer tx.Rollback()

      r := &roleRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      got, err := r.FindByName(tt.args.ctx, tt.args.name)
      if (err != nil) != tt.wantErr {
        t.Errorf("FindByName() error = %v, wantErr %v", err, tt.wantErr)
        return
      }

      ignoreRoleFields(&got)
      ignoreRoleFields(&tt.want)

      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("FindByName() got = %v, want %v", got, tt.want)
      }
    })
  }
}

func (f *roleTestSuite) Test_roleRepository_FindByUserId() {
  type args struct {
    ctx    context.Context
    userId types.Id
  }
  tests := []struct {
    name    string
    args    args
    want    []entity.Role
    wantErr bool
  }{
    {
      name: "Find valid user id roles",
      args: args{
        ctx:    context.Background(),
        userId: userSeed[0],
      },
      want:    roleSeed[:2],
      wantErr: false,
    },
    {
      name: "Find valid user id roles have single role",
      args: args{
        ctx:    context.Background(),
        userId: userSeed[1],
      },
      want:    roleSeed[0:1],
      wantErr: false,
    },
    {
      name: "Find user id with no roles",
      args: args{
        ctx:    context.Background(),
        userId: userSeed[2],
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

      r := &roleRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      got, err := r.FindByUserId(tt.args.ctx, tt.args.userId)
      if res := err != nil; res {
        if res != tt.wantErr {
          t.Errorf("FindByUserId() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      ignoreRolesFields(got...)
      ignoreRolesFields(tt.want...)

      comparatorFunc := func(e *entity.Role, e2 *entity.Role) bool { return e.Id == e2.Id }
      if !util.ArbitraryCheck(got, tt.want, comparatorFunc) != tt.wantErr {
        t.Errorf("FindByUserId() got = %v, want %v", got, tt.want)
      }
    })
  }
}

func (f *roleTestSuite) Test_roleRepository_Patch() {
  type args struct {
    ctx     context.Context
    role    *entity.PatchedRole
    baseIdx int
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "Update all fields",
      args: args{
        ctx: context.Background(),
        role: &entity.PatchedRole{
          Id:          roleSeed[0].Id,
          Name:        gofakeit.Country(),
          Description: types.SomeNullable(gofakeit.LoremIpsumSentence(5)),
        },
        baseIdx: 0,
      },
      wantErr: false,
    },
    {
      name: "Update role name",
      args: args{
        ctx: context.Background(),
        role: &entity.PatchedRole{
          Id:   roleSeed[0].Id,
          Name: gofakeit.Country(),
        },
        baseIdx: 0,
      },
      wantErr: false,
    },
    {
      name: "Update role description",
      args: args{
        ctx: context.Background(),
        role: &entity.PatchedRole{
          Id:          roleSeed[0].Id,
          Description: types.SomeNullable(gofakeit.LoremIpsumSentence(5)),
        },
        baseIdx: 0,
      },
      wantErr: false,
    },
    {
      name: "Remove role description",
      args: args{
        ctx: context.Background(),
        role: &entity.PatchedRole{
          Id:          roleSeed[0].Id,
          Description: types.SomeNullable(""),
        },
        baseIdx: 0,
      },
      wantErr: false,
    },
    {
      name: "Remove role name",
      args: args{
        ctx: context.Background(),
        role: &entity.PatchedRole{
          Id: roleSeed[0].Id,
        },
        baseIdx: 0,
      },
      wantErr: false, // NOTE: Because the time is updated
    },
    {
      name: "Role not found",
      args: args{
        ctx: context.Background(),
        role: &entity.PatchedRole{
          Id:          types.MustCreateId(),
          Name:        gofakeit.Country(),
          Description: types.SomeNullable(gofakeit.LoremIpsumSentence(5)),
        },
        baseIdx: -1,
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {
    f.Run(tt.name, func() {
      tx, err := f.db.BeginTx(tt.args.ctx, nil)
      f.Require().NoError(err)
      defer tx.Rollback()

      r := &roleRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      err = r.Patch(tt.args.ctx, tt.args.role)
      if res := err != nil; res {
        if res != tt.wantErr {
          t.Errorf("Patch() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      got, err := r.FindByIds(tt.args.ctx, tt.args.role.Id)
      f.Require().NoError(err)
      f.Require().Len(got, 1)

      comparator := roleSeed[tt.args.baseIdx]
      if tt.args.role.Name != "" {
        comparator.Name = tt.args.role.Name
      }
      if tt.args.role.Description.HasValue() {
        comparator.Description = tt.args.role.Description.RawValue()
      }

      ignoreRolesFields(got...)
      //ignoreRoleFields(tt.args.role)

      if !reflect.DeepEqual(got[0], comparator) != tt.wantErr {
        t.Errorf("Patch() \ngot = %v, \nwant = %v", got[0], comparator)
      }
    })
  }
}

func (f *roleTestSuite) Test_roleRepository_RemovePermissions() {
  type args struct {
    ctx           context.Context
    roleId        types.Id
    permissionIds []types.Id
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    //{
    //  name: "Role with single valid permissions",
    //  args: args{
    //    ctx:           context.Background(),
    //    roleId:        roleSeed[1].Id,
    //    permissionIds: []types.Id{permSeed[0].Id},
    //  },
    //  wantErr: false,
    //},
    //{
    //  name: "Role with multiple valid permissions",
    //  args: args{
    //    ctx:           context.Background(),
    //    roleId:        roleSeed[0].Id,
    //    permissionIds: []types.Id{permSeed[0].Id},
    //  },
    //  wantErr: false,
    //},
    //{
    //  name: "Remove all roles permissions",
    //  args: args{
    //    ctx:           context.Background(),
    //    roleId:        roleSeed[0].Id,
    //    permissionIds: []types.Id{permSeed[0].Id, permSeed[1].Id},
    //  },
    //  wantErr: false,
    //},
    {
      name: "Role with combination permissions",
      args: args{
        ctx:           context.Background(),
        roleId:        roleSeed[0].Id,
        permissionIds: []types.Id{types.MustCreateId(), permSeed[0].Id},
      },
      wantErr: false,
    },
    {
      name: "Role not found",
      args: args{
        ctx:           context.Background(),
        roleId:        types.MustCreateId(),
        permissionIds: []types.Id{permSeed[0].Id, permSeed[1].Id},
      },
      wantErr: true,
    },
    {
      name: "Role with multiple invalid permissions",
      args: args{
        ctx:           context.Background(),
        roleId:        roleSeed[0].Id,
        permissionIds: util.GenerateMultiple(2, types.MustCreateId),
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {
    f.Run(tt.name, func() {
      tx, err := f.db.BeginTx(tt.args.ctx, nil)
      f.Require().NoError(err)
      defer tx.Rollback()

      r := &roleRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      err = r.RemovePermissions(tt.args.ctx, tt.args.roleId, tt.args.permissionIds...)
      if res := err != nil; res {
        if res != tt.wantErr {
          t.Errorf("RemovePermissions() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      repo := &permissionRepository{
        db:     tx,
        tracer: f.tracer,
      }

      got, err := repo.FindByRoleIds(tt.args.ctx, tt.args.roleId)
      //if (err != nil) != tt.wantErr {
      //  t.Errorf("RemovePermissions() error = %v, wantErr %v", err, tt.wantErr)
      //  return
      //}

      for _, permId := range tt.args.permissionIds {
        if slices.ContainsFunc(got, func(perm entity.Permission) bool {
          return perm.Id.Eq(permId)
        }) {
          t.Errorf("RemovePermissions() error = role %v not removed", permId)
        }
      }

    })
  }
}

func (f *roleTestSuite) Test_roleRepository_RemoveUser() {
  type args struct {
    ctx     context.Context
    userId  types.Id
    roleIds []types.Id
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "User with valid roles",
      args: args{
        ctx:     context.Background(),
        userId:  userSeed[1],
        roleIds: []types.Id{roleSeed[0].Id},
      },
      wantErr: false,
    },
    {
      name: "User with multiple valid roles",
      args: args{
        ctx:     context.Background(),
        userId:  userSeed[0],
        roleIds: []types.Id{roleSeed[0].Id, roleSeed[1].Id},
      },
      wantErr: false,
    },
    {
      name: "User with combination roles",
      args: args{
        ctx:     context.Background(),
        userId:  userSeed[0],
        roleIds: []types.Id{types.MustCreateId(), roleSeed[0].Id},
      },
      wantErr: false,
    },
    {
      name: "User with invalid roles",
      args: args{
        ctx:     context.Background(),
        userId:  userSeed[1],
        roleIds: []types.Id{types.MustCreateId()},
      },
      wantErr: true,
    },
    {
      name: "User not found",
      args: args{
        ctx:     context.Background(),
        userId:  types.MustCreateId(),
        roleIds: []types.Id{roleSeed[0].Id},
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {
    f.Run(tt.name, func() {
      tx, err := f.db.BeginTx(tt.args.ctx, nil)
      f.Require().NoError(err)
      defer tx.Rollback()

      r := &roleRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      err = r.RemoveUser(tt.args.ctx, tt.args.userId, tt.args.roleIds...)
      if res := err != nil; res {
        if res != tt.wantErr {
          t.Errorf("RemoveUser() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      got, err := r.FindByUserId(tt.args.ctx, tt.args.userId)

      ignoreRolesFields(got...)

      for _, roleId := range tt.args.roleIds {
        if slices.ContainsFunc(got, func(role entity.Role) bool {
          return role.Id.Eq(roleId)
        }) {
          t.Errorf("RemoveUser() error = role %v not removed", roleId)
        }
      }

    })
  }
}

func TestRoles(t *testing.T) {
  seedPermsData()
  seedRolesData()
  seedUserData()

  suite.Run(t, &roleTestSuite{})
}

func seedRolesData() {
  for i := 0; i < SEED_ROLE_DATA_SIZE; i += 1 {
    roleSeed = append(roleSeed, generateRole())
  }
}

func seedUserData() {
  for i := 0; i < SEED_USER_DATA_SIZE; i += 1 {
    userSeed = append(userSeed, types.MustCreateId())
  }
}

func generateRole() entity.Role {
  return entity.Role{
    Id:          types.MustCreateId(),
    Name:        gofakeit.Country(),
    Description: gofakeit.LoremIpsumSentence(5),
  }
}

func generateRoleP() *entity.Role {
  temp := generateRole()
  return &temp
}

func ignoreRoleFields(role *entity.Role) {
  role.Permissions = nil
}

func ignoreRolesFields(role ...entity.Role) {
  for i := 0; i < len(role); i++ {
    ignoreRoleFields(&role[i])
  }
}
