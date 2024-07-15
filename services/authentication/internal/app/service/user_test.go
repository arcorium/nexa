package service

import (
  "context"
  "database/sql"
  sharedConst "github.com/arcorium/nexa/shared/constant"
  sharedDto "github.com/arcorium/nexa/shared/dto"
  sharedErr "github.com/arcorium/nexa/shared/errors"
  sharedJwt "github.com/arcorium/nexa/shared/jwt"
  "github.com/arcorium/nexa/shared/status"
  "github.com/arcorium/nexa/shared/types"
  sharedUow "github.com/arcorium/nexa/shared/uow"
  uowMock "github.com/arcorium/nexa/shared/uow/mocks"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  "github.com/arcorium/nexa/shared/util/repo"
  "github.com/brianvoe/gofakeit/v7"
  "github.com/stretchr/testify/mock"
  "github.com/stretchr/testify/require"
  "go.opentelemetry.io/otel/trace"
  "go.opentelemetry.io/otel/trace/noop"
  "nexa/services/authentication/constant"
  "nexa/services/authentication/internal/app/uow"
  "nexa/services/authentication/internal/domain/dto"
  "nexa/services/authentication/internal/domain/entity"
  extMock "nexa/services/authentication/internal/domain/external/mocks"
  "nexa/services/authentication/internal/domain/mapper"
  repoMock "nexa/services/authentication/internal/domain/repository/mocks"
  "reflect"
  "testing"
  "time"
)

func newUserMocked(t *testing.T) userMocked {
  // Tracer
  provider := noop.NewTracerProvider()
  return userMocked{
    UOW:         uowMock.NewUnitOfWorkMock[uow.UserStorage](t),
    Cred:        repoMock.NewCredentialMock(t),
    User:        repoMock.NewUserMock(t),
    Profile:     repoMock.NewProfileMock(t),
    RoleClient:  extMock.NewRoleClientMock(t),
    MailClient:  extMock.NewMailerClientMock(t),
    TokenClient: extMock.NewTokenClientMock(t),
    Tracer:      provider.Tracer("MOCK"),
  }
}

type userMocked struct {
  UOW         *uowMock.UnitOfWorkMock[uow.UserStorage]
  Cred        *repoMock.CredentialMock
  User        *repoMock.UserMock
  Profile     *repoMock.ProfileMock
  RoleClient  *extMock.RoleClientMock
  MailClient  *extMock.MailerClientMock
  TokenClient *extMock.TokenClientMock
  Tracer      trace.Tracer
}

func (m *userMocked) defaultUOWMock() {
  m.UOW.EXPECT().
    Repositories().
    Return(uow.NewStorage(m.Profile, m.User))
}

func (m *userMocked) txProxy() {
  m.UOW.On("DoTx", mock.Anything, mock.Anything).
      Return(func(ctx context.Context, f sharedUow.UOWBlock[uow.UserStorage]) error {
        return f(ctx, uow.NewStorage(m.Profile, m.User))
      })
}

type setupUserTestFunc func(mocked *userMocked, arg any, want any)

func Test_userService_BannedUser(t *testing.T) {
  t.Parallel()

  type args struct {
    ctx       context.Context
    bannedDto *dto.UserBannedDTO
  }
  tests := []struct {
    name  string
    setup setupUserTestFunc
    args  args
    want  status.Object
  }{
    {
      name: "Success banned user",
      setup: func(mocked *userMocked, arg any, want any) {
        //a := arg.(*args)
        mocked.defaultUOWMock()

        //entity := a.bannedDto.ToDomain()
        mocked.User.
          EXPECT().
          Patch(mock.Anything, mock.Anything).
          Return(nil)
      },
      args: args{
        ctx: context.Background(),
        bannedDto: &dto.UserBannedDTO{
          Id:       types.MustCreateId(),
          Duration: time.Hour * time.Duration(gofakeit.Hour()),
        },
      },
      want: status.Success(),
    },
    {
      name: "Failed banned user",
      setup: func(mocked *userMocked, arg any, want any) {
        //a := arg.(*args)

        mocked.defaultUOWMock()

        //entity := a.bannedDto.ToDomain()
        mocked.User.
          EXPECT().
          Patch(mock.Anything, mock.Anything).
          Return(sql.ErrNoRows)
      },
      args: args{
        ctx: context.Background(),
        bannedDto: &dto.UserBannedDTO{
          Id:       types.MustCreateId(),
          Duration: time.Hour * time.Duration(gofakeit.Hour()),
        },
      },
      want: status.FromRepository(sql.ErrNoRows, status.NullCode),
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newUserMocked(t)

      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, tt.want)
      }

      u := userService{
        unit:     mocked.UOW,
        credRepo: mocked.Cred,
        tracer:   mocked.Tracer,
        config: UserConfig{
          RoleClient:  mocked.RoleClient,
          MailClient:  mocked.MailClient,
          TokenClient: mocked.TokenClient,
        },
      }

      if got := u.BannedUser(tt.args.ctx, tt.args.bannedDto); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("BannedUser() = %v, want %v", got, tt.want)
      }
    })
  }
}

func Test_userService_Create(t *testing.T) {
  t.Parallel()

  type args struct {
    ctx   context.Context
    input *dto.UserCreateDTO
  }
  tests := []struct {
    name         string
    setup        setupUserTestFunc
    args         args
    expectNullId bool
    want         status.Object
  }{
    {
      name: "Success create user",
      setup: func(mocked *userMocked, arg any, want any) {
        mocked.UOW.
          EXPECT().
          DoTx(mock.Anything, mock.Anything).
          Return(nil)

        //mocked.User.
        //  EXPECT().
        //  Create(mock.Anything, mock.Anything).
        //  Return(nil)
        //
        //mocked.Profile.
        //  EXPECT().
        //  Create(mock.Anything, mock.Anything).
        //  Return(nil)
      },
      args: args{
        ctx: context.Background(),
        input: &dto.UserCreateDTO{
          Username:  gofakeit.Username(),
          Email:     types.Email(gofakeit.Email()),
          Password:  types.Password(gofakeit.AppName()),
          FirstName: gofakeit.FirstName(),
          LastName:  types.SomeNullable(gofakeit.LastName()),
          Bio:       types.NullableString{},
        },
      },
      expectNullId: false,
      want:         status.Created(),
    },
    {
      name: "Failed create user on transaction",
      setup: func(mocked *userMocked, arg any, want any) {
        var err error = sql.ErrNoRows

        mocked.UOW.
          EXPECT().
          DoTx(mock.Anything, mock.Anything).
          Return(err)

        //mocked.User.
        //  EXPECT().
        //  Create(mock.Anything, mock.Anything).
        //  Return(err)
      },
      args: args{
        ctx: context.Background(),
        input: &dto.UserCreateDTO{
          Username:  gofakeit.Username(),
          Email:     types.Email(gofakeit.Email()),
          Password:  types.Password(gofakeit.AppName()),
          FirstName: gofakeit.FirstName(),
          LastName:  types.SomeNullable(gofakeit.LastName()),
        },
      },
      expectNullId: true,
      want:         status.FromRepository(sql.ErrNoRows, status.NullCode),
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newUserMocked(t)

      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, tt.want)
      }

      u := userService{
        unit:     mocked.UOW,
        credRepo: mocked.Cred,
        tracer:   mocked.Tracer,
        config: UserConfig{
          RoleClient:  mocked.RoleClient,
          MailClient:  mocked.MailClient,
          TokenClient: mocked.TokenClient,
        },
      }

      got, got1 := u.Create(tt.args.ctx, tt.args.input)
      if (got == types.NullId()) != tt.expectNullId {
        t.Errorf("Create() got = %v, want %v", got, tt.want)
      }
      if !reflect.DeepEqual(got1, tt.want) {
        t.Errorf("Create() got1 = %v, want %v", got1, tt.want)
      }
    })
  }
}

func Test_userService_DeleteById(t *testing.T) {
  t.Parallel()

  self := generateUserClaims(generateRole(constant.AUTHN_DELETE_USER))

  type args struct {
    ctx    context.Context
    userId types.Id
  }
  tests := []struct {
    name  string
    setup setupUserTestFunc
    args  args
    want  status.Object
  }{
    {
      name: "Success self delete user",
      setup: func(mocked *userMocked, arg any, want any) {
        a := arg.(*args)

        mocked.txProxy()

        mocked.User.EXPECT().
          Delete(mock.Anything, a.userId).
          Return(nil)

        mocked.Cred.EXPECT().
          DeleteByUserId(mock.Anything, mock.Anything).
          Return(nil)

        mocked.RoleClient.EXPECT().
          RemoveUserRoles(mock.Anything, a.userId).
          Return(nil)
      },
      args: args{
        userId: types.Must(types.IdFromString(self.UserId)),
        ctx: context.WithValue(context.Background(), sharedConst.USER_CLAIMS_CONTEXT_KEY,
          self),
      },
      want: status.Deleted(),
    },
    {
      name: "Failed to delete other user due to permissions",
      setup: func(mocked *userMocked, arg any, want any) {
        a := arg.(*args)

        mocked.txProxy()

        mocked.User.EXPECT().
          Delete(mock.Anything, a.userId).
          Return(sql.ErrNoRows)

      },
      args: args{
        userId: types.MustCreateId(),
        ctx: context.WithValue(context.Background(), sharedConst.USER_CLAIMS_CONTEXT_KEY,
          generateUserClaims(generateRole(constant.AUTHN_DELETE_USER, constant.AUTHN_DELETE_USER_ARB))),
      },
      want: status.FromRepository(sql.ErrNoRows, status.NullCode),
    },
    {
      name: "Failed to delete user credentials",
      setup: func(mocked *userMocked, arg any, want any) {
        a := arg.(*args)

        mocked.txProxy()

        mocked.User.EXPECT().
          Delete(mock.Anything, a.userId).
          Return(nil)

        mocked.Cred.EXPECT().
          DeleteByUserId(mock.Anything, mock.Anything).
          Return(dummyErr)
      },
      args: args{
        userId: types.Must(types.IdFromString(self.UserId)),
        ctx: context.WithValue(context.Background(), sharedConst.USER_CLAIMS_CONTEXT_KEY,
          self),
      },
      want: status.FromRepository(dummyErr, status.NullCode),
    },
    {
      name: "User has no credentials",
      setup: func(mocked *userMocked, arg any, want any) {
        a := arg.(*args)

        mocked.txProxy()

        mocked.User.EXPECT().
          Delete(mock.Anything, a.userId).
          Return(nil)

        mocked.Cred.EXPECT().
          DeleteByUserId(mock.Anything, mock.Anything).
          Return(sql.ErrNoRows)

        mocked.RoleClient.EXPECT().
          RemoveUserRoles(mock.Anything, a.userId).
          Return(nil)
      },
      args: args{
        userId: types.Must(types.IdFromString(self.UserId)),
        ctx: context.WithValue(context.Background(), sharedConst.USER_CLAIMS_CONTEXT_KEY,
          self),
      },
      want: status.Deleted(),
    },
    {
      name: "Failed to remove user roles",
      setup: func(mocked *userMocked, arg any, want any) {
        a := arg.(*args)

        mocked.txProxy()

        mocked.User.EXPECT().
          Delete(mock.Anything, a.userId).
          Return(nil)

        mocked.Cred.EXPECT().
          DeleteByUserId(mock.Anything, mock.Anything).
          Return(sql.ErrNoRows)

        mocked.RoleClient.EXPECT().
          RemoveUserRoles(mock.Anything, a.userId).
          Return(dummyErr)
      },
      args: args{
        userId: types.Must(types.IdFromString(self.UserId)),
        ctx: context.WithValue(context.Background(), sharedConst.USER_CLAIMS_CONTEXT_KEY,
          self),
      },
      want: status.ErrExternal(dummyErr),
    },
    {
      name: "Success delete other user",
      setup: func(mocked *userMocked, arg any, want any) {
        a := arg.(*args)

        mocked.txProxy()

        mocked.User.EXPECT().
          Delete(mock.Anything, a.userId).
          Return(nil)

        mocked.Cred.EXPECT().
          DeleteByUserId(mock.Anything, mock.Anything).
          Return(nil)

        mocked.RoleClient.EXPECT().
          RemoveUserRoles(mock.Anything, a.userId).
          Return(nil)
      },
      args: args{
        userId: types.MustCreateId(),
        ctx: context.WithValue(context.Background(), sharedConst.USER_CLAIMS_CONTEXT_KEY,
          generateUserClaims(generateRole(constant.AUTHN_DELETE_USER_ARB, constant.AUTHN_DELETE_USER))),
      },
      want: status.Deleted(),
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newUserMocked(t)

      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, tt.want)
      }

      u := userService{
        unit:     mocked.UOW,
        credRepo: mocked.Cred,
        tracer:   mocked.Tracer,
        config: UserConfig{
          RoleClient:  mocked.RoleClient,
          MailClient:  mocked.MailClient,
          TokenClient: mocked.TokenClient,
        },
      }

      if got := u.DeleteById(tt.args.ctx, tt.args.userId); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("DeleteById() = %v, want %v", got, tt.want)
      }
    })
  }
}

func Test_userService_EmailVerificationRequest(t *testing.T) {
  t.Parallel()

  type args struct {
    ctx context.Context
  }
  tests := []struct {
    name  string
    setup setupUserTestFunc
    args  args
    want1 status.Object
  }{
    {
      name: "Success request an email verification",
      setup: func(mocked *userMocked, arg any, want any) {
        a := arg.(*args)
        claims, err := sharedJwt.GetUserClaimsFromCtx(a.ctx)
        require.NoError(t, err)
        userId := types.Must(types.IdFromString(claims.UserId))

        mocked.TokenClient.EXPECT().
            Generate(mock.Anything, &dto.TokenGenerationDTO{
              UserId: userId,
              Usage:  dto.TokenUsageEmailVerification,
            }).
            Return(dto.TokenResponseDTO{
              Token:     sharedUtil.RandomString(32),
              Usage:     dto.TokenUsageEmailVerification,
              ExpiredAt: time.Now().Add(time.Hour * 24),
            }, nil)

        mocked.UOW.EXPECT().
          Repositories().
          Return(uow.NewStorage(mocked.Profile, mocked.User))

        mocked.User.EXPECT().
          FindByIds(mock.Anything, userId).
          Return([]entity.User{{}}, nil)

        mocked.MailClient.EXPECT().
          SendEmailVerification(mock.Anything, mock.Anything).
          Return(nil)
      },
      args: args{
        ctx: generateClaimsCtx(),
      },

      want1: status.Success(),
    },
    {
      name: "User not found",
      setup: func(mocked *userMocked, arg any, want any) {
        a := arg.(*args)
        claims, err := sharedJwt.GetUserClaimsFromCtx(a.ctx)
        require.NoError(t, err)
        userId := types.Must(types.IdFromString(claims.UserId))

        mocked.UOW.EXPECT().
          Repositories().
          Return(uow.NewStorage(mocked.Profile, mocked.User))

        mocked.User.EXPECT().
          FindByIds(mock.Anything, userId).
          Return(nil, sql.ErrNoRows)

      },
      args: args{
        ctx: generateClaimsCtx(),
      },

      want1: status.FromRepository(sql.ErrNoRows, status.NullCode),
    },
    {
      name: "Failed generate token",
      setup: func(mocked *userMocked, arg any, want any) {
        a := arg.(*args)
        mocked.defaultUOWMock()

        claims, err := sharedJwt.GetUserClaimsFromCtx(a.ctx)
        require.NoError(t, err)
        userId := types.Must(types.IdFromString(claims.UserId))

        mocked.TokenClient.EXPECT().
            Generate(mock.Anything, &dto.TokenGenerationDTO{
              UserId: userId,
              Usage:  dto.TokenUsageEmailVerification,
            }).
          Return(dto.TokenResponseDTO{}, dummyErr)

        mocked.User.EXPECT().
          FindByIds(mock.Anything, userId).
          Return([]entity.User{{}}, nil)
      },
      args: args{
        ctx: generateClaimsCtx(),
      },

      want1: status.ErrExternal(dummyErr),
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newUserMocked(t)

      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, nil)
      }

      u := userService{
        unit:     mocked.UOW,
        credRepo: mocked.Cred,
        tracer:   mocked.Tracer,
        config: UserConfig{
          RoleClient:  mocked.RoleClient,
          MailClient:  mocked.MailClient,
          TokenClient: mocked.TokenClient,
        },
      }

      got, got1 := u.EmailVerificationRequest(tt.args.ctx)

      if !got1.IsError() {
        require.Greater(t, len(got.Token), 0)
        require.True(t, got.Usage == dto.TokenUsageEmailVerification)
        require.True(t, got.ExpiredAt.After(time.Now()))
      }

      if !reflect.DeepEqual(got1, tt.want1) {
        t.Errorf("EmailVerificationRequest() got1 = %v, want %v", got1, tt.want1)
      }
    })
  }
}

func Test_userService_FindAll(t *testing.T) {
  t.Parallel()

  type args struct {
    ctx      context.Context
    pagedDto sharedDto.PagedElementDTO
  }
  tests := []struct {
    name  string
    setup setupUserTestFunc
    args  args
    want  sharedDto.PagedElementResult[dto.UserResponseDTO]
    want1 status.Object
  }{
    {
      name: "Success get users",
      setup: func(mocked *userMocked, arg any, want any) {
        a := arg.(*args)
        w := want.(*sharedDto.PagedElementResult[dto.UserResponseDTO])
        mocked.defaultUOWMock()

        mocked.User.EXPECT().
          Get(mock.Anything, a.pagedDto.ToQueryParam()).
          Return(repo.PaginatedResult[entity.User]{nil, w.TotalElements, w.Element}, nil)

      },
      args: args{
        ctx: context.Background(),
        pagedDto: sharedDto.PagedElementDTO{
          Element: 2,
          Page:    2,
        },
      },
      want: sharedDto.PagedElementResult[dto.UserResponseDTO]{
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
      setup: func(mocked *userMocked, arg any, want any) {
        a := arg.(*args)
        //w := want.(*sharedDto.PagedElementResult[dto.UserResponseDTO])
        mocked.defaultUOWMock()

        mocked.User.EXPECT().
          Get(mock.Anything, a.pagedDto.ToQueryParam()).
          Return(repo.PaginatedResult[entity.User]{nil, 0, 0}, sql.ErrNoRows)

      },
      args: args{
        ctx: context.Background(),
        pagedDto: sharedDto.PagedElementDTO{
          Element: 2,
          Page:    2,
        },
      },
      want: sharedDto.PagedElementResult[dto.UserResponseDTO]{
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
      mocked := newUserMocked(t)

      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, &tt.want)
      }

      u := userService{
        unit:     mocked.UOW,
        credRepo: mocked.Cred,
        tracer:   mocked.Tracer,
        config: UserConfig{
          RoleClient:  mocked.RoleClient,
          MailClient:  mocked.MailClient,
          TokenClient: mocked.TokenClient,
        },
      }

      got, got1 := u.GetAll(tt.args.ctx, tt.args.pagedDto)
      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("GetAll() got = %v, want %v", got, tt.want)
      }
      if !reflect.DeepEqual(got1, tt.want1) {
        t.Errorf("GetAll() got1 = %v, want %v", got1, tt.want1)
      }
    })
  }
}

func Test_userService_FindByIds(t *testing.T) {
  t.Parallel()

  type args struct {
    ctx context.Context
    ids []types.Id
  }
  tests := []struct {
    name  string
    setup setupUserTestFunc
    args  args
    want  []dto.UserResponseDTO
    want1 status.Object
  }{
    {
      name: "Success find by ids",
      setup: func(mocked *userMocked, arg any, want any) {
        a := arg.(*args)
        //w := want.(*dto.UserResponseDTO)
        mocked.defaultUOWMock()

        userIds := sharedUtil.CastSlice(a.ids, sharedUtil.ToAny[types.Id])

        mocked.User.
          EXPECT().
          FindByIds(mock.Anything, userIds...).
            Return(sharedUtil.GenerateMultiple(2, func() entity.User {
              return entity.User{
                Id:       dummyId,
                Username: "Username",
                Email:    types.Email("email@email.com"),
              }
            }), nil)
      },
      args: args{
        ctx: context.Background(),
        ids: []types.Id{types.MustCreateId(), types.MustCreateId()},
      },
      want: sharedUtil.GenerateMultiple(2, func() dto.UserResponseDTO {
        return mapper.ToUserResponse(&entity.User{
          Id:       dummyId,
          Username: "Username",
          Email:    types.Email("email@email.com"),
        })
      }),
      want1: status.Success(),
    },
    {
      name: "Failed find by ids",
      setup: func(mocked *userMocked, arg any, want any) {
        a := arg.(*args)
        //w := want.(*dto.UserResponseDTO)
        mocked.defaultUOWMock()

        userIds := sharedUtil.CastSlice(a.ids, sharedUtil.ToAny[types.Id])

        mocked.User.
          EXPECT().
          FindByIds(mock.Anything, userIds...).
          Return(nil, sql.ErrNoRows)
      },
      args: args{
        ctx: context.Background(),
        ids: []types.Id{types.MustCreateId(), types.MustCreateId()},
      },
      want:  nil,
      want1: status.FromRepository(sql.ErrNoRows, status.NullCode),
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newUserMocked(t)

      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, &tt.want)
      }

      u := userService{
        unit:     mocked.UOW,
        credRepo: mocked.Cred,
        tracer:   mocked.Tracer,
        config: UserConfig{
          RoleClient:  mocked.RoleClient,
          MailClient:  mocked.MailClient,
          TokenClient: mocked.TokenClient,
        },
      }

      got, got1 := u.FindByIds(tt.args.ctx, tt.args.ids...)
      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("FindByIds() got = %v, want %v", got, tt.want)
      }
      if !reflect.DeepEqual(got1, tt.want1) {
        t.Errorf("FindByIds() got1 = %v, want %v", got1, tt.want1)
      }
    })
  }
}

func Test_userService_ForgotPassword(t *testing.T) {
  t.Parallel()

  type args struct {
    ctx   context.Context
    email types.Email
  }
  tests := []struct {
    name  string
    setup setupUserTestFunc
    args  args
    want  dto.TokenResponseDTO
    want1 status.Object
  }{
    {
      name: "Success get token for reset password",
      setup: func(mocked *userMocked, arg any, want any) {
        a := arg.(*args)
        w := want.(*dto.TokenResponseDTO)
        mocked.defaultUOWMock()

        email := sharedUtil.ToAny(a.email)

        user := entity.User{
          Id:       types.MustCreateId(),
          Username: gofakeit.Username(),
          Email:    types.Email(gofakeit.Email()),
        }

        mocked.User.EXPECT().
          FindByEmails(mock.Anything, email).
          Return([]entity.User{user}, nil)

        mocked.TokenClient.EXPECT().
            Generate(mock.Anything, &dto.TokenGenerationDTO{
              UserId: user.Id,
              Usage:  dto.TokenUsageResetPassword,
            }).Return(*w, nil)

        mocked.MailClient.EXPECT().
            SendForgotPassword(mock.Anything, &dto.SendForgotPasswordDTO{
              Recipient: user.Email,
              Token:     w.Token,
            }).Return(nil)

      },
      args: args{
        ctx:   context.Background(),
        email: types.Email(gofakeit.Email()),
      },
      want: dto.TokenResponseDTO{
        Token:     sharedUtil.RandomString(32),
        Usage:     dto.TokenUsageResetPassword,
        ExpiredAt: time.Now().Add(time.Hour * 24),
      },
      want1: status.Success(),
    },
    {
      name: "failed to get token for reset password",
      setup: func(mocked *userMocked, arg any, want any) {
        a := arg.(*args)
        //w := want.(*dto.TokenResponseDTO)
        mocked.defaultUOWMock()

        email := sharedUtil.ToAny(a.email)

        user := entity.User{
          Id:       types.MustCreateId(),
          Username: gofakeit.Username(),
          Email:    types.Email(gofakeit.Email()),
        }

        mocked.User.EXPECT().
          FindByEmails(mock.Anything, email).
          Return([]entity.User{user}, nil)

        mocked.TokenClient.EXPECT().
            Generate(mock.Anything, &dto.TokenGenerationDTO{
              UserId: user.Id,
              Usage:  dto.TokenUsageResetPassword,
            }).Return(dto.TokenResponseDTO{}, dummyErr)

      },
      args: args{
        ctx:   context.Background(),
        email: types.Email(gofakeit.Email()),
      },
      want:  dto.TokenResponseDTO{},
      want1: status.ErrExternal(dummyErr),
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newUserMocked(t)

      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, &tt.want)
      }

      u := userService{
        unit:     mocked.UOW,
        credRepo: mocked.Cred,
        tracer:   mocked.Tracer,
        config: UserConfig{
          RoleClient:  mocked.RoleClient,
          MailClient:  mocked.MailClient,
          TokenClient: mocked.TokenClient,
        },
      }

      got, got1 := u.ForgotPassword(tt.args.ctx, tt.args.email)

      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("ForgotPassword() got = %v, want %v", got, tt.want)
      }
      if !reflect.DeepEqual(got1, tt.want1) {
        t.Errorf("ForgotPassword() got1 = %v, want %v", got1, tt.want1)
      }
    })
  }
}

func Test_userService_ResetPassword(t *testing.T) {
  t.Parallel()

  type args struct {
    ctx   context.Context
    input *dto.ResetUserPasswordDTO
  }
  tests := []struct {
    name  string
    setup setupUserTestFunc
    args  args
    want  status.Object
  }{
    {
      name: "Success reset password by token and force logout",
      setup: func(mocked *userMocked, arg any, want any) {
        a := arg.(*args)
        mocked.defaultUOWMock()

        mocked.User.EXPECT().
          Patch(mock.Anything, mock.Anything).Return(nil)

        mocked.Cred.EXPECT().
          DeleteByUserId(mock.Anything, a.input.UserId).
          Return(nil)
      },
      args: args{
        ctx: context.Background(),
        input: &dto.ResetUserPasswordDTO{
          UserId:      types.MustCreateId(),
          LogoutAll:   true,
          NewPassword: types.Password(gofakeit.Username()),
        },
      },
      want: status.Updated(),
    },
    {
      name: "Failed to update user on repository",
      setup: func(mocked *userMocked, arg any, want any) {
        mocked.defaultUOWMock()

        mocked.User.EXPECT().
          Patch(mock.Anything, mock.Anything).
          Return(sql.ErrNoRows)
      },
      args: args{
        ctx: generateClaimsCtx(constant.AUTHN_UPDATE_USER_ARB),
        input: &dto.ResetUserPasswordDTO{
          UserId:      types.MustCreateId(),
          LogoutAll:   true,
          NewPassword: types.Password(gofakeit.Username()),
        },
      },
      want: status.FromRepository(sql.ErrNoRows, status.NullCode),
    },
    {
      name: "Failed to force logout",
      setup: func(mocked *userMocked, arg any, want any) {
        a := arg.(*args)
        mocked.defaultUOWMock()

        mocked.User.EXPECT().
          Patch(mock.Anything, mock.Anything).Return(nil)

        mocked.Cred.EXPECT().
          DeleteByUserId(mock.Anything, a.input.UserId).
          Return(dummyErr)
      },
      args: args{
        ctx: context.Background(),
        input: &dto.ResetUserPasswordDTO{
          UserId:      types.MustCreateId(),
          LogoutAll:   true,
          NewPassword: types.Password(gofakeit.Username()),
        },
      },
      want: status.Updated(),
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newUserMocked(t)

      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, &tt.want)
      }

      u := userService{
        unit:     mocked.UOW,
        credRepo: mocked.Cred,
        tracer:   mocked.Tracer,
        config: UserConfig{
          RoleClient:  mocked.RoleClient,
          MailClient:  mocked.MailClient,
          TokenClient: mocked.TokenClient,
        },
      }

      if got := u.ResetPassword(tt.args.ctx, tt.args.input); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("ResetPassword() = %v, want %v", got, tt.want)
      }
    })
  }
}

func Test_userService_ResetPasswordWithToken(t *testing.T) {
  t.Parallel()

  type args struct {
    ctx   context.Context
    input *dto.ResetPasswordWithTokenDTO
  }
  tests := []struct {
    name  string
    setup setupUserTestFunc
    args  args
    want  status.Object
  }{
    {
      name: "Success reset password by token and force logout",
      setup: func(mocked *userMocked, arg any, want any) {
        a := arg.(*args)
        mocked.defaultUOWMock()

        mocked.TokenClient.EXPECT().
            Verify(mock.Anything, &dto.TokenVerificationDTO{
              Token:   a.input.Token,
              Purpose: dto.TokenUsageResetPassword,
            }).Return(dummyId, nil)

        mocked.User.EXPECT().
          Patch(mock.Anything, mock.Anything).Return(nil)

        mocked.Cred.EXPECT().
          DeleteByUserId(mock.Anything, dummyId).
          Return(nil)
      },
      args: args{
        ctx: context.Background(),
        input: &dto.ResetPasswordWithTokenDTO{
          Token:       sharedUtil.RandomString(32),
          LogoutAll:   true,
          NewPassword: types.Password(gofakeit.Username()),
        },
      },
      want: status.Updated(),
    },
    {
      name: "Failed reset password with token due to token verification",
      setup: func(mocked *userMocked, arg any, want any) {
        a := arg.(*args)
        mocked.TokenClient.EXPECT().
            Verify(mock.Anything, &dto.TokenVerificationDTO{
              Token:   a.input.Token,
              Purpose: dto.TokenUsageResetPassword,
            }).Return(types.NullId(), dummyErr)

      },
      args: args{
        ctx: generateClaimsCtx(constant.AUTHN_UPDATE_USER_ARB),
        input: &dto.ResetPasswordWithTokenDTO{
          Token:       sharedUtil.RandomString(32),
          LogoutAll:   true,
          NewPassword: types.Password(gofakeit.Username()),
        },
      },
      want: status.ErrBadRequest(dummyErr),
    },
    {
      name: "Failed to force logout",
      setup: func(mocked *userMocked, arg any, want any) {
        a := arg.(*args)
        mocked.defaultUOWMock()

        mocked.TokenClient.EXPECT().
            Verify(mock.Anything, &dto.TokenVerificationDTO{
              Token:   a.input.Token,
              Purpose: dto.TokenUsageResetPassword,
            }).Return(dummyId, nil)

        mocked.User.EXPECT().
          Patch(mock.Anything, mock.Anything).Return(nil)

        mocked.Cred.EXPECT().
          DeleteByUserId(mock.Anything, dummyId).
          Return(dummyErr)
      },
      args: args{
        ctx: generateClaimsCtx(constant.AUTHN_UPDATE_USER_ARB),
        input: &dto.ResetPasswordWithTokenDTO{
          Token:       sharedUtil.RandomString(32),
          LogoutAll:   true,
          NewPassword: types.Password(gofakeit.Username()),
        },
      },
      want: status.Updated(),
    },
    {
      name: "Failed to update user on repository",
      setup: func(mocked *userMocked, arg any, want any) {
        a := arg.(*args)
        mocked.defaultUOWMock()

        mocked.TokenClient.EXPECT().
            Verify(mock.Anything, &dto.TokenVerificationDTO{
              Token:   a.input.Token,
              Purpose: dto.TokenUsageResetPassword,
            }).Return(dummyId, nil)

        mocked.User.EXPECT().
          Patch(mock.Anything, mock.Anything).
          Return(sql.ErrNoRows)
      },
      args: args{
        ctx: generateClaimsCtx(constant.AUTHN_UPDATE_USER_ARB),
        input: &dto.ResetPasswordWithTokenDTO{
          Token:       sharedUtil.RandomString(32),
          LogoutAll:   false,
          NewPassword: types.Password(gofakeit.Username()),
        },
      },
      want: status.FromRepository(sql.ErrNoRows, status.NullCode),
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newUserMocked(t)

      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, &tt.want)
      }

      u := userService{
        unit:     mocked.UOW,
        credRepo: mocked.Cred,
        tracer:   mocked.Tracer,
        config: UserConfig{
          RoleClient:  mocked.RoleClient,
          MailClient:  mocked.MailClient,
          TokenClient: mocked.TokenClient,
        },
      }

      if got := u.ResetPasswordWithToken(tt.args.ctx, tt.args.input); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("ResetPassword() = %v, want %v", got, tt.want)
      }
    })
  }
}

func Test_userService_Update(t *testing.T) {
  t.Parallel()

  type args struct {
    ctx   context.Context
    input *dto.UserUpdateDTO
  }
  tests := []struct {
    name  string
    setup setupUserTestFunc
    args  args
    want  status.Object
  }{
    {
      name: "Success update self user",
      setup: func(mocked *userMocked, arg any, want any) {
        a := arg.(*args)

        mocked.defaultUOWMock()

        claims, err := sharedJwt.GetUserClaimsFromCtx(a.ctx)
        require.NoError(t, err)

        a.input.Id = types.Must(types.IdFromString(claims.UserId))

        ent := a.input.ToDomain()
        mocked.User.EXPECT().
          Patch(mock.Anything, &ent).
          Return(nil)

      },
      args: args{
        ctx: generateClaimsCtx(constant.AUTHN_UPDATE_USER),
        input: &dto.UserUpdateDTO{
          Id:       types.MustCreateId(),
          Username: types.SomeNullable(gofakeit.Username()),
          Email:    types.SomeNullable(types.Email(gofakeit.Email())),
        },
      },
      want: status.Updated(),
    },
    {
      name: "Success update other user",
      setup: func(mocked *userMocked, arg any, want any) {
        a := arg.(*args)

        mocked.defaultUOWMock()

        ent := a.input.ToDomain()
        mocked.User.EXPECT().
          Patch(mock.Anything, &ent).
          Return(nil)
      },
      args: args{
        ctx: generateClaimsCtx(constant.AUTHN_UPDATE_USER_ARB),
        input: &dto.UserUpdateDTO{
          Id:       types.MustCreateId(),
          Username: types.SomeNullable(gofakeit.Username()),
          Email:    types.SomeNullable(types.Email(gofakeit.Email())),
        },
      },
      want: status.Updated(),
    },
    {
      name: "Failed to update other user due to permission",
      setup: func(mocked *userMocked, arg any, want any) {

      },
      args: args{
        ctx: generateClaimsCtx(constant.AUTHN_UPDATE_USER),
        input: &dto.UserUpdateDTO{
          Id:       types.MustCreateId(),
          Username: types.SomeNullable(gofakeit.Username()),
          Email:    types.SomeNullable(types.Email(gofakeit.Email())),
        },
      },
      want: status.ErrUnAuthorized(sharedErr.ErrUnauthorizedPermission),
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newUserMocked(t)

      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, &tt.want)
      }

      u := userService{
        unit:     mocked.UOW,
        credRepo: mocked.Cred,
        tracer:   mocked.Tracer,
        config: UserConfig{
          RoleClient:  mocked.RoleClient,
          MailClient:  mocked.MailClient,
          TokenClient: mocked.TokenClient,
        },
      }

      if got := u.Update(tt.args.ctx, tt.args.input); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("Update() = %v, want %v", got, tt.want)
      }
    })
  }
}

func Test_userService_UpdatePassword(t *testing.T) {
  t.Parallel()

  type args struct {
    ctx   context.Context
    input *dto.UserUpdatePasswordDTO
  }
  tests := []struct {
    name  string
    setup setupUserTestFunc
    args  args
    want  status.Object
  }{
    {
      name: "Success update self password",
      setup: func(mocked *userMocked, arg any, want any) {
        a := arg.(*args)

        mocked.defaultUOWMock()

        claims, err := sharedJwt.GetUserClaimsFromCtx(a.ctx)
        require.NoError(t, err)

        a.input.Id = types.Must(types.IdFromString(claims.UserId))

        mocked.User.EXPECT().
          FindByIds(mock.Anything, a.input.Id).
            Return([]entity.User{
              {
                Id:       a.input.Id,
                Username: gofakeit.Username(),
                Email:    types.Email(gofakeit.Email()),
                Password: types.Must(a.input.LastPassword.Hash()),
              },
            }, nil)

        mocked.User.EXPECT().
          Patch(mock.Anything, mock.Anything).
          Return(nil)
      },
      args: args{
        ctx: generateClaimsCtx(constant.AUTHN_UPDATE_USER),
        input: &dto.UserUpdatePasswordDTO{
          Id:           types.MustCreateId(),
          LastPassword: types.Password(gofakeit.Username()),
          NewPassword:  types.Password(gofakeit.Username()),
        },
      },
      want: status.Updated(),
    },
    {
      name: "Success update other password",
      setup: func(mocked *userMocked, arg any, want any) {
        a := arg.(*args)

        mocked.defaultUOWMock()

        mocked.User.EXPECT().
          FindByIds(mock.Anything, a.input.Id).
            Return([]entity.User{
              {
                Id:       a.input.Id,
                Username: gofakeit.Username(),
                Email:    types.Email(gofakeit.Email()),
                Password: types.Must(a.input.LastPassword.Hash()),
              },
            }, nil)

        mocked.User.EXPECT().
          Patch(mock.Anything, mock.Anything).
          Return(nil)
      },
      args: args{
        ctx: generateClaimsCtx(constant.AUTHN_UPDATE_USER_ARB),
        input: &dto.UserUpdatePasswordDTO{
          Id:           types.MustCreateId(),
          LastPassword: types.Password(gofakeit.Username()),
          NewPassword:  types.Password(gofakeit.Username()),
        },
      },
      want: status.Updated(),
    },
    {
      name: "Failed to update other password due to permission",
      setup: func(mocked *userMocked, arg any, want any) {

      },
      args: args{
        ctx: generateClaimsCtx(constant.AUTHN_UPDATE_USER),
        input: &dto.UserUpdatePasswordDTO{
          Id:           types.MustCreateId(),
          LastPassword: types.Password(gofakeit.Username()),
          NewPassword:  types.Password(gofakeit.Username()),
        },
      },
      want: status.ErrUnAuthorized(sharedErr.ErrUnauthorizedPermission),
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newUserMocked(t)

      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, &tt.want)
      }

      u := userService{
        unit:     mocked.UOW,
        credRepo: mocked.Cred,
        tracer:   mocked.Tracer,
        config: UserConfig{
          RoleClient:  mocked.RoleClient,
          MailClient:  mocked.MailClient,
          TokenClient: mocked.TokenClient,
        },
      }

      if got := u.UpdatePassword(tt.args.ctx, tt.args.input); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("UpdatePassword() = %v, want %v", got, tt.want)
      }
    })
  }
}

func Test_userService_VerifyEmail(t *testing.T) {
  t.Parallel()

  type args struct {
    ctx   context.Context
    token string
  }
  tests := []struct {
    name  string
    setup setupUserTestFunc
    args  args
    want  status.Object
  }{
    {
      name: "Success verify email",
      setup: func(mocked *userMocked, arg any, want any) {
        a := arg.(*args)
        mocked.defaultUOWMock()

        mocked.TokenClient.EXPECT().
            Verify(mock.Anything, &dto.TokenVerificationDTO{
              Token:   a.token,
              Purpose: dto.TokenUsageEmailVerification,
            }).
          Return(dummyId, nil)

        mocked.User.EXPECT().
            Patch(mock.Anything, &entity.PatchedUser{
              Id:         dummyId,
              IsVerified: types.SomeNullable(true),
            }).Return(nil)
      },
      args: args{
        ctx:   context.Background(),
        token: sharedUtil.RandomString(32),
      },
      want: status.Updated(),
    },
    {
      name: "Failed verify email",
      setup: func(mocked *userMocked, arg any, want any) {
        a := arg.(*args)

        mocked.TokenClient.EXPECT().
            Verify(mock.Anything, &dto.TokenVerificationDTO{
              Token:   a.token,
              Purpose: dto.TokenUsageEmailVerification,
            }).
          Return(types.NullId(), dummyErr)
      },
      args: args{
        ctx:   context.Background(),
        token: sharedUtil.RandomString(32),
      },
      want: status.ErrExternal(dummyErr),
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newUserMocked(t)

      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, &tt.want)
      }

      u := userService{
        unit:     mocked.UOW,
        credRepo: mocked.Cred,
        tracer:   mocked.Tracer,
        config: UserConfig{
          RoleClient:  mocked.RoleClient,
          MailClient:  mocked.MailClient,
          TokenClient: mocked.TokenClient,
        },
      }

      if got := u.VerifyEmail(tt.args.ctx, tt.args.token); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("VerifyEmail() = %v, want %v", got, tt.want)
      }
    })
  }
}

func Test_userService_checkPermission(t *testing.T) {
  t.Parallel()

  type args struct {
    ctx         context.Context
    targetId    types.Id
    permissions string
  }
  tests := []struct {
    name    string
    setup   setupUserTestFunc
    args    args
    wantErr bool
  }{
    {
      name: "Success modify other",
      setup: func(mocked *userMocked, arg any, want any) {
      },
      args: args{
        ctx:         generateClaimsCtx(constant.AUTHN_UPDATE_USER, constant.AUTHN_UPDATE_USER_ARB),
        targetId:    types.MustCreateId(),
        permissions: constant.AUTHN_PERMISSIONS[constant.AUTHN_UPDATE_USER_ARB],
      },
      wantErr: false,
    },
    {
      name: "Failed modify other",
      setup: func(mocked *userMocked, arg any, want any) {
      },
      args: args{
        ctx:         generateClaimsCtx(constant.AUTHN_UPDATE_USER),
        targetId:    types.MustCreateId(),
        permissions: constant.AUTHN_PERMISSIONS[constant.AUTHN_UPDATE_USER_ARB],
      },
      wantErr: true,
    },
    {
      name: "Success self modify",
      setup: func(mocked *userMocked, arg any, want any) {
        a := arg.(*args)

        claims, err := sharedJwt.GetUserClaimsFromCtx(a.ctx)
        require.NoError(t, err)

        a.targetId = types.Must(types.IdFromString(claims.UserId))
      },
      args: args{
        ctx:         generateClaimsCtx(constant.AUTHN_UPDATE_USER),
        targetId:    types.Id{}, // Set on setup
        permissions: constant.AUTHN_PERMISSIONS[constant.AUTHN_UPDATE_USER_ARB],
      },
      wantErr: false,
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newUserMocked(t)

      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, nil)
      }

      u := userService{
        unit:     mocked.UOW,
        credRepo: mocked.Cred,
        tracer:   mocked.Tracer,
        config: UserConfig{
          RoleClient:  mocked.RoleClient,
          MailClient:  mocked.MailClient,
          TokenClient: mocked.TokenClient,
        },
      }

      if err := u.checkPermission(tt.args.ctx, tt.args.targetId, tt.args.permissions); (err != nil) != tt.wantErr {
        t.Errorf("checkPermission() error = %v, wantErr %v", err, tt.wantErr)
      }
    })
  }
}
