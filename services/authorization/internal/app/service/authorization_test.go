package service

import (
  "context"
  "database/sql"
  "github.com/brianvoe/gofakeit/v7"
  "github.com/stretchr/testify/mock"
  "go.opentelemetry.io/otel/trace"
  "go.opentelemetry.io/otel/trace/noop"
  "nexa/services/authorization/internal/domain/dto"
  "nexa/services/authorization/internal/domain/entity"
  repoMock "nexa/services/authorization/internal/domain/repository/mocks"
  sharedErr "nexa/shared/errors"
  "nexa/shared/status"
  "nexa/shared/types"
  "nexa/shared/util"
  "reflect"
  "testing"
  "time"
)

func newAuthZMock(t *testing.T) authZMocked {
  // Tracer
  provider := noop.NewTracerProvider()

  return authZMocked{
    Role:   repoMock.NewRoleMock(t),
    Tracer: provider.Tracer("MOCK"),
  }
}

type authZMocked struct {
  Role   *repoMock.RoleMock
  Tracer trace.Tracer
}

type setupAuthZMock func(mocked *authZMocked, arg any)

func Test_authorizationService_IsAuthorized(t *testing.T) {
  type args struct {
    ctx     context.Context
    authDto *dto.IsAuthorizationDTO
  }
  tests := []struct {
    name  string
    setup setupAuthZMock
    args  args
    want  status.Object
  }{
    {
      name: "User authorized",
      setup: func(mocked *authZMocked, arg any) {
        a := arg.(*args)

        roles := util.GenerateMultiple(3, func() entity.Role {
          return entity.Role{
            Id:          types.MustCreateId(),
            Name:        gofakeit.AnimalType(),
            Description: gofakeit.LoremIpsumSentence(5),
            Permissions: []entity.Permission{
              {
                Id:        types.MustCreateId(),
                Resource:  "should",
                Action:    "work",
                CreatedAt: time.Now(),
              },
              {
                Id:        types.MustCreateId(),
                Resource:  "hooh",
                Action:    "work",
                CreatedAt: time.Now(),
              },
            },
          }
        })

        mocked.Role.EXPECT().
          FindByUserId(mock.Anything, a.authDto.UserId).
          Return(roles, nil)
      },
      args: args{
        ctx: context.Background(),
        authDto: &dto.IsAuthorizationDTO{
          UserId:             types.MustCreateId(),
          ExpectedPermission: "should:work",
        },
      },
      want: status.Success(),
    },
    {
      name: "User unauthorized",
      setup: func(mocked *authZMocked, arg any) {
        a := arg.(*args)

        roles := util.GenerateMultiple(3, func() entity.Role {
          return entity.Role{
            Id:          types.MustCreateId(),
            Name:        gofakeit.AnimalType(),
            Description: gofakeit.LoremIpsumSentence(5),
            Permissions: []entity.Permission{
              {
                Id:        types.MustCreateId(),
                Resource:  "should",
                Action:    "work",
                CreatedAt: time.Now(),
              },
              {
                Id:        types.MustCreateId(),
                Resource:  "not",
                Action:    "work",
                CreatedAt: time.Now(),
              },
            },
          }
        })

        mocked.Role.EXPECT().
          FindByUserId(mock.Anything, a.authDto.UserId).
          Return(roles, nil)
      },
      args: args{
        ctx: context.Background(),
        authDto: &dto.IsAuthorizationDTO{
          UserId:             types.MustCreateId(),
          ExpectedPermission: "hooh:work",
        },
      },
      want: status.ErrUnAuthorized(sharedErr.ErrUnauthorizedPermission),
    },
    {
      name: "User doesn't have roles",
      setup: func(mocked *authZMocked, arg any) {
        a := arg.(*args)

        mocked.Role.EXPECT().
          FindByUserId(mock.Anything, a.authDto.UserId).
          Return(nil, sql.ErrNoRows)
      },
      args: args{
        ctx: context.Background(),
        authDto: &dto.IsAuthorizationDTO{
          UserId:             types.MustCreateId(),
          ExpectedPermission: "hooh:work",
        },
      },
      want: status.ErrUnAuthorized(sharedErr.ErrUnauthorized),
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newAuthZMock(t)

      if tt.setup != nil {
        tt.setup(&mocked, &tt.args)
      }

      a := &authorizationService{
        roleRepo: mocked.Role,
        tracer:   mocked.Tracer,
      }

      if got := a.IsAuthorized(tt.args.ctx, tt.args.authDto); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("IsAuthorized() = %v, want %v", got, tt.want)
      }
    })
  }
}
