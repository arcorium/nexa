package service

import (
  "context"
  "database/sql"
  "errors"
  "github.com/brianvoe/gofakeit/v7"
  "github.com/golang-jwt/jwt/v5"
  "github.com/stretchr/testify/mock"
  "github.com/stretchr/testify/require"
  "go.opentelemetry.io/otel/trace"
  "go.opentelemetry.io/otel/trace/noop"
  "nexa/services/user/constant"
  userUow "nexa/services/user/internal/app/uow"
  "nexa/services/user/internal/domain/dto"
  "nexa/services/user/internal/domain/entity"
  extMock "nexa/services/user/internal/domain/external/mocks"
  "nexa/services/user/internal/domain/mapper"
  repoMock "nexa/services/user/internal/domain/repository/mocks"
  sharedConst "nexa/shared/constant"
  sharedDto "nexa/shared/dto"
  sharedErr "nexa/shared/errors"
  sharedJwt "nexa/shared/jwt"
  "nexa/shared/status"
  "nexa/shared/types"
  uowMock "nexa/shared/uow/mocks"
  sharedUtil "nexa/shared/util"
  "nexa/shared/util/repo"
  "reflect"
  "testing"
  "time"
)

var dummyErr = errors.New("dummy error")

var dummyId = types.MustCreateId()

func newUserMocked(t *testing.T) userMocked {
  // Tracer
  provider := noop.NewTracerProvider()
  return userMocked{
    UOW:     uowMock.NewUnitOfWorkMock[userUow.UserStorage](t),
    Profile: repoMock.NewProfileMock(t),
    User:    repoMock.NewUserMock(t),
    AuthN:   extMock.NewAuthenticationClientMock(t),
    Mailer:  extMock.NewMailerClientMock(t),
    Tracer:  provider.Tracer("MOCK"),
  }
}

type userMocked struct {
  UOW     *uowMock.UnitOfWorkMock[userUow.UserStorage]
  Profile *repoMock.ProfileMock
  User    *repoMock.UserMock
  AuthN   *extMock.AuthenticationClientMock
  Mailer  *extMock.MailerClientMock
  Tracer  trace.Tracer
}

func (m *userMocked) defaultUOWMock() {
  m.UOW.EXPECT().
    Repositories().
    Return(userUow.NewStorage(m.Profile, m.User))
}

type setupUserTestFunc func(mocked *userMocked, arg any, want any)

func Test_userService_BannedUser(t *testing.T) {
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

        mocked.UOW.
          EXPECT().
          Repositories().
          Return(userUow.NewStorage(mocked.Profile, mocked.User))

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

        mocked.UOW.
          EXPECT().
          Repositories().
          Return(userUow.NewStorage(mocked.Profile, mocked.User))

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
        unit:       mocked.UOW,
        tracer:     mocked.Tracer,
        mailClient: mocked.Mailer,
        authClient: mocked.AuthN,
      }

      if got := u.BannedUser(tt.args.ctx, tt.args.bannedDto); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("BannedUser() = %v, want %v", got, tt.want)
      }
    })
  }
}

func Test_userService_Create(t *testing.T) {
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
        unit:       mocked.UOW,
        tracer:     mocked.Tracer,
        mailClient: mocked.Mailer,
        authClient: mocked.AuthN,
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
  self := generateUserClaims(generateRole(constant.USER_DELETE))

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

        mocked.UOW.EXPECT().
          Repositories().
          Return(userUow.NewStorage(mocked.Profile, mocked.User))

        mocked.User.EXPECT().
          Delete(mock.Anything, a.userId).
          Return(nil)

        mocked.AuthN.EXPECT().
          DeleteCredentials(mock.Anything, mock.Anything).
          Return(nil)
      },
      args: args{
        userId: types.Must(types.IdFromString(self.UserId)),
        ctx: context.WithValue(context.Background(), sharedConst.CLAIMS_CONTEXT_KEY,
          self),
      },
      want: status.Deleted(),
    },
    {
      name: "Failed delete other user",
      setup: func(mocked *userMocked, arg any, want any) {
        a := arg.(*args)

        mocked.UOW.EXPECT().
          Repositories().
          Return(userUow.NewStorage(mocked.Profile, mocked.User))

        mocked.User.EXPECT().
          Delete(mock.Anything, a.userId).
          Return(sql.ErrNoRows)

      },
      args: args{
        userId: types.MustCreateId(),
        ctx: context.WithValue(context.Background(), sharedConst.CLAIMS_CONTEXT_KEY,
          generateUserClaims(generateRole(constant.USER_DELETE, constant.USER_DELETE_OTHER))),
      },
      want: status.FromRepository(sql.ErrNoRows, status.NullCode),
    },
    {
      name: "Success delete other user",
      setup: func(mocked *userMocked, arg any, want any) {
        a := arg.(*args)

        mocked.UOW.EXPECT().
          Repositories().
          Return(userUow.NewStorage(mocked.Profile, mocked.User))

        mocked.User.EXPECT().
          Delete(mock.Anything, a.userId).
          Return(nil)

        mocked.AuthN.EXPECT().
          DeleteCredentials(mock.Anything, mock.Anything).
          Return(nil)
      },
      args: args{
        userId: types.MustCreateId(),
        ctx: context.WithValue(context.Background(), sharedConst.CLAIMS_CONTEXT_KEY,
          generateUserClaims(generateRole(constant.USER_DELETE, constant.USER_DELETE_OTHER))),
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
        unit:       mocked.UOW,
        tracer:     mocked.Tracer,
        mailClient: mocked.Mailer,
        authClient: mocked.AuthN,
      }

      if got := u.DeleteById(tt.args.ctx, tt.args.userId); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("DeleteById() = %v, want %v", got, tt.want)
      }
    })
  }
}

func Test_userService_EmailVerificationRequest(t *testing.T) {
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
        claims, _ := sharedJwt.GetClaimsFromCtx(a.ctx)
        userId := types.Must(types.IdFromString(claims.UserId))

        mocked.AuthN.EXPECT().
            GenerateToken(mock.Anything, &dto.TokenGenerationDTO{
              UserId:  userId,
              Purpose: dto.EmailVerificationToken,
              TTL:     constant.EMAIL_VERIFICAITON_TOKEN_TTL,
            }).
            Return(dto.TokenResponseDTO{
              Token:     sharedUtil.RandomString(32),
              Purpose:   dto.EmailVerificationToken,
              ExpiredAt: time.Now().Add(constant.EMAIL_VERIFICAITON_TOKEN_TTL),
            }, nil)

        mocked.UOW.EXPECT().
          Repositories().
          Return(userUow.NewStorage(mocked.Profile, mocked.User))

        mocked.User.EXPECT().
          FindByIds(mock.Anything, userId).
          Return([]entity.User{{}}, nil)

        mocked.Mailer.EXPECT().
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
        claims, _ := sharedJwt.GetClaimsFromCtx(a.ctx)
        userId := types.Must(types.IdFromString(claims.UserId))

        mocked.UOW.EXPECT().
          Repositories().
          Return(userUow.NewStorage(mocked.Profile, mocked.User))

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
        claims, _ := sharedJwt.GetClaimsFromCtx(a.ctx)
        userId := types.Must(types.IdFromString(claims.UserId))

        mocked.AuthN.EXPECT().
            GenerateToken(mock.Anything, &dto.TokenGenerationDTO{
              UserId:  userId,
              Purpose: dto.EmailVerificationToken,
              TTL:     constant.EMAIL_VERIFICAITON_TOKEN_TTL,
            }).
          Return(dto.TokenResponseDTO{}, dummyErr)

        mocked.UOW.EXPECT().
          Repositories().
          Return(userUow.NewStorage(mocked.Profile, mocked.User))

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
        unit:       mocked.UOW,
        tracer:     mocked.Tracer,
        mailClient: mocked.Mailer,
        authClient: mocked.AuthN,
      }

      got, got1 := u.EmailVerificationRequest(tt.args.ctx)

      if !got1.IsError() {
        require.Greater(t, len(got.Token), 0)
        require.True(t, got.Purpose == dto.EmailVerificationToken)
        require.True(t, got.ExpiredAt.After(time.Now()))
      }

      if !reflect.DeepEqual(got1, tt.want1) {
        t.Errorf("EmailVerificationRequest() got1 = %v, want %v", got1, tt.want1)
      }
    })
  }
}

func Test_userService_FindAll(t *testing.T) {
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
        unit:       mocked.UOW,
        tracer:     mocked.Tracer,
        mailClient: mocked.Mailer,
        authClient: mocked.AuthN,
      }

      got, got1 := u.FindAll(tt.args.ctx, tt.args.pagedDto)
      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("FindAll() got = %v, want %v", got, tt.want)
      }
      if !reflect.DeepEqual(got1, tt.want1) {
        t.Errorf("FindAll() got1 = %v, want %v", got1, tt.want1)
      }
    })
  }
}

func Test_userService_FindByIds(t *testing.T) {
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
        unit:       mocked.UOW,
        tracer:     mocked.Tracer,
        mailClient: mocked.Mailer,
        authClient: mocked.AuthN,
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

        mocked.AuthN.EXPECT().
            GenerateToken(mock.Anything, &dto.TokenGenerationDTO{
              UserId:  user.Id,
              Purpose: dto.ForgotPasswordToken,
              TTL:     constant.FORGOT_PASSWORD_TOKEN_TTL,
            }).Return(*w, nil)

        mocked.Mailer.EXPECT().
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
        Purpose:   dto.ForgotPasswordToken,
        ExpiredAt: time.Now().Add(constant.FORGOT_PASSWORD_TOKEN_TTL),
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

        mocked.AuthN.EXPECT().
            GenerateToken(mock.Anything, &dto.TokenGenerationDTO{
              UserId:  user.Id,
              Purpose: dto.ForgotPasswordToken,
              TTL:     constant.FORGOT_PASSWORD_TOKEN_TTL,
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
        unit:       mocked.UOW,
        tracer:     mocked.Tracer,
        mailClient: mocked.Mailer,
        authClient: mocked.AuthN,
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

        mocked.AuthN.EXPECT().
            VerifyToken(mock.Anything, &dto.TokenVerificationDTO{
              Token:   a.input.Token.RawValue(),
              Purpose: dto.ForgotPasswordToken,
            }).Return(dummyId, nil)

        mocked.User.EXPECT().
          Patch(mock.Anything, mock.Anything).Return(nil)

        mocked.AuthN.EXPECT().
          DeleteCredentials(mock.Anything, dummyId).
          Return(nil)
      },
      args: args{
        ctx: context.Background(),
        input: &dto.ResetUserPasswordDTO{
          Token:       types.SomeNullable(sharedUtil.RandomString(32)),
          LogoutAll:   true,
          NewPassword: types.Password(gofakeit.Username()),
        },
      },
      want: status.Updated(),
    },
    {
      name: "Failed reset password with token due to verification",
      setup: func(mocked *userMocked, arg any, want any) {
        a := arg.(*args)
        mocked.AuthN.EXPECT().
            VerifyToken(mock.Anything, &dto.TokenVerificationDTO{
              Token:   a.input.Token.RawValue(),
              Purpose: dto.ForgotPasswordToken,
            }).Return(types.NullId(), dummyErr)

      },
      args: args{
        ctx: generateClaimsCtx(constant.USER_UPDATE_OTHER),
        input: &dto.ResetUserPasswordDTO{
          Token:       types.SomeNullable(sharedUtil.RandomString(32)),
          LogoutAll:   true,
          NewPassword: types.Password(gofakeit.Username()),
        },
      },
      want: status.ErrExternal(dummyErr),
    },
    {
      name: "Success reset password without token",
      setup: func(mocked *userMocked, arg any, want any) {
        a := arg.(*args)
        mocked.defaultUOWMock()

        mocked.User.EXPECT().
          Patch(mock.Anything, mock.Anything).Return(nil)

        mocked.AuthN.EXPECT().
          DeleteCredentials(mock.Anything, a.input.UserId.RawValue()).
          Return(nil)
      },
      args: args{
        ctx: generateClaimsCtx(constant.USER_UPDATE_OTHER),
        input: &dto.ResetUserPasswordDTO{
          //Token:       types.SomeNullable(sharedUtil.RandomString(32)),
          UserId:      types.SomeNullable(types.MustCreateId()),
          LogoutAll:   true,
          NewPassword: types.Password(gofakeit.Username()),
        },
      },
      want: status.Updated(),
    },
    {
      name: "Failed reset password due to user id doesn't present",
      setup: func(mocked *userMocked, arg any, want any) {
        //a := arg.(*args)
        //mocked.defaultUOWMock()

      },
      args: args{
        ctx: generateClaimsCtx(constant.USER_UPDATE_OTHER),
        input: &dto.ResetUserPasswordDTO{
          //UserId:      types.MustCreateId(),
          LogoutAll:   true,
          NewPassword: types.Password(gofakeit.Username()),
        },
      },
      want: status.ErrBadRequest(errors.New("expected user id, when there are no token provided")),
    },
    {
      name: "Failed reset other password due to permissions",
      setup: func(mocked *userMocked, arg any, want any) {

      },
      args: args{
        ctx: generateClaimsCtx(constant.USER_UPDATE),
        input: &dto.ResetUserPasswordDTO{
          UserId:      types.SomeNullable(types.MustCreateId()),
          LogoutAll:   true,
          NewPassword: types.Password(gofakeit.Username()),
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
        unit:       mocked.UOW,
        tracer:     mocked.Tracer,
        mailClient: mocked.Mailer,
        authClient: mocked.AuthN,
      }

      if got := u.ResetPassword(tt.args.ctx, tt.args.input); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("ResetPassword() = %v, want %v", got, tt.want)
      }
    })
  }
}

func Test_userService_Update(t *testing.T) {
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

        claims, err := sharedJwt.GetClaimsFromCtx(a.ctx)
        require.NoError(t, err)

        a.input.Id = types.Must(types.IdFromString(claims.UserId))

        ent := a.input.ToDomain()
        mocked.User.EXPECT().
          Patch(mock.Anything, &ent).
          Return(nil)

      },
      args: args{
        ctx: generateClaimsCtx(constant.USER_GET),
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
        ctx: generateClaimsCtx(constant.USER_UPDATE_OTHER),
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
        ctx: generateClaimsCtx(constant.USER_GET),
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
        unit:       mocked.UOW,
        tracer:     mocked.Tracer,
        mailClient: mocked.Mailer,
        authClient: mocked.AuthN,
      }

      if got := u.Update(tt.args.ctx, tt.args.input); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("Update() = %v, want %v", got, tt.want)
      }
    })
  }
}

func Test_userService_UpdatePassword(t *testing.T) {
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

        claims, err := sharedJwt.GetClaimsFromCtx(a.ctx)
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
        ctx: generateClaimsCtx(constant.USER_UPDATE),
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
        ctx: generateClaimsCtx(constant.USER_UPDATE_OTHER),
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
        ctx: generateClaimsCtx(constant.USER_UPDATE),
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
        unit:       mocked.UOW,
        tracer:     mocked.Tracer,
        mailClient: mocked.Mailer,
        authClient: mocked.AuthN,
      }

      if got := u.UpdatePassword(tt.args.ctx, tt.args.input); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("UpdatePassword() = %v, want %v", got, tt.want)
      }
    })
  }
}

func Test_userService_Validate(t *testing.T) {
  type args struct {
    ctx      context.Context
    email    types.Email
    password types.Password
  }
  tests := []struct {
    name  string
    setup setupUserTestFunc
    args  args
    want  dto.UserResponseDTO
    want1 status.Object
  }{
    {
      name: "Success validate",
      setup: func(mocked *userMocked, arg any, want any) {
        a := arg.(*args)
        w := want.(*dto.UserResponseDTO)

        mocked.defaultUOWMock()

        user := entity.User{
          Id:         w.Id,
          Username:   w.Username,
          Email:      w.Email,
          Password:   types.Must(a.password.Hash()),
          IsVerified: w.IsVerified,
        }

        mocked.User.EXPECT().
          FindByEmails(mock.Anything, a.email).
          Return([]entity.User{user}, nil)
      },
      args: args{
        ctx:      context.Background(),
        email:    types.Email(gofakeit.Email()),
        password: types.Password(gofakeit.Username()),
      },
      want: dto.UserResponseDTO{
        Id:       types.MustCreateId(),
        Username: gofakeit.Username(),
        Email:    types.Email(gofakeit.Email()),
      },
      want1: status.Success(),
    },
    {
      name: "Failed validate",
      setup: func(mocked *userMocked, arg any, want any) {
        a := arg.(*args)

        mocked.defaultUOWMock()

        mocked.User.EXPECT().
          FindByEmails(mock.Anything, a.email).
          Return(nil, sql.ErrNoRows)
      },
      args: args{
        ctx:      context.Background(),
        email:    types.Email(gofakeit.Email()),
        password: types.Password(gofakeit.Username()),
      },
      want:  dto.UserResponseDTO{},
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
        unit:       mocked.UOW,
        tracer:     mocked.Tracer,
        mailClient: mocked.Mailer,
        authClient: mocked.AuthN,
      }

      got, got1 := u.Validate(tt.args.ctx, tt.args.email, tt.args.password)
      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("Validate() got = %v, want %v", got, tt.want)
      }
      if !reflect.DeepEqual(got1, tt.want1) {
        t.Errorf("Validate() got1 = %v, want %v", got1, tt.want1)
      }
    })
  }
}

func Test_userService_VerifyEmail(t *testing.T) {
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

        mocked.AuthN.EXPECT().
            VerifyToken(mock.Anything, &dto.TokenVerificationDTO{
              Token:   a.token,
              Purpose: dto.EmailVerificationToken,
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

        mocked.AuthN.EXPECT().
            VerifyToken(mock.Anything, &dto.TokenVerificationDTO{
              Token:   a.token,
              Purpose: dto.EmailVerificationToken,
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
        unit:       mocked.UOW,
        tracer:     mocked.Tracer,
        mailClient: mocked.Mailer,
        authClient: mocked.AuthN,
      }

      if got := u.VerifyEmail(tt.args.ctx, tt.args.token); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("VerifyEmail() = %v, want %v", got, tt.want)
      }
    })
  }
}

func Test_userService_checkPermission(t *testing.T) {
  type args struct {
    ctx         context.Context
    targetId    types.Id
    permissions []string
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
        ctx:         generateClaimsCtx(constant.USER_UPDATE, constant.USER_UPDATE_OTHER),
        targetId:    types.MustCreateId(),
        permissions: []string{constant.USER_PERMISSIONS[constant.USER_UPDATE_OTHER]},
      },
      wantErr: false,
    },
    {
      name: "Failed modify other",
      setup: func(mocked *userMocked, arg any, want any) {
      },
      args: args{
        ctx:         generateClaimsCtx(constant.USER_UPDATE),
        targetId:    types.MustCreateId(),
        permissions: []string{constant.USER_PERMISSIONS[constant.USER_UPDATE_OTHER]},
      },
      wantErr: true,
    },
    {
      name: "Success self modify",
      setup: func(mocked *userMocked, arg any, want any) {
        a := arg.(*args)

        claims, err := sharedJwt.GetClaimsFromCtx(a.ctx)
        require.NoError(t, err)

        a.targetId = types.Must(types.IdFromString(claims.UserId))
      },
      args: args{
        ctx:         generateClaimsCtx(constant.USER_UPDATE),
        targetId:    types.Id{}, // Set on setup
        permissions: []string{constant.USER_PERMISSIONS[constant.USER_UPDATE_OTHER]},
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
        unit:       mocked.UOW,
        tracer:     mocked.Tracer,
        mailClient: mocked.Mailer,
        authClient: mocked.AuthN,
      }

      if err := u.checkPermission(tt.args.ctx, tt.args.targetId, tt.args.permissions...); (err != nil) != tt.wantErr {
        t.Errorf("checkPermission() error = %v, wantErr %v", err, tt.wantErr)
      }
    })
  }
}

func generateRole(actions ...string) sharedJwt.Role {
  return sharedJwt.Role{
    Id:   types.MustCreateId().String(),
    Role: gofakeit.AnimalType(),
    Permissions: sharedUtil.CastSlice(actions, func(action string) string {
      return constant.USER_PERMISSIONS[action]
    }),
  }
}

func generateUserClaims(roles ...sharedJwt.Role) *sharedJwt.UserClaims {
  return &sharedJwt.UserClaims{
    RegisteredClaims: jwt.RegisteredClaims{
      Issuer:    gofakeit.AppName(),
      Subject:   gofakeit.AppName(),
      ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 12)),
      NotBefore: jwt.NewNumericDate(time.Now()),
      IssuedAt:  jwt.NewNumericDate(time.Now()),
      ID:        types.MustCreateId().String(),
    },
    RefreshTokenId: types.MustCreateId().String(),
    UserId:         types.MustCreateId().String(),
    Username:       gofakeit.Username(),
    Roles:          roles,
  }
}

func generateClaimsCtx(actions ...string) context.Context {
  return context.WithValue(context.Background(), sharedConst.CLAIMS_CONTEXT_KEY, generateUserClaims(generateRole(actions...)))
}
