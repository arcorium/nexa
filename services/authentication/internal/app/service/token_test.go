package service

import (
  "context"
  "database/sql"
  "errors"
  "github.com/arcorium/nexa/shared/status"
  "github.com/arcorium/nexa/shared/types"
  "github.com/arcorium/nexa/shared/util"
  "github.com/stretchr/testify/mock"
  "github.com/stretchr/testify/require"
  "go.opentelemetry.io/otel/trace"
  "go.opentelemetry.io/otel/trace/noop"
  "nexa/services/authentication/internal/domain/dto"
  "nexa/services/authentication/internal/domain/entity"
  repoMock "nexa/services/authentication/internal/domain/repository/mocks"
  sharedErr "nexa/services/authentication/util/errors"
  "reflect"
  "testing"
  "time"
)

var dummyErr = errors.New("dummy error")

var claimsUserId = types.MustCreateId()

var duration = time.Hour * 5

func newTokenMocked(t *testing.T) tokenMocked {
  // Tracer
  provider := noop.NewTracerProvider()

  return tokenMocked{
    Config: TokenServiceConfig{
      VerificationTokenExpiration: duration,
      ResetTokenExpiration:        duration,
    },
    Token:  repoMock.NewTokenMock(t),
    Tracer: provider.Tracer("MOCK"),
  }
}

type tokenMocked struct {
  Config TokenServiceConfig
  Token  *repoMock.TokenMock
  Tracer trace.Tracer
}

type setupTokenTestFunc func(mocked *tokenMocked, arg any)

func Test_tokenService_Request(t1 *testing.T) {
  type args struct {
    ctx context.Context
    req *dto.TokenCreateDTO
  }
  tests := []struct {
    name  string
    setup setupTokenTestFunc
    args  args
    want  dto.TokenResponseDTO
    want1 status.Object
  }{
    {
      name: "Success request reset password token",
      setup: func(mocked *tokenMocked, arg any) {
        //a := arg.(*args)
        mocked.Token.EXPECT().
          Create(mock.Anything, mock.Anything).
          Return(nil)
      },
      args: args{
        ctx: context.Background(),
        req: &dto.TokenCreateDTO{
          UserId: types.MustCreateId(),
          Usage:  entity.TokenUsageResetPassword,
        },
      },
      want: dto.TokenResponseDTO{
        Usage:     entity.TokenUsageResetPassword,
        ExpiredAt: time.Now().Add(duration),
      },
      want1: status.Created(),
    },
    {
      name: "Success request verification token",
      setup: func(mocked *tokenMocked, arg any) {
        //a := arg.(*args)
        mocked.Token.EXPECT().
          Create(mock.Anything, mock.Anything).
          Return(nil)
      },
      args: args{
        ctx: context.Background(),
        req: &dto.TokenCreateDTO{
          UserId: types.MustCreateId(),
          Usage:  entity.TokenUsageVerification,
        },
      },
      want: dto.TokenResponseDTO{
        Usage:     entity.TokenUsageVerification,
        ExpiredAt: time.Now().Add(duration),
      },
      want1: status.Created(),
    },
    {
      name: "Failed to request verification token",
      setup: func(mocked *tokenMocked, arg any) {
        //a := arg.(*args)
        mocked.Token.EXPECT().
          Create(mock.Anything, mock.Anything).
          Return(sql.ErrNoRows)
      },
      args: args{
        ctx: context.Background(),
        req: &dto.TokenCreateDTO{
          UserId: types.MustCreateId(),
          Usage:  entity.TokenUsageVerification,
        },
      },
      want:  dto.TokenResponseDTO{},
      want1: status.FromRepository(sql.ErrNoRows, status.NullCode),
    },
  }
  for _, tt := range tests {
    t1.Run(tt.name, func(t1 *testing.T) {
      mocked := newTokenMocked(t1)

      if tt.setup != nil {
        tt.setup(&mocked, &tt.args)
      }

      t := &tokenService{
        config:    mocked.Config,
        tokenRepo: mocked.Token,
        tracer:    mocked.Tracer,
      }

      got, got1 := t.Request(tt.args.ctx, tt.args.req)
      if !reflect.DeepEqual(got1, tt.want1) {
        t1.Errorf("Request() got1 = %v, want %v", got1, tt.want1)
      }
      if got1.IsError() {
        return
      }

      require.NotEmpty(t1, got.Token)
      require.Equal(t1, got.Usage, tt.want.Usage)
    })
  }
}

func Test_tokenService_Verify(t1 *testing.T) {
  type args struct {
    ctx       context.Context
    verifyDTO *dto.TokenVerifyDTO
  }
  tests := []struct {
    name  string
    setup setupTokenTestFunc
    args  args
    want1 status.Object
  }{
    {
      name: "Success verify token",
      setup: func(mocked *tokenMocked, arg any) {
        a := arg.(*args)

        token := entity.Token{
          Token:     a.verifyDTO.Token,
          UserId:    types.MustCreateId(),
          Usage:     a.verifyDTO.Usage,
          ExpiredAt: time.Now().Add(time.Hour * 5),
        }

        mocked.Token.EXPECT().
          Find(mock.Anything, a.verifyDTO.Token).
          Return(token, nil)

        mocked.Token.EXPECT().
          Delete(mock.Anything, a.verifyDTO.Token).
          Return(nil)

      },
      args: args{
        ctx: context.Background(),
        verifyDTO: &dto.TokenVerifyDTO{
          Token: util.RandomString(32),
          Usage: 0,
        },
      },
      want1: status.Success(),
    },
    {
      name: "Token not found",
      setup: func(mocked *tokenMocked, arg any) {
        a := arg.(*args)

        mocked.Token.EXPECT().
          Find(mock.Anything, a.verifyDTO.Token).
          Return(entity.Token{}, sql.ErrNoRows)
      },
      args: args{
        ctx: context.Background(),
        verifyDTO: &dto.TokenVerifyDTO{
          Token: util.RandomString(32),
          Usage: entity.TokenUsageResetPassword,
        },
      },
      want1: status.ErrBadRequest(sharedErr.ErrTokenNotFound),
    },
    {
      name: "Failed due to repository error",
      setup: func(mocked *tokenMocked, arg any) {
        a := arg.(*args)

        token := entity.Token{
          Token:     a.verifyDTO.Token,
          UserId:    types.MustCreateId(),
          Usage:     a.verifyDTO.Usage,
          ExpiredAt: time.Now().Add(time.Hour * 5),
        }

        mocked.Token.EXPECT().
          Find(mock.Anything, a.verifyDTO.Token).
          Return(token, dummyErr)
      },
      args: args{
        ctx: context.Background(),
        verifyDTO: &dto.TokenVerifyDTO{
          Token: util.RandomString(32),
          Usage: 0,
        },
      },
      want1: status.FromRepository(dummyErr, status.NullCode),
    },
    {
      name: "Failed due to different usage token",
      setup: func(mocked *tokenMocked, arg any) {
        a := arg.(*args)

        token := entity.Token{
          Token:     a.verifyDTO.Token,
          UserId:    types.MustCreateId(),
          Usage:     entity.TokenUsageVerification,
          ExpiredAt: time.Now().Add(time.Hour * 5),
        }

        mocked.Token.EXPECT().
          Find(mock.Anything, a.verifyDTO.Token).
          Return(token, nil)

      },
      args: args{
        ctx: context.Background(),
        verifyDTO: &dto.TokenVerifyDTO{
          Token: util.RandomString(32),
          Usage: entity.TokenUsageResetPassword,
        },
      },
      want1: status.ErrBadRequest(sharedErr.ErrTokenDifferentUsage),
    },
    {
      name: "Failed to delete token",
      setup: func(mocked *tokenMocked, arg any) {
        a := arg.(*args)

        token := entity.Token{
          Token:     a.verifyDTO.Token,
          UserId:    types.MustCreateId(),
          Usage:     entity.TokenUsageResetPassword,
          ExpiredAt: time.Now().Add(time.Hour * 5),
        }

        mocked.Token.EXPECT().
          Find(mock.Anything, a.verifyDTO.Token).
          Return(token, nil)

        mocked.Token.EXPECT().
          Delete(mock.Anything, a.verifyDTO.Token).
          Return(sql.ErrNoRows)

      },
      args: args{
        ctx: context.Background(),
        verifyDTO: &dto.TokenVerifyDTO{
          Token: util.RandomString(32),
          Usage: entity.TokenUsageResetPassword,
        },
      },
      want1: status.FromRepository(sql.ErrNoRows, status.NullCode),
    },
    {
      name: "Token expired",
      setup: func(mocked *tokenMocked, arg any) {
        a := arg.(*args)

        token := entity.Token{
          Token:     a.verifyDTO.Token,
          UserId:    types.MustCreateId(),
          Usage:     entity.TokenUsageResetPassword,
          ExpiredAt: time.Now().Add(time.Hour * -5),
        }

        mocked.Token.EXPECT().
          Find(mock.Anything, a.verifyDTO.Token).
          Return(token, nil)

        mocked.Token.EXPECT().
          Delete(mock.Anything, a.verifyDTO.Token).
          Return(nil)

      },
      args: args{
        ctx: context.Background(),
        verifyDTO: &dto.TokenVerifyDTO{
          Token: util.RandomString(32),
          Usage: entity.TokenUsageResetPassword,
        },
      },
      want1: status.ErrBadRequest(sharedErr.ErrTokenExpired),
    },
  }
  for _, tt := range tests {
    t1.Run(tt.name, func(t1 *testing.T) {
      mocked := newTokenMocked(t1)

      if tt.setup != nil {
        tt.setup(&mocked, &tt.args)
      }

      t := &tokenService{
        config:    mocked.Config,
        tokenRepo: mocked.Token,
        tracer:    mocked.Tracer,
      }

      got, got1 := t.Verify(tt.args.ctx, tt.args.verifyDTO)
      if !reflect.DeepEqual(got1, tt.want1) {
        t1.Errorf("Verify() got1 = %v, want %v", got1, tt.want1)
      }

      if got1.IsError() {
        return
      }

      require.True(t1, !got.Eq(types.NullId()))
    })
  }
}
