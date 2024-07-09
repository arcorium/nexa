package service

import (
  "context"
  "database/sql"
  "errors"
  "github.com/brianvoe/gofakeit/v7"
  "github.com/stretchr/testify/mock"
  "github.com/stretchr/testify/require"
  "go.opentelemetry.io/otel/trace"
  "go.opentelemetry.io/otel/trace/noop"
  "nexa/services/authorization/internal/domain/dto"
  "nexa/services/authorization/internal/domain/entity"
  repoMock "nexa/services/authorization/internal/domain/repository/mocks"
  sharedDto "nexa/shared/dto"
  "nexa/shared/status"
  "nexa/shared/types"
  "nexa/shared/util"
  "nexa/shared/util/repo"
  "reflect"
  "testing"
  "time"
)

var dummyErr = errors.New("dummy error")

var dummyId = types.MustCreateId()

func newPermsMocked(t *testing.T) permsMocked {
  // Tracer
  provider := noop.NewTracerProvider()
  return permsMocked{
    Perm:   repoMock.NewPermissionMock(t),
    Tracer: provider.Tracer("MOCK"),
  }
}

type permsMocked struct {
  Perm   *repoMock.PermissionMock
  Tracer trace.Tracer
}

type setupPermsTestFunc func(mocked *permsMocked, arg any, want any)

func Test_permissionService_Create(t *testing.T) {
  type args struct {
    ctx       context.Context
    createDTO *dto.PermissionCreateDTO
  }
  tests := []struct {
    name  string
    setup setupPermsTestFunc
    args  args
    want1 status.Object
  }{
    {
      name: "Success create permission",
      setup: func(mocked *permsMocked, arg any, want any) {
        //a := arg.(*args)

        mocked.Perm.EXPECT().
          Create(mock.Anything, mock.Anything).
          Return(nil)
      },
      args: args{
        ctx: context.Background(),
        createDTO: &dto.PermissionCreateDTO{
          Resource: gofakeit.AppName(),
          Action:   gofakeit.Username(),
        },
      },
      want1: status.Created(),
    },
    {
      name: "Permission already exists",
      setup: func(mocked *permsMocked, arg any, want any) {
        mocked.Perm.EXPECT().
          Create(mock.Anything, mock.Anything).
          Return(sql.ErrNoRows)
      },
      args: args{
        ctx: context.Background(),
        createDTO: &dto.PermissionCreateDTO{
          Resource: gofakeit.AppName(),
          Action:   gofakeit.Username(),
        },
      },
      want1: status.FromRepositoryExist(sql.ErrNoRows),
    },
    {
      name: "Failed to create permission",
      setup: func(mocked *permsMocked, arg any, want any) {
        mocked.Perm.EXPECT().
          Create(mock.Anything, mock.Anything).
          Return(dummyErr)
      },
      args: args{
        ctx: context.Background(),
        createDTO: &dto.PermissionCreateDTO{
          Resource: gofakeit.AppName(),
          Action:   gofakeit.Username(),
        },
      },
      want1: status.FromRepository(dummyErr, status.NullCode),
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newPermsMocked(t)
      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, nil)
      }

      p := &permissionService{
        permRepo: mocked.Perm,
        tracer:   mocked.Tracer,
      }

      got, got1 := p.Create(tt.args.ctx, tt.args.createDTO)
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

func Test_permissionService_Delete(t *testing.T) {
  type args struct {
    ctx    context.Context
    permId types.Id
  }
  tests := []struct {
    name  string
    setup setupPermsTestFunc
    args  args
    want  status.Object
  }{
    {
      name: "Success delete permission",
      setup: func(mocked *permsMocked, arg any, want any) {
        a := arg.(*args)

        mocked.Perm.EXPECT().
          Delete(mock.Anything, a.permId).
          Return(nil)
      },
      args: args{
        ctx:    context.Background(),
        permId: types.MustCreateId(),
      },
      want: status.Deleted(),
    },
    {
      name: "Failed to delete permission",
      setup: func(mocked *permsMocked, arg any, want any) {
        a := arg.(*args)

        mocked.Perm.EXPECT().
          Delete(mock.Anything, a.permId).
          Return(sql.ErrNoRows)
      },
      args: args{
        ctx:    context.Background(),
        permId: types.MustCreateId(),
      },
      want: status.FromRepository(sql.ErrNoRows, status.NullCode),
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newPermsMocked(t)
      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, nil)
      }

      p := &permissionService{
        permRepo: mocked.Perm,
        tracer:   mocked.Tracer,
      }

      if got := p.Delete(tt.args.ctx, tt.args.permId); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("Delete() = %v, want %v", got, tt.want)
      }
    })
  }
}

func Test_permissionService_Find(t *testing.T) {
  type args struct {
    ctx     context.Context
    permIds []types.Id
  }
  tests := []struct {
    name    string
    setup   setupPermsTestFunc
    args    args
    wantLen int
    want1   status.Object
  }{
    {
      name: "Success find permissions",
      setup: func(mocked *permsMocked, arg any, want any) {
        a := arg.(*args)

        permIds := util.CastSlice(a.permIds, util.ToAny[types.Id])

        perms := util.CastSlice(a.permIds, func(permId types.Id) entity.Permission {
          return entity.Permission{
            Id:        permId,
            Resource:  gofakeit.AppName(),
            Action:    gofakeit.AnimalType(),
            CreatedAt: time.Now(),
          }
        })

        mocked.Perm.EXPECT().
          FindByIds(mock.Anything, permIds...).
          Return(perms, nil)

      },
      args: args{
        ctx:     context.Background(),
        permIds: util.GenerateMultiple(2, types.MustCreateId),
      },
      wantLen: 2,
      want1:   status.Success(),
    },
    {
      name: "Failed to find permissions",
      setup: func(mocked *permsMocked, arg any, want any) {
        a := arg.(*args)

        permIds := util.CastSlice(a.permIds, util.ToAny[types.Id])

        mocked.Perm.EXPECT().
          FindByIds(mock.Anything, permIds...).
          Return(nil, sql.ErrNoRows)
      },
      args: args{
        ctx:     context.Background(),
        permIds: util.GenerateMultiple(2, types.MustCreateId),
      },
      wantLen: 0,
      want1:   status.FromRepository(sql.ErrNoRows, status.NullCode),
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newPermsMocked(t)
      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, nil)
      }

      p := &permissionService{
        permRepo: mocked.Perm,
        tracer:   mocked.Tracer,
      }

      got, got1 := p.Find(tt.args.ctx, tt.args.permIds...)
      require.Equal(t, tt.wantLen, len(got))
      if !reflect.DeepEqual(got1, tt.want1) {
        t.Errorf("Find() got1 = %v, want %v", got1, tt.want1)
      }
    })
  }
}

func Test_permissionService_GetAll(t *testing.T) {
  type args struct {
    ctx      context.Context
    pagedDto *sharedDto.PagedElementDTO
  }
  tests := []struct {
    name  string
    setup setupPermsTestFunc
    args  args
    want  sharedDto.PagedElementResult[dto.PermissionResponseDTO]
    want1 status.Object
  }{
    {
      name: "Success get permissions",
      setup: func(mocked *permsMocked, arg any, want any) {
        a := arg.(*args)
        w := want.(*sharedDto.PagedElementResult[dto.PermissionResponseDTO])

        mocked.Perm.EXPECT().
          Get(mock.Anything, a.pagedDto.ToQueryParam()).
          Return(repo.PaginatedResult[entity.Permission]{nil, w.TotalElements, w.Element}, nil)

      },
      args: args{
        ctx: context.Background(),
        pagedDto: &sharedDto.PagedElementDTO{
          Element: 2,
          Page:    2,
        },
      },
      want: sharedDto.PagedElementResult[dto.PermissionResponseDTO]{
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
      setup: func(mocked *permsMocked, arg any, want any) {
        a := arg.(*args)
        //w := want.(*sharedDto.PagedElementResult[dto.UserResponseDTO])

        mocked.Perm.EXPECT().
          Get(mock.Anything, a.pagedDto.ToQueryParam()).
          Return(repo.PaginatedResult[entity.Permission]{nil, 0, 0}, sql.ErrNoRows)

      },
      args: args{
        ctx: context.Background(),
        pagedDto: &sharedDto.PagedElementDTO{
          Element: 2,
          Page:    2,
        },
      },
      want: sharedDto.PagedElementResult[dto.PermissionResponseDTO]{
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
      mocked := newPermsMocked(t)
      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, &tt.want)
      }

      p := &permissionService{
        permRepo: mocked.Perm,
        tracer:   mocked.Tracer,
      }

      got, got1 := p.GetAll(tt.args.ctx, tt.args.pagedDto)
      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("GetAll() got = %v, want %v", got, tt.want)
      }
      if !reflect.DeepEqual(got1, tt.want1) {
        t.Errorf("GetAll() got1 = %v, want %v", got1, tt.want1)
      }
    })
  }
}

func Test_permissionService_FindByRoles(t *testing.T) {
  type args struct {
    ctx     context.Context
    roleIds []types.Id
  }
  tests := []struct {
    name    string
    setup   setupPermsTestFunc
    args    args
    wantLen int
    want1   status.Object
  }{
    {
      name: "Success find role's permissions",
      setup: func(mocked *permsMocked, arg any, want any) {
        a := arg.(*args)

        roleIds := util.CastSlice(a.roleIds, util.ToAny[types.Id])
        perms := util.GenerateMultiple(3, func() entity.Permission {
          return entity.Permission{
            Id:        types.MustCreateId(),
            Resource:  gofakeit.AppName(),
            Action:    gofakeit.AnimalType(),
            CreatedAt: time.Now(),
          }
        })

        mocked.Perm.EXPECT().
          FindByRoleIds(mock.Anything, roleIds...).
          Return(perms, nil)
      },
      args: args{
        ctx:     context.Background(),
        roleIds: util.GenerateMultiple(2, types.MustCreateId),
      },
      wantLen: 3,
      want1:   status.Success(),
    },

    {
      name: "Failed to find role's permissions",
      setup: func(mocked *permsMocked, arg any, want any) {
        a := arg.(*args)

        roleIds := util.CastSlice(a.roleIds, util.ToAny[types.Id])

        mocked.Perm.EXPECT().
          FindByRoleIds(mock.Anything, roleIds...).
          Return(nil, sql.ErrNoRows)
      },
      args: args{
        ctx:     context.Background(),
        roleIds: util.GenerateMultiple(2, types.MustCreateId),
      },
      wantLen: 0,
      want1:   status.FromRepository(sql.ErrNoRows, status.NullCode),
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newPermsMocked(t)
      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, nil)
      }

      p := &permissionService{
        permRepo: mocked.Perm,
        tracer:   mocked.Tracer,
      }

      got, got1 := p.FindByRoles(tt.args.ctx, tt.args.roleIds...)
      require.Equal(t, tt.wantLen, len(got))
      if !reflect.DeepEqual(got1, tt.want1) {
        t.Errorf("FindByRoles() got1 = %v, want %v", got1, tt.want1)
      }
    })
  }
}
