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
  "nexa/services/token/internal/domain/dto"
  "nexa/services/token/internal/domain/entity"
  repoMock "nexa/services/token/internal/domain/repository/mocks"
  errs "nexa/services/token/util/errors"
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
          Upsert(mock.Anything, mock.Anything).
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
          Upsert(mock.Anything, mock.Anything).
          Return(nil)
      },
      args: args{
        ctx: context.Background(),
        req: &dto.TokenCreateDTO{
          UserId: types.MustCreateId(),
          Usage:  entity.TokenUsageEmailVerification,
        },
      },
      want: dto.TokenResponseDTO{
        Usage:     entity.TokenUsageEmailVerification,
        ExpiredAt: time.Now().Add(duration),
      },
      want1: status.Created(),
    },
    {
      name: "Failed to request verification token",
      setup: func(mocked *tokenMocked, arg any) {
        //a := arg.(*args)
        mocked.Token.EXPECT().
          Upsert(mock.Anything, mock.Anything).
          Return(sql.ErrNoRows)
      },
      args: args{
        ctx: context.Background(),
        req: &dto.TokenCreateDTO{
          UserId: types.MustCreateId(),
          Usage:  entity.TokenUsageEmailVerification,
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
          Usage:     a.verifyDTO.ExpectedUsage,
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
          Token:         util.RandomString(32),
          ExpectedUsage: 0,
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
          Token:         util.RandomString(32),
          ExpectedUsage: entity.TokenUsageResetPassword,
        },
      },
      want1: status.ErrBadRequest(errs.ErrTokenNotFound),
    },
    {
      name: "Failed due to repository error",
      setup: func(mocked *tokenMocked, arg any) {
        a := arg.(*args)

        token := entity.Token{
          Token:     a.verifyDTO.Token,
          UserId:    types.MustCreateId(),
          Usage:     a.verifyDTO.ExpectedUsage,
          ExpiredAt: time.Now().Add(time.Hour * 5),
        }

        mocked.Token.EXPECT().
          Find(mock.Anything, a.verifyDTO.Token).
          Return(token, dummyErr)
      },
      args: args{
        ctx: context.Background(),
        verifyDTO: &dto.TokenVerifyDTO{
          Token:         util.RandomString(32),
          ExpectedUsage: 0,
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
          Usage:     entity.TokenUsageEmailVerification,
          ExpiredAt: time.Now().Add(time.Hour * 5),
        }

        mocked.Token.EXPECT().
          Delete(mock.Anything, a.verifyDTO.Token).
          Return(nil)

        mocked.Token.EXPECT().
          Find(mock.Anything, a.verifyDTO.Token).
          Return(token, nil)
      },
      args: args{
        ctx: context.Background(),
        verifyDTO: &dto.TokenVerifyDTO{
          Token:         util.RandomString(32),
          ExpectedUsage: entity.TokenUsageResetPassword,
        },
      },
      want1: status.ErrBadRequest(errs.ErrTokenDifferentUsage),
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
          Token:         util.RandomString(32),
          ExpectedUsage: entity.TokenUsageResetPassword,
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
          Token:         util.RandomString(32),
          ExpectedUsage: entity.TokenUsageResetPassword,
        },
      },
      want1: status.ErrBadRequest(errs.ErrTokenExpired),
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

func Test_tokenService_AuthVerify(t1 *testing.T) {
  type args struct {
    ctx       context.Context
    verifyDTO *dto.TokenAuthVerifyDTO
  }
  tests := []struct {
    name  string
    setup setupTokenTestFunc
    args  args
    want  status.Object
  }{
    {
      name: "Success auth verify",
      setup: func(mocked *tokenMocked, arg any) {
        a := arg.(*args)

        token := entity.Token{
          Token:     a.verifyDTO.Token,
          UserId:    a.verifyDTO.ExpectedUserId,
          Usage:     a.verifyDTO.ExpectedUsage,
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
        verifyDTO: &dto.TokenAuthVerifyDTO{
          TokenVerifyDTO: dto.TokenVerifyDTO{
            Token:         util.RandomString(0),
            ExpectedUsage: entity.TokenUsageLogin,
          },
          ExpectedUserId: types.MustCreateId(),
        },
      },
      want: status.Success(),
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
        verifyDTO: &dto.TokenAuthVerifyDTO{
          TokenVerifyDTO: dto.TokenVerifyDTO{
            Token:         util.RandomString(0),
            ExpectedUsage: entity.TokenUsageLogin,
          },
          ExpectedUserId: types.MustCreateId(),
        },
      },
      want: status.FromRepositoryOverride(sql.ErrNoRows, types.NewPair(status.BAD_REQUEST_ERROR, errs.ErrTokenNotFound)),
    },
    {
      name: "Token has different user",
      setup: func(mocked *tokenMocked, arg any) {
        a := arg.(*args)

        token := entity.Token{
          Token:     a.verifyDTO.Token,
          UserId:    types.MustCreateId(),
          Usage:     a.verifyDTO.ExpectedUsage,
          ExpiredAt: time.Now().Add(time.Hour * 5),
        }

        mocked.Token.EXPECT().
          Find(mock.Anything, a.verifyDTO.Token).
          Return(token, nil)
      },
      args: args{
        ctx: context.Background(),
        verifyDTO: &dto.TokenAuthVerifyDTO{
          TokenVerifyDTO: dto.TokenVerifyDTO{
            Token:         util.RandomString(0),
            ExpectedUsage: entity.TokenUsageLogin,
          },
          ExpectedUserId: types.MustCreateId(),
        },
      },
      want: status.ErrUnAuthorized(errs.ErrTokenDifferentUser),
    },
    {
      name: "Failed to delete token",
      setup: func(mocked *tokenMocked, arg any) {
        a := arg.(*args)

        token := entity.Token{
          Token:     a.verifyDTO.Token,
          UserId:    a.verifyDTO.ExpectedUserId,
          Usage:     a.verifyDTO.ExpectedUsage,
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
        verifyDTO: &dto.TokenAuthVerifyDTO{
          TokenVerifyDTO: dto.TokenVerifyDTO{
            Token:         util.RandomString(0),
            ExpectedUsage: entity.TokenUsageLogin,
          },
          ExpectedUserId: types.MustCreateId(),
        },
      },
      want: status.ErrInternal(sql.ErrNoRows),
    },
    {
      name: "Failed due to different usage token",
      setup: func(mocked *tokenMocked, arg any) {
        a := arg.(*args)

        token := entity.Token{
          Token:     a.verifyDTO.Token,
          UserId:    a.verifyDTO.ExpectedUserId,
          Usage:     entity.TokenUsageEmailVerification,
          ExpiredAt: time.Now().Add(time.Hour * 5),
        }

        mocked.Token.EXPECT().
          Delete(mock.Anything, a.verifyDTO.Token).
          Return(nil)

        mocked.Token.EXPECT().
          Find(mock.Anything, a.verifyDTO.Token).
          Return(token, nil)
      },
      args: args{
        ctx: context.Background(),
        verifyDTO: &dto.TokenAuthVerifyDTO{
          TokenVerifyDTO: dto.TokenVerifyDTO{
            Token:         util.RandomString(32),
            ExpectedUsage: entity.TokenUsageResetPassword,
          },
          ExpectedUserId: types.MustCreateId(),
        },
      },
      want: status.ErrBadRequest(errs.ErrTokenDifferentUsage),
    },
    {
      name: "Token expired",
      setup: func(mocked *tokenMocked, arg any) {
        a := arg.(*args)

        token := entity.Token{
          Token:     a.verifyDTO.Token,
          UserId:    a.verifyDTO.ExpectedUserId,
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
        verifyDTO: &dto.TokenAuthVerifyDTO{
          TokenVerifyDTO: dto.TokenVerifyDTO{
            Token:         util.RandomString(32),
            ExpectedUsage: entity.TokenUsageResetPassword,
          },
          ExpectedUserId: types.MustCreateId(),
        },
      },
      want: status.ErrBadRequest(errs.ErrTokenExpired),
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

      if got := t.AuthVerify(tt.args.ctx, tt.args.verifyDTO); !reflect.DeepEqual(got, tt.want) {
        t1.Errorf("AuthVerify() = %v, want %v", got, tt.want)
      }
    })
  }
}
