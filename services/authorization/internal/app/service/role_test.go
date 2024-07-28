package service

import (
  "context"
  "database/sql"
  sharedDto "github.com/arcorium/nexa/shared/dto"
  "github.com/arcorium/nexa/shared/status"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  "github.com/arcorium/nexa/shared/util/repo"
  "github.com/brianvoe/gofakeit/v7"
  "github.com/stretchr/testify/mock"
  "github.com/stretchr/testify/require"
  "go.opentelemetry.io/otel/trace"
  "go.opentelemetry.io/otel/trace/noop"
  "nexa/services/authorization/constant"
  "nexa/services/authorization/internal/domain/dto"
  "nexa/services/authorization/internal/domain/entity"
  repoMock "nexa/services/authorization/internal/domain/repository/mocks"
  "nexa/services/authorization/util/errors"
  "reflect"
  "testing"
)

type roleMocked struct {
  Role   *repoMock.RoleMock
  Tracer trace.Tracer
}

func newRoleMocked(t *testing.T) roleMocked {
  // Tracer
  provider := noop.NewTracerProvider()
  return roleMocked{
    Role:   repoMock.NewRoleMock(t),
    Tracer: provider.Tracer("MOCK"),
  }
}

type setupRoleTestFunc func(mocked *roleMocked, arg any, want any)

func Test_roleService_AddPermissions(t *testing.T) {
  type args struct {
    ctx            context.Context
    permissionsDTO *dto.ModifyRolesPermissionsDTO
  }
  tests := []struct {
    name  string
    setup setupRoleTestFunc
    args  args
    want  status.Object
  }{
    {
      name: "Success add role's permissions",
      setup: func(mocked *roleMocked, arg any, want any) {
        a := arg.(*args)

        permIds := sharedUtil.CastSlice(a.permissionsDTO.PermissionIds, sharedUtil.ToAny[types.Id])

        mocked.Role.EXPECT().
          AddPermissions(mock.Anything, a.permissionsDTO.RoleId, permIds...).
          Return(nil)
      },
      args: args{
        ctx: context.Background(),
        permissionsDTO: &dto.ModifyRolesPermissionsDTO{
          RoleId:        types.MustCreateId(),
          PermissionIds: sharedUtil.GenerateMultiple(3, types.MustCreateId),
        },
      },
      want: status.Created(),
    },
    {
      name: "Role already has the permissions",
      setup: func(mocked *roleMocked, arg any, want any) {
        a := arg.(*args)

        permIds := sharedUtil.CastSlice(a.permissionsDTO.PermissionIds, sharedUtil.ToAny[types.Id])

        mocked.Role.EXPECT().
          AddPermissions(mock.Anything, a.permissionsDTO.RoleId, permIds...).
          Return(sql.ErrNoRows)
      },
      args: args{
        ctx: context.Background(),
        permissionsDTO: &dto.ModifyRolesPermissionsDTO{
          RoleId:        types.MustCreateId(),
          PermissionIds: sharedUtil.GenerateMultiple(3, types.MustCreateId),
        },
      },
      want: status.FromRepositoryExist(sql.ErrNoRows),
    },
    {
      name: "Failed to add role's permissions",
      setup: func(mocked *roleMocked, arg any, want any) {
        a := arg.(*args)

        permIds := sharedUtil.CastSlice(a.permissionsDTO.PermissionIds, sharedUtil.ToAny[types.Id])

        mocked.Role.EXPECT().
          AddPermissions(mock.Anything, a.permissionsDTO.RoleId, permIds...).
          Return(dummyErr)
      },
      args: args{
        ctx: context.Background(),
        permissionsDTO: &dto.ModifyRolesPermissionsDTO{
          RoleId:        types.MustCreateId(),
          PermissionIds: sharedUtil.GenerateMultiple(3, types.MustCreateId),
        },
      },
      want: status.FromRepository(dummyErr, status.NullCode),
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newRoleMocked(t)
      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, nil)
      }

      r := &roleService{
        roleRepo: mocked.Role,
        tracer:   mocked.Tracer,
      }
      if got := r.AddPermissions(tt.args.ctx, tt.args.permissionsDTO); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("AddPermissions() = %v, want %v", got, tt.want)
      }
    })
  }
}

func Test_roleService_AddUsers(t *testing.T) {
  type args struct {
    ctx      context.Context
    usersDTO *dto.ModifyUserRolesDTO
  }
  tests := []struct {
    name  string
    setup setupRoleTestFunc
    args  args
    want  status.Object
  }{
    {
      name: "Success to add user's roles",
      setup: func(mocked *roleMocked, arg any, want any) {
        a := arg.(*args)

        roleIds := sharedUtil.CastSlice(a.usersDTO.RoleIds, sharedUtil.ToAny[types.Id])
        mocked.Role.EXPECT().
          AddUser(mock.Anything, a.usersDTO.UserId, roleIds...).
          Return(nil)

      },
      args: args{
        ctx: context.Background(),
        usersDTO: &dto.ModifyUserRolesDTO{
          UserId:  types.MustCreateId(),
          RoleIds: sharedUtil.GenerateMultiple(2, types.MustCreateId),
        },
      },
      want: status.Created(),
    },
    {
      name: "User already has the roles",
      setup: func(mocked *roleMocked, arg any, want any) {
        a := arg.(*args)

        roleIds := sharedUtil.CastSlice(a.usersDTO.RoleIds, sharedUtil.ToAny[types.Id])
        mocked.Role.EXPECT().
          AddUser(mock.Anything, a.usersDTO.UserId, roleIds...).
          Return(sql.ErrNoRows)

      },
      args: args{
        ctx: context.Background(),
        usersDTO: &dto.ModifyUserRolesDTO{
          UserId:  types.MustCreateId(),
          RoleIds: sharedUtil.GenerateMultiple(2, types.MustCreateId),
        },
      },
      want: status.FromRepositoryExist(sql.ErrNoRows),
    },
    {
      name: "Failed to add user's roles",
      setup: func(mocked *roleMocked, arg any, want any) {
        a := arg.(*args)

        roleIds := sharedUtil.CastSlice(a.usersDTO.RoleIds, sharedUtil.ToAny[types.Id])
        mocked.Role.EXPECT().
          AddUser(mock.Anything, a.usersDTO.UserId, roleIds...).
          Return(dummyErr)

      },
      args: args{
        ctx: context.Background(),
        usersDTO: &dto.ModifyUserRolesDTO{
          UserId:  types.MustCreateId(),
          RoleIds: sharedUtil.GenerateMultiple(2, types.MustCreateId),
        },
      },
      want: status.FromRepositoryExist(dummyErr),
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newRoleMocked(t)
      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, nil)
      }

      r := &roleService{
        roleRepo: mocked.Role,
        tracer:   mocked.Tracer,
      }

      if got := r.AddUsers(tt.args.ctx, tt.args.usersDTO); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("AddUsers() = %v, want %v", got, tt.want)
      }
    })
  }
}

func Test_roleService_AppendSuperRolesPermission(t *testing.T) {
  type args struct {
    ctx     context.Context
    permIds []types.Id
  }
  tests := []struct {
    name  string
    setup setupRoleTestFunc
    args  args
    want  status.Object
  }{
    {
      name: "Success to append super roles permissions",
      setup: func(mocked *roleMocked, arg any, want any) {
        a := arg.(*args)

        role := entity.Role{
          Id:          types.MustCreateId(),
          Name:        gofakeit.AnimalType(),
          Description: gofakeit.LoremIpsumSentence(6),
        }

        mocked.Role.EXPECT().
          FindByName(mock.Anything, constant.SUPER_ROLE_NAME).
          Return(role, nil)

        permIds := sharedUtil.CastSlice(a.permIds, sharedUtil.ToAny[types.Id])

        mocked.Role.EXPECT().
          AddPermissions(mock.Anything, role.Id, permIds...).
          Return(nil)
      },
      args: args{
        ctx:     context.Background(),
        permIds: sharedUtil.GenerateMultiple(3, types.MustCreateId),
      },
      want: status.Created(),
    },
    {
      name: "Role not found",
      setup: func(mocked *roleMocked, arg any, want any) {

        mocked.Role.EXPECT().
          FindByName(mock.Anything, constant.SUPER_ROLE_NAME).
          Return(entity.Role{}, sql.ErrNoRows)
      },
      args: args{
        ctx:     context.Background(),
        permIds: sharedUtil.GenerateMultiple(3, types.MustCreateId),
      },
      want: status.ErrInternal(errs.ErrDefaultRoleNotSeeded),
    },
    {
      name: "Permission role failed to get role",
      setup: func(mocked *roleMocked, arg any, want any) {

        mocked.Role.EXPECT().
          FindByName(mock.Anything, constant.SUPER_ROLE_NAME).
          Return(entity.Role{}, dummyErr)
      },
      args: args{
        ctx:     context.Background(),
        permIds: sharedUtil.GenerateMultiple(3, types.MustCreateId),
      },
      want: status.New(status.REPOSITORY_ERROR, dummyErr),
    },
    {
      name: "Super role already has the permissions",
      setup: func(mocked *roleMocked, arg any, want any) {
        a := arg.(*args)

        role := entity.Role{
          Id:          types.MustCreateId(),
          Name:        gofakeit.AnimalType(),
          Description: gofakeit.LoremIpsumSentence(6),
        }

        mocked.Role.EXPECT().
          FindByName(mock.Anything, constant.SUPER_ROLE_NAME).
          Return(role, nil)

        permIds := sharedUtil.CastSlice(a.permIds, sharedUtil.ToAny[types.Id])

        mocked.Role.EXPECT().
          AddPermissions(mock.Anything, role.Id, permIds...).
          Return(sql.ErrNoRows)
      },
      args: args{
        ctx:     context.Background(),
        permIds: sharedUtil.GenerateMultiple(3, types.MustCreateId),
      },
      want: status.FromRepositoryExist(sql.ErrNoRows),
    },
    {
      name: "Failed to add super roles permissions",
      setup: func(mocked *roleMocked, arg any, want any) {
        a := arg.(*args)

        role := entity.Role{
          Id:          types.MustCreateId(),
          Name:        gofakeit.AnimalType(),
          Description: gofakeit.LoremIpsumSentence(6),
        }

        mocked.Role.EXPECT().
          FindByName(mock.Anything, constant.SUPER_ROLE_NAME).
          Return(role, nil)

        permIds := sharedUtil.CastSlice(a.permIds, sharedUtil.ToAny[types.Id])

        mocked.Role.EXPECT().
          AddPermissions(mock.Anything, role.Id, permIds...).
          Return(dummyErr)
      },
      args: args{
        ctx:     context.Background(),
        permIds: sharedUtil.GenerateMultiple(3, types.MustCreateId),
      },
      want: status.FromRepositoryExist(dummyErr),
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newRoleMocked(t)
      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, nil)
      }

      r := &roleService{
        roleRepo: mocked.Role,
        tracer:   mocked.Tracer,
      }

      if got := r.AppendSuperRolesPermission(tt.args.ctx, tt.args.permIds...); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("AppendSuperRolesPermission() = %v, want %v", got, tt.want)
      }
    })
  }
}

func Test_roleService_Create(t *testing.T) {
  type args struct {
    ctx       context.Context
    createDTO *dto.RoleCreateDTO
  }
  tests := []struct {
    name  string
    setup setupRoleTestFunc
    args  args
    want1 status.Object
  }{
    {
      name: "Success create roles",
      setup: func(mocked *roleMocked, arg any, want any) {
        mocked.Role.EXPECT().
          Create(mock.Anything, mock.Anything).
          Return(nil)
      },
      args: args{
        ctx: context.Background(),
        createDTO: &dto.RoleCreateDTO{
          Name:        gofakeit.AnimalType(),
          Description: types.SomeNullable(gofakeit.LoremIpsumSentence(4)),
        },
      },
      want1: status.Created(),
    },
    {
      name: "Roles already exist",
      setup: func(mocked *roleMocked, arg any, want any) {
        mocked.Role.EXPECT().
          Create(mock.Anything, mock.Anything).
          Return(sql.ErrNoRows)
      },
      args: args{
        ctx: context.Background(),
        createDTO: &dto.RoleCreateDTO{
          Name:        gofakeit.AnimalType(),
          Description: types.SomeNullable(gofakeit.LoremIpsumSentence(4)),
        },
      },
      want1: status.FromRepositoryExist(sql.ErrNoRows),
    },
    {
      name: "Failed to create roles",
      setup: func(mocked *roleMocked, arg any, want any) {
        mocked.Role.EXPECT().
          Create(mock.Anything, mock.Anything).
          Return(dummyErr)
      },
      args: args{
        ctx: context.Background(),
        createDTO: &dto.RoleCreateDTO{
          Name:        gofakeit.AnimalType(),
          Description: types.SomeNullable(gofakeit.LoremIpsumSentence(4)),
        },
      },
      want1: status.FromRepositoryExist(dummyErr),
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newRoleMocked(t)
      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, nil)
      }

      r := &roleService{
        roleRepo: mocked.Role,
        tracer:   mocked.Tracer,
      }

      got, got1 := r.Create(tt.args.ctx, tt.args.createDTO)
      if !reflect.DeepEqual(got1, tt.want1) {
        t.Errorf("Create() got1 = %v, want %v", got1, tt.want1)
      }

      if got1.IsError() {
        return
      }

      require.False(t, got.Eq(types.NullId()))
    })
  }
}

func Test_roleService_Delete(t *testing.T) {
  type args struct {
    ctx    context.Context
    roleId types.Id
  }
  tests := []struct {
    name  string
    setup setupRoleTestFunc
    args  args
    want  status.Object
  }{
    {
      name: "Success to delete roles",
      setup: func(mocked *roleMocked, arg any, want any) {
        a := arg.(*args)

        mocked.Role.EXPECT().
          Delete(mock.Anything, a.roleId).
          Return(nil)
      },
      args: args{
        ctx:    context.Background(),
        roleId: types.MustCreateId(),
      },
      want: status.Deleted(),
    },
    {
      name: "Failed to delete roles",
      setup: func(mocked *roleMocked, arg any, want any) {
        a := arg.(*args)

        mocked.Role.EXPECT().
          Delete(mock.Anything, a.roleId).
          Return(sql.ErrNoRows)
      },
      args: args{
        ctx:    context.Background(),
        roleId: types.MustCreateId(),
      },
      want: status.FromRepository(sql.ErrNoRows, status.NullCode),
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newRoleMocked(t)
      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, nil)
      }

      r := &roleService{
        roleRepo: mocked.Role,
        tracer:   mocked.Tracer,
      }

      if got := r.Delete(tt.args.ctx, tt.args.roleId); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("Delete() = %v, want %v", got, tt.want)
      }
    })
  }
}

func Test_roleService_GetAll(t *testing.T) {
  type args struct {
    ctx      context.Context
    pagedDto *sharedDto.PagedElementDTO
  }
  tests := []struct {
    name  string
    setup setupRoleTestFunc
    args  args
    want  sharedDto.PagedElementResult[dto.RoleResponseDTO]
    want1 status.Object
  }{
    {
      name: "Success get permissions",
      setup: func(mocked *roleMocked, arg any, want any) {
        a := arg.(*args)
        w := want.(*sharedDto.PagedElementResult[dto.RoleResponseDTO])

        mocked.Role.EXPECT().
          Get(mock.Anything, a.pagedDto.ToQueryParam()).
          Return(repo.PaginatedResult[entity.Role]{nil, w.TotalElements, w.Element}, nil)

      },
      args: args{
        ctx: context.Background(),
        pagedDto: &sharedDto.PagedElementDTO{
          Element: 2,
          Page:    2,
        },
      },
      want: sharedDto.PagedElementResult[dto.RoleResponseDTO]{
        Data:          nil,
        Element:       0,
        Page:          2,
        TotalElements: 20,
        TotalPages:    10,
      },
      want1: status.Success(),
    },
    {
      name: "Failed to get users",
      setup: func(mocked *roleMocked, arg any, want any) {
        a := arg.(*args)
        //w := want.(*sharedDto.PagedElementResult[dto.UserResponseDTO])

        mocked.Role.EXPECT().
          Get(mock.Anything, a.pagedDto.ToQueryParam()).
          Return(repo.PaginatedResult[entity.Role]{nil, 0, 0}, sql.ErrNoRows)

      },
      args: args{
        ctx: context.Background(),
        pagedDto: &sharedDto.PagedElementDTO{
          Element: 2,
          Page:    2,
        },
      },
      want: sharedDto.PagedElementResult[dto.RoleResponseDTO]{
        Data:          nil,
        Element:       0,
        Page:          0,
        TotalElements: 0,
        TotalPages:    0,
      },
      want1: status.FromRepository(sql.ErrNoRows, status.NullCode),
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newRoleMocked(t)
      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, &tt.want)
      }

      r := &roleService{
        roleRepo: mocked.Role,
        tracer:   mocked.Tracer,
      }

      got, got1 := r.GetAll(tt.args.ctx, tt.args.pagedDto)
      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("GetAll() got = %v, want %v", got, tt.want)
      }
      if !reflect.DeepEqual(got1, tt.want1) {
        t.Errorf("GetAll() got1 = %v, want %v", got1, tt.want1)
      }
    })
  }
}

func Test_roleService_FindByIds(t *testing.T) {
  type args struct {
    ctx     context.Context
    roleIds []types.Id
  }
  tests := []struct {
    name    string
    setup   setupRoleTestFunc
    args    args
    wantLen int
    want1   status.Object
  }{
    {
      name: "Success get permissions by id",
      setup: func(mocked *roleMocked, arg any, want any) {
        a := arg.(*args)

        roleIds := sharedUtil.CastSlice(a.roleIds, sharedUtil.ToAny[types.Id])

        roles := sharedUtil.CastSlice(a.roleIds, func(from types.Id) entity.Role {
          return entity.Role{
            Id:          from,
            Name:        gofakeit.AnimalType(),
            Description: gofakeit.LoremIpsumSentence(3),
          }
        })

        mocked.Role.EXPECT().
          FindByIds(mock.Anything, roleIds...).
          Return(roles, nil)
      },
      args: args{
        ctx:     context.Background(),
        roleIds: sharedUtil.GenerateMultiple(2, types.MustCreateId),
      },
      wantLen: 2,
      want1:   status.Success(),
    },
    {
      name: "Failed to get permissions by id",
      setup: func(mocked *roleMocked, arg any, want any) {
        a := arg.(*args)

        roleIds := sharedUtil.CastSlice(a.roleIds, sharedUtil.ToAny[types.Id])

        mocked.Role.EXPECT().
          FindByIds(mock.Anything, roleIds...).
          Return(nil, sql.ErrNoRows)
      },
      args: args{
        ctx:     context.Background(),
        roleIds: sharedUtil.GenerateMultiple(2, types.MustCreateId),
      },
      wantLen: 0,
      want1:   status.FromRepository(sql.ErrNoRows, status.NullCode),
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newRoleMocked(t)
      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, nil)
      }

      r := &roleService{
        roleRepo: mocked.Role,
        tracer:   mocked.Tracer,
      }

      got, got1 := r.FindByIds(tt.args.ctx, tt.args.roleIds...)
      require.Equal(t, len(got), tt.wantLen)
      if !reflect.DeepEqual(got1, tt.want1) {
        t.Errorf("FindByIds() got1 = %v, want %v", got1, tt.want1)
      }
    })
  }
}

func Test_roleService_FindByUserId(t *testing.T) {
  type args struct {
    ctx    context.Context
    userId types.Id
  }
  tests := []struct {
    name  string
    setup setupRoleTestFunc
    args  args
    want1 status.Object
  }{
    {
      name: "Success get user permissions",
      setup: func(mocked *roleMocked, arg any, want any) {
        a := arg.(*args)

        roles := sharedUtil.GenerateMultiple(2, func() entity.Role {
          return entity.Role{
            Id:          types.MustCreateId(),
            Name:        gofakeit.AppName(),
            Description: gofakeit.LoremIpsumSentence(3),
          }
        })

        mocked.Role.EXPECT().
          FindByUserId(mock.Anything, a.userId).
          Return(roles, nil)
      },
      args: args{
        ctx:    context.Background(),
        userId: types.MustCreateId(),
      },
      want1: status.Success(),
    },
    {
      name: "Failed to get user permissions",
      setup: func(mocked *roleMocked, arg any, want any) {
        a := arg.(*args)

        mocked.Role.EXPECT().
          FindByUserId(mock.Anything, a.userId).
          Return(nil, sql.ErrNoRows)
      },
      args: args{
        ctx:    context.Background(),
        userId: types.MustCreateId(),
      },
      want1: status.FromRepository(sql.ErrNoRows, status.NullCode),
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newRoleMocked(t)
      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, nil)
      }

      r := &roleService{
        roleRepo: mocked.Role,
        tracer:   mocked.Tracer,
      }

      _, got1 := r.FindByUserId(tt.args.ctx, tt.args.userId)
      if !reflect.DeepEqual(got1, tt.want1) {
        t.Errorf("FindByUserId() got1 = %v, want %v", got1, tt.want1)
      }
    })
  }
}

func Test_roleService_RemovePermissions(t *testing.T) {
  type args struct {
    ctx            context.Context
    permissionsDTO *dto.ModifyRolesPermissionsDTO
  }
  tests := []struct {
    name  string
    setup setupRoleTestFunc
    args  args
    want  status.Object
  }{
    {
      name: "Success to remove role's permissions",
      setup: func(mocked *roleMocked, arg any, want any) {
        a := arg.(*args)

        permIds := sharedUtil.CastSlice(a.permissionsDTO.PermissionIds, sharedUtil.ToAny[types.Id])

        mocked.Role.EXPECT().
          RemovePermissions(mock.Anything, a.permissionsDTO.RoleId, permIds...).
          Return(nil)

      },
      args: args{
        ctx: context.Background(),
        permissionsDTO: &dto.ModifyRolesPermissionsDTO{
          RoleId:        types.MustCreateId(),
          PermissionIds: sharedUtil.GenerateMultiple(2, types.MustCreateId),
        },
      },
      want: status.Deleted(),
    },
    {
      name: "Failed to remove role's permissions",
      setup: func(mocked *roleMocked, arg any, want any) {
        a := arg.(*args)

        permIds := sharedUtil.CastSlice(a.permissionsDTO.PermissionIds, sharedUtil.ToAny[types.Id])

        mocked.Role.EXPECT().
          RemovePermissions(mock.Anything, a.permissionsDTO.RoleId, permIds...).
          Return(sql.ErrNoRows)

      },
      args: args{
        ctx: context.Background(),
        permissionsDTO: &dto.ModifyRolesPermissionsDTO{
          RoleId:        types.MustCreateId(),
          PermissionIds: sharedUtil.GenerateMultiple(2, types.MustCreateId),
        },
      },
      want: status.FromRepository(sql.ErrNoRows, status.NullCode),
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newRoleMocked(t)
      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, nil)
      }

      r := &roleService{
        roleRepo: mocked.Role,
        tracer:   mocked.Tracer,
      }

      if got := r.RemovePermissions(tt.args.ctx, tt.args.permissionsDTO); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("RemovePermissions() = %v, want %v", got, tt.want)
      }
    })
  }
}

func Test_roleService_RemoveUsers(t *testing.T) {
  type args struct {
    ctx      context.Context
    usersDTO *dto.ModifyUserRolesDTO
  }
  tests := []struct {
    name  string
    setup setupRoleTestFunc
    args  args
    want  status.Object
  }{
    {
      name: "Success to remove user's roles",
      setup: func(mocked *roleMocked, arg any, want any) {
        a := arg.(*args)

        roleIds := sharedUtil.CastSlice(a.usersDTO.RoleIds, sharedUtil.ToAny[types.Id])

        mocked.Role.EXPECT().
          RemoveUser(mock.Anything, a.usersDTO.UserId, roleIds...).
          Return(nil)
      },
      args: args{
        ctx: context.Background(),
        usersDTO: &dto.ModifyUserRolesDTO{
          UserId:  types.Id{},
          RoleIds: sharedUtil.GenerateMultiple(3, types.MustCreateId),
        },
      },
      want: status.Deleted(),
    },
    {
      name: "Failed to remove user's roles",
      setup: func(mocked *roleMocked, arg any, want any) {
        a := arg.(*args)

        roleIds := sharedUtil.CastSlice(a.usersDTO.RoleIds, sharedUtil.ToAny[types.Id])

        mocked.Role.EXPECT().
          RemoveUser(mock.Anything, a.usersDTO.UserId, roleIds...).
          Return(sql.ErrNoRows)
      },
      args: args{
        ctx: context.Background(),
        usersDTO: &dto.ModifyUserRolesDTO{
          UserId:  types.Id{},
          RoleIds: sharedUtil.GenerateMultiple(3, types.MustCreateId),
        },
      },
      want: status.FromRepository(sql.ErrNoRows, status.NullCode),
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newRoleMocked(t)
      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, nil)
      }

      r := &roleService{
        roleRepo: mocked.Role,
        tracer:   mocked.Tracer,
      }

      if got := r.RemoveUsers(tt.args.ctx, tt.args.usersDTO); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("RemoveUsers() = %v, want %v", got, tt.want)
      }
    })
  }
}

func Test_roleService_Update(t *testing.T) {
  type args struct {
    ctx       context.Context
    updateDTO *dto.RoleUpdateDTO
  }
  tests := []struct {
    name  string
    setup setupRoleTestFunc
    args  args
    want  status.Object
  }{
    {
      name: "Success update roles",
      setup: func(mocked *roleMocked, arg any, want any) {
        mocked.Role.EXPECT().
          Patch(mock.Anything, mock.Anything).
          Return(nil)
      },
      args: args{
        ctx: context.Background(),
        updateDTO: &dto.RoleUpdateDTO{
          RoleId:      types.MustCreateId(),
          Name:        types.SomeNullable(gofakeit.AnimalType()),
          Description: types.SomeNullable(gofakeit.LoremIpsumSentence(4)),
        },
      },
      want: status.Updated(),
    },
    {
      name: "Failed to update roles",
      setup: func(mocked *roleMocked, arg any, want any) {
        mocked.Role.EXPECT().
          Patch(mock.Anything, mock.Anything).
          Return(sql.ErrNoRows)
      },
      args: args{
        ctx: context.Background(),
        updateDTO: &dto.RoleUpdateDTO{
          RoleId:      types.MustCreateId(),
          Name:        types.SomeNullable(gofakeit.AnimalType()),
          Description: types.SomeNullable(gofakeit.LoremIpsumSentence(4)),
        },
      },
      want: status.FromRepository(sql.ErrNoRows, status.NullCode),
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newRoleMocked(t)
      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, nil)
      }

      r := &roleService{
        roleRepo: mocked.Role,
        tracer:   mocked.Tracer,
      }

      if got := r.Update(tt.args.ctx, tt.args.updateDTO); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("Update() = %v, want %v", got, tt.want)
      }
    })
  }
}
