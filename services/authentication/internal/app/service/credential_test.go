package service

import (
  "context"
  "crypto/rand"
  "crypto/rsa"
  "database/sql"
  "fmt"
  sharedConst "github.com/arcorium/nexa/shared/constant"
  sharedErr "github.com/arcorium/nexa/shared/errors"
  sharedJwt "github.com/arcorium/nexa/shared/jwt"
  "github.com/arcorium/nexa/shared/optional"
  "github.com/arcorium/nexa/shared/status"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  "github.com/brianvoe/gofakeit/v7"
  "github.com/golang-jwt/jwt/v5"
  "github.com/stretchr/testify/mock"
  "github.com/stretchr/testify/require"
  "go.opentelemetry.io/otel/trace"
  "go.opentelemetry.io/otel/trace/noop"
  "nexa/services/authentication/constant"
  "nexa/services/authentication/internal/domain/dto"
  "nexa/services/authentication/internal/domain/entity"
  extMock "nexa/services/authentication/internal/domain/external/mocks"
  repoMock "nexa/services/authentication/internal/domain/repository/mocks"
  "nexa/services/authentication/util/errors"
  "reflect"
  "testing"
  "time"
)

func newCredentialMocked(t *testing.T) credentialMocked {
  // Tracer
  provider := noop.NewTracerProvider()

  privkey, err := rsa.GenerateKey(rand.Reader, 2048)
  require.NoError(t, err)

  return credentialMocked{
    Config: CredentialServiceConfig{
      SigningMethod:          jwt.SigningMethodRS256,
      AccessTokenExpiration:  time.Minute * 5,
      RefreshTokenExpiration: time.Hour * 24 * 7,
      PrivateKey:             privkey,
      PublicKey:              &privkey.PublicKey,
    },
    Cred:       repoMock.NewCredentialMock(t),
    Token:      repoMock.NewTokenMock(t),
    UserClient: extMock.NewUserClientMock(t),
    RoleClient: extMock.NewRoleClientMock(t),
    MailClient: extMock.NewMailClientMock(t),
    Tracer:     provider.Tracer("MOCK"),
  }
}

type credentialMocked struct {
  Config     CredentialServiceConfig
  Cred       *repoMock.CredentialMock
  Token      *repoMock.TokenMock
  UserClient *extMock.UserClientMock
  RoleClient *extMock.RoleClientMock
  MailClient *extMock.MailClientMock
  Tracer     trace.Tracer
}

type setupCredTestFunc func(mocked *credentialMocked, arg any, want any)

func Test_credentialService_GetCredentials(t *testing.T) {
  type args struct {
    ctx    context.Context
    userId types.Id
  }
  tests := []struct {
    name  string
    setup setupCredTestFunc
    args  args
    want  []dto.CredentialResponseDTO
    want1 status.Object
  }{
    {
      name: "Success get self user credentials",
      setup: func(mocked *credentialMocked, arg any, want any) {
        a := arg.(*args)
        w := want.([]dto.CredentialResponseDTO)

        claims, err := sharedJwt.GetUserClaimsFromCtx(a.ctx)
        require.NoError(t, err)

        a.userId = types.Must(types.IdFromString(claims.UserId))
        creds := sharedUtil.CastSliceP(w, func(cred *dto.CredentialResponseDTO) entity.Credential {
          return entity.Credential{
            Id:            cred.Id,
            UserId:        a.userId,
            AccessTokenId: types.MustCreateId(),
            Device:        entity.Device{Name: cred.Device},
            RefreshToken:  sharedUtil.RandomString(32),
            ExpiresAt:     time.Now().Add(time.Hour * 3),
          }
        })

        mocked.Cred.EXPECT().
          FindByUserId(mock.Anything, a.userId).
          Return(creds, nil)
      },
      args: args{
        ctx:    generateClaimsCtx(), // without roles
        userId: types.MustCreateId(),
      },
      want: sharedUtil.GenerateMultiple(3, func() dto.CredentialResponseDTO {
        return dto.CredentialResponseDTO{
          Id:     types.MustCreateId(),
          Device: sharedUtil.RandomString(12),
        }
      }),
      want1: status.Success(),
    },
    {
      name: "Success get other user credentials",
      setup: func(mocked *credentialMocked, arg any, want any) {
        a := arg.(*args)
        w := want.([]dto.CredentialResponseDTO)

        creds := sharedUtil.CastSliceP(w, func(cred *dto.CredentialResponseDTO) entity.Credential {
          return entity.Credential{
            Id:            cred.Id,
            UserId:        a.userId,
            AccessTokenId: types.MustCreateId(),
            Device:        entity.Device{Name: cred.Device},
            RefreshToken:  sharedUtil.RandomString(32),
            ExpiresAt:     time.Now().Add(time.Hour * 3),
          }
        })

        mocked.Cred.EXPECT().
          FindByUserId(mock.Anything, a.userId).
          Return(creds, nil)
      },
      args: args{
        ctx:    generateClaimsCtx(constant.AUTHN_GET_OTHER_CREDENTIALS), // without roles
        userId: types.MustCreateId(),
      },
      want: sharedUtil.GenerateMultiple(3, func() dto.CredentialResponseDTO {
        return dto.CredentialResponseDTO{
          Id:     types.MustCreateId(),
          Device: sharedUtil.RandomString(12),
        }
      }),
      want1: status.Success(),
    },
    {
      name: "Has no permission",
      setup: func(mocked *credentialMocked, arg any, want any) {
      },
      args: args{
        ctx:    generateClaimsCtx(), // without roles
        userId: types.MustCreateId(),
      },
      want:  nil,
      want1: status.ErrUnAuthorized(sharedErr.ErrUnauthorizedPermission),
    },
    {
      name: "User has no credentials",
      setup: func(mocked *credentialMocked, arg any, want any) {
        a := arg.(*args)

        mocked.Cred.EXPECT().
          FindByUserId(mock.Anything, a.userId).
          Return(nil, sql.ErrNoRows)
      },
      args: args{
        ctx:    generateClaimsCtx(constant.AUTHN_GET_OTHER_CREDENTIALS), // without roles
        userId: types.MustCreateId(),
      },
      want:  nil,
      want1: status.FromRepository(sql.ErrNoRows, status.NullCode),
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newCredentialMocked(t)

      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, tt.want)
      }

      c := &credentialService{
        credRepo:  mocked.Cred,
        tokenRepo: mocked.Token,
        userExt:   mocked.UserClient,
        roleExt:   mocked.RoleClient,
        mailExt:   mocked.MailClient,
        config:    mocked.Config,
        tracer:    mocked.Tracer,
      }

      got, got1 := c.GetCredentials(tt.args.ctx, tt.args.userId)
      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("GetCredentials() got = %v, want %v", got, tt.want)
      }
      if !reflect.DeepEqual(got1, tt.want1) {
        t.Errorf("GetCredentials() got1 = %v, want %v", got1, tt.want1)
      }
    })
  }
}

func Test_credentialService_Login(t *testing.T) {
  type args struct {
    ctx      context.Context
    loginDto *dto.LoginDTO
  }
  tests := []struct {
    name  string
    setup setupCredTestFunc
    args  args
    want  dto.LoginResponseDTO
    want1 status.Object
  }{
    {
      name: "Success login",
      setup: func(mocked *credentialMocked, arg any, want any) {
        a := arg.(*args)

        userResp := dto.UserResponseDTO{
          UserId:   types.MustCreateId(),
          Username: gofakeit.Username(),
        }

        rolesResp := sharedUtil.GenerateMultiple(3, func() dto.RoleResponseDTO {
          return dto.RoleResponseDTO{
            Id:   types.MustCreateId(),
            Role: gofakeit.AnimalType(),
            Permissions: sharedUtil.GenerateMultiple(2, func() dto.Permission {
              return dto.Permission{
                Id:   types.MustCreateId(),
                Code: fmt.Sprintf("%s:%s", gofakeit.AppName(), gofakeit.AnimalType()),
              }
            }),
          }
        })

        mocked.UserClient.EXPECT().
          Validate(mock.Anything, a.loginDto.Email, a.loginDto.Password).
          Return(userResp, nil)

        mocked.RoleClient.EXPECT().
          GetUserRoles(mock.Anything, userResp.UserId).
          Return(rolesResp, nil)

        mocked.Cred.EXPECT().
          Create(mock.Anything, mock.Anything).
          Return(nil)
      },
      args: args{
        ctx: context.Background(),
        loginDto: &dto.LoginDTO{
          Email:      types.Email(gofakeit.Email()),
          Password:   types.Password(sharedUtil.RandomString(12)),
          DeviceName: gofakeit.Username(),
        },
      },
      want: dto.LoginResponseDTO{
        TokenType: constant.TOKEN_TYPE,
      },
      want1: status.Success(),
    },
    {
      name: "Failed to validate",
      setup: func(mocked *credentialMocked, arg any, want any) {
        a := arg.(*args)

        mocked.UserClient.EXPECT().
          Validate(mock.Anything, a.loginDto.Email, a.loginDto.Password).
          Return(dto.UserResponseDTO{}, dummyErr)
      },
      args: args{
        ctx: context.Background(),
        loginDto: &dto.LoginDTO{
          Email:      types.Email(gofakeit.Email()),
          Password:   types.Password(sharedUtil.RandomString(12)),
          DeviceName: gofakeit.Username(),
        },
      },
      want: dto.LoginResponseDTO{
        TokenType: constant.TOKEN_TYPE,
      },
      want1: status.ErrExternal(dummyErr),
    },
    {
      name: "Failed to get roles",
      setup: func(mocked *credentialMocked, arg any, want any) {
        a := arg.(*args)

        userResp := dto.UserResponseDTO{
          UserId:   types.MustCreateId(),
          Username: gofakeit.Username(),
        }

        mocked.UserClient.EXPECT().
          Validate(mock.Anything, a.loginDto.Email, a.loginDto.Password).
          Return(userResp, nil)

        mocked.RoleClient.EXPECT().
          GetUserRoles(mock.Anything, userResp.UserId).
          Return(nil, dummyErr)
      },
      args: args{
        ctx: context.Background(),
        loginDto: &dto.LoginDTO{
          Email:      types.Email(gofakeit.Email()),
          Password:   types.Password(sharedUtil.RandomString(12)),
          DeviceName: gofakeit.Username(),
        },
      },
      want: dto.LoginResponseDTO{
        TokenType: constant.TOKEN_TYPE,
      },
      want1: status.ErrExternal(dummyErr),
    },
    {
      name: "Success login with no roles",
      setup: func(mocked *credentialMocked, arg any, want any) {
        a := arg.(*args)

        userResp := dto.UserResponseDTO{
          UserId:   types.MustCreateId(),
          Username: gofakeit.Username(),
        }

        mocked.UserClient.EXPECT().
          Validate(mock.Anything, a.loginDto.Email, a.loginDto.Password).
          Return(userResp, nil)

        mocked.RoleClient.EXPECT().
          GetUserRoles(mock.Anything, userResp.UserId).
          Return(nil, nil)

        mocked.Cred.EXPECT().
          Create(mock.Anything, mock.Anything).
          Return(nil)
      },
      args: args{
        ctx: context.Background(),
        loginDto: &dto.LoginDTO{
          Email:      types.Email(gofakeit.Email()),
          Password:   types.Password(sharedUtil.RandomString(12)),
          DeviceName: gofakeit.Username(),
        },
      },
      want: dto.LoginResponseDTO{
        TokenType: constant.TOKEN_TYPE,
      },
      want1: status.Success(),
    },
    {
      name: "Success login with roles without permission",
      setup: func(mocked *credentialMocked, arg any, want any) {
        a := arg.(*args)

        userResp := dto.UserResponseDTO{
          UserId:   types.MustCreateId(),
          Username: gofakeit.Username(),
        }

        rolesResp := sharedUtil.GenerateMultiple(3, func() dto.RoleResponseDTO {
          return dto.RoleResponseDTO{
            Id:          types.MustCreateId(),
            Role:        gofakeit.AnimalType(),
            Permissions: nil,
          }
        })

        mocked.UserClient.EXPECT().
          Validate(mock.Anything, a.loginDto.Email, a.loginDto.Password).
          Return(userResp, nil)

        mocked.RoleClient.EXPECT().
          GetUserRoles(mock.Anything, userResp.UserId).
          Return(rolesResp, nil)

        mocked.Cred.EXPECT().
          Create(mock.Anything, mock.Anything).
          Return(nil)
      },
      args: args{
        ctx: context.Background(),
        loginDto: &dto.LoginDTO{
          Email:      types.Email(gofakeit.Email()),
          Password:   types.Password(sharedUtil.RandomString(12)),
          DeviceName: gofakeit.Username(),
        },
      },
      want: dto.LoginResponseDTO{
        TokenType: constant.TOKEN_TYPE,
      },
      want1: status.Success(),
    },
    {
      name: "Failed to save to repository",
      setup: func(mocked *credentialMocked, arg any, want any) {
        a := arg.(*args)

        userResp := dto.UserResponseDTO{
          UserId:   types.MustCreateId(),
          Username: gofakeit.Username(),
        }

        rolesResp := sharedUtil.GenerateMultiple(3, func() dto.RoleResponseDTO {
          return dto.RoleResponseDTO{
            Id:   types.MustCreateId(),
            Role: gofakeit.AnimalType(),
            Permissions: sharedUtil.GenerateMultiple(2, func() dto.Permission {
              return dto.Permission{
                Id:   types.MustCreateId(),
                Code: fmt.Sprintf("%s:%s", gofakeit.AppName(), gofakeit.AnimalType()),
              }
            }),
          }
        })

        mocked.UserClient.EXPECT().
          Validate(mock.Anything, a.loginDto.Email, a.loginDto.Password).
          Return(userResp, nil)

        mocked.RoleClient.EXPECT().
          GetUserRoles(mock.Anything, userResp.UserId).
          Return(rolesResp, nil)

        mocked.Cred.EXPECT().
          Create(mock.Anything, mock.Anything).
          Return(sql.ErrNoRows)
      },
      args: args{
        ctx: context.Background(),
        loginDto: &dto.LoginDTO{
          Email:      types.Email(gofakeit.Email()),
          Password:   types.Password(sharedUtil.RandomString(12)),
          DeviceName: gofakeit.Username(),
        },
      },
      want: dto.LoginResponseDTO{
        TokenType: constant.TOKEN_TYPE,
      },
      want1: status.FromRepository(sql.ErrNoRows, status.NullCode),
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newCredentialMocked(t)

      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, &tt.want)
      }

      c := &credentialService{
        credRepo:  mocked.Cred,
        tokenRepo: mocked.Token,
        userExt:   mocked.UserClient,
        roleExt:   mocked.RoleClient,
        mailExt:   mocked.MailClient,
        config:    mocked.Config,
        tracer:    mocked.Tracer,
      }

      got, got1 := c.Login(tt.args.ctx, tt.args.loginDto)
      //if !reflect.DeepEqual(got, tt.want) {
      //  t.Errorf("Login() got = %v, want %v", got, tt.want)
      //}
      if !reflect.DeepEqual(got1, tt.want1) {
        t.Errorf("Login() got1 = %v, want %v", got1, tt.want1)
      }

      if got1.IsError() {
        return
      }

      require.Equal(t, got.TokenType, tt.want.TokenType)
      require.NotEmptyf(t, got.Token, "Token should not empty")
    })
  }
}

func Test_credentialService_Logout(t *testing.T) {
  type args struct {
    ctx       context.Context
    logoutDTO *dto.LogoutDTO
  }
  tests := []struct {
    name  string
    setup setupCredTestFunc
    args  args
    want  status.Object
  }{
    {
      name: "Success self logout",
      setup: func(mocked *credentialMocked, arg any, want any) {
        a := arg.(*args)

        credIds := sharedUtil.CastSlice(a.logoutDTO.CredentialIds, sharedUtil.ToAny[types.Id])

        mocked.Cred.EXPECT().
          DeleteByUserId(mock.Anything, a.logoutDTO.UserId, credIds...).
          Return(nil)
      },
      args: args{
        ctx: generateClaimsCtx(),
        logoutDTO: &dto.LogoutDTO{
          UserId:        claimsUserId,
          CredentialIds: sharedUtil.GenerateMultiple(3, types.MustCreateId),
        },
      },
      want: status.Deleted(),
    },
    {
      name: "Success logout other",
      setup: func(mocked *credentialMocked, arg any, want any) {
        a := arg.(*args)

        credIds := sharedUtil.CastSlice(a.logoutDTO.CredentialIds, sharedUtil.ToAny[types.Id])

        mocked.Cred.EXPECT().
          Delete(mock.Anything, credIds...).
          Return(nil)
      },
      args: args{
        ctx: generateClaimsCtx(constant.AUTHN_LOGOUT_OTHER),
        logoutDTO: &dto.LogoutDTO{
          UserId:        types.MustCreateId(),
          CredentialIds: sharedUtil.GenerateMultiple(3, types.MustCreateId),
        },
      },
      want: status.Deleted(),
    },
    {
      name: "Failed to delete credentials",
      setup: func(mocked *credentialMocked, arg any, want any) {
        a := arg.(*args)

        credIds := sharedUtil.CastSlice(a.logoutDTO.CredentialIds, sharedUtil.ToAny[types.Id])

        mocked.Cred.EXPECT().
          DeleteByUserId(mock.Anything, a.logoutDTO.UserId, credIds...).
          Return(sql.ErrNoRows)
      },
      args: args{
        ctx: generateClaimsCtx(),
        logoutDTO: &dto.LogoutDTO{
          UserId:        claimsUserId,
          CredentialIds: sharedUtil.GenerateMultiple(3, types.MustCreateId),
        },
      },
      want: status.FromRepository(sql.ErrNoRows, status.NullCode),
    },
    {
      name: "Failed to delete other users credentials",
      setup: func(mocked *credentialMocked, arg any, want any) {
        a := arg.(*args)

        credIds := sharedUtil.CastSlice(a.logoutDTO.CredentialIds, sharedUtil.ToAny[types.Id])

        mocked.Cred.EXPECT().
          Delete(mock.Anything, credIds...).
          Return(sql.ErrNoRows)
      },
      args: args{
        ctx: generateClaimsCtx(constant.AUTHN_LOGOUT_OTHER),
        logoutDTO: &dto.LogoutDTO{
          UserId:        types.MustCreateId(),
          CredentialIds: sharedUtil.GenerateMultiple(3, types.MustCreateId),
        },
      },
      want: status.FromRepository(sql.ErrNoRows, status.NullCode),
    },
    {
      name: "Has no permission to logout other user",
      setup: func(mocked *credentialMocked, arg any, want any) {

      },
      args: args{
        ctx: generateClaimsCtx(),
        logoutDTO: &dto.LogoutDTO{
          UserId:        types.MustCreateId(),
          CredentialIds: sharedUtil.GenerateMultiple(3, types.MustCreateId),
        },
      },
      want: status.ErrUnAuthorized(sharedErr.ErrUnauthorizedPermission),
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newCredentialMocked(t)

      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, &tt.want)
      }

      c := &credentialService{
        credRepo:  mocked.Cred,
        tokenRepo: mocked.Token,
        userExt:   mocked.UserClient,
        roleExt:   mocked.RoleClient,
        mailExt:   mocked.MailClient,
        config:    mocked.Config,
        tracer:    mocked.Tracer,
      }

      if got := c.Logout(tt.args.ctx, tt.args.logoutDTO); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("Logout() = %v, want %v", got, tt.want)
      }
    })
  }
}

func Test_credentialService_LogoutAll(t *testing.T) {
  type args struct {
    ctx    context.Context
    userId types.Id
  }
  tests := []struct {
    name  string
    setup setupCredTestFunc
    args  args
    want  status.Object
  }{
    {
      name: "Success logout all for self user",
      setup: func(mocked *credentialMocked, arg any, want any) {
        a := arg.(*args)

        mocked.Cred.EXPECT().
          DeleteByUserId(mock.Anything, a.userId).
          Return(nil)
      },
      args: args{
        ctx:    generateClaimsCtx(),
        userId: claimsUserId,
      },
      want: status.Deleted(),
    },
    {
      name: "Success logout all other users",
      setup: func(mocked *credentialMocked, arg any, want any) {
        a := arg.(*args)

        mocked.Cred.EXPECT().
          DeleteByUserId(mock.Anything, a.userId).
          Return(nil)
      },
      args: args{
        ctx:    generateClaimsCtx(constant.AUTHN_LOGOUT_OTHER),
        userId: types.MustCreateId(),
      },
      want: status.Deleted(),
    },
    {
      name: "Failed to logout other user due to no permission",
      setup: func(mocked *credentialMocked, arg any, want any) {
      },
      args: args{
        ctx:    generateClaimsCtx(),
        userId: types.MustCreateId(),
      },
      want: status.ErrUnAuthorized(sharedErr.ErrUnauthorizedPermission),
    },
    {
      name: "Failed to delete credentials from repository",
      setup: func(mocked *credentialMocked, arg any, want any) {
        a := arg.(*args)

        mocked.Cred.EXPECT().
          DeleteByUserId(mock.Anything, a.userId).
          Return(sql.ErrNoRows)
      },
      args: args{
        ctx:    generateClaimsCtx(),
        userId: claimsUserId,
      },
      want: status.FromRepository(sql.ErrNoRows, status.NullCode),
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newCredentialMocked(t)

      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, &tt.want)
      }

      c := &credentialService{
        credRepo:  mocked.Cred,
        tokenRepo: mocked.Token,
        userExt:   mocked.UserClient,
        roleExt:   mocked.RoleClient,
        mailExt:   mocked.MailClient,
        config:    mocked.Config,
        tracer:    mocked.Tracer,
      }

      if got := c.LogoutAll(tt.args.ctx, tt.args.userId); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("LogoutAll() = %v, want %v", got, tt.want)
      }
    })
  }
}

func Test_credentialService_RefreshToken(t *testing.T) {
  type args struct {
    ctx        context.Context
    refreshDto *dto.RefreshTokenDTO
  }
  tests := []struct {
    name  string
    setup setupCredTestFunc
    args  args
    want1 status.Object
  }{
    {
      name: "Success refresh token",
      setup: func(mocked *credentialMocked, arg any, want any) {
        a := arg.(*args)

        // Generate access token
        claims := generateUserClaims()
        token := jwt.NewWithClaims(mocked.Config.SigningMethod, claims)
        tokenStr, err := token.SignedString(mocked.Config.PrivateKey)
        require.NoError(t, err)

        a.refreshDto.AccessToken = tokenStr

        cred := entity.Credential{
          Id:            types.Must(types.IdFromString(claims.CredentialId)),
          UserId:        types.Must(types.IdFromString(claims.UserId)),
          AccessTokenId: types.Must(types.IdFromString(claims.ID)),
          Device: entity.Device{
            Name: gofakeit.AppName(),
          },
          RefreshToken: sharedUtil.RandomString(32),
          ExpiresAt:    time.Now().Add(time.Hour * 5),
        }

        rolesResp := sharedUtil.GenerateMultiple(3, func() dto.RoleResponseDTO {
          return dto.RoleResponseDTO{
            Id:   types.MustCreateId(),
            Role: gofakeit.AnimalType(),
            Permissions: sharedUtil.GenerateMultiple(2, func() dto.Permission {
              return dto.Permission{
                Id:   types.MustCreateId(),
                Code: fmt.Sprintf("%s:%s", gofakeit.AppName(), gofakeit.AnimalType()),
              }
            }),
          }
        })

        mocked.Cred.EXPECT().
          Find(mock.Anything, cred.Id).
          Return(&cred, nil)

        mocked.RoleClient.EXPECT().
          GetUserRoles(mock.Anything, cred.UserId).
          Return(rolesResp, nil)

        mocked.Cred.EXPECT().
          Patch(mock.Anything, mock.Anything).
          Return(nil)
      },
      args: args{
        ctx: context.Background(),
        refreshDto: &dto.RefreshTokenDTO{
          TokenType:   constant.TOKEN_TYPE,
          AccessToken: "",
        },
      },
      want1: status.Updated(),
    },
    {
      name: "Token has different scheme",
      setup: func(mocked *credentialMocked, arg any, want any) {
      },
      args: args{
        ctx: context.Background(),
        refreshDto: &dto.RefreshTokenDTO{
          TokenType:   "Other scheme",
          AccessToken: "",
        },
      },
      want1: status.ErrBadRequest(errors.ErrDifferentScheme),
    },
    {
      name: "Access token malformed",
      setup: func(mocked *credentialMocked, arg any, want any) {

      },
      args: args{
        ctx: context.Background(),
        refreshDto: &dto.RefreshTokenDTO{
          TokenType:   constant.TOKEN_TYPE,
          AccessToken: sharedUtil.RandomString(128),
        },
      },
      want1: status.ErrBadRequest(errors.ErrMalformedToken),
    },
    {
      name: "Access token has no refresh token",
      setup: func(mocked *credentialMocked, arg any, want any) {
        a := arg.(*args)

        // Generate access token
        claims := generateUserClaims()
        token := jwt.NewWithClaims(mocked.Config.SigningMethod, claims)
        tokenStr, err := token.SignedString(mocked.Config.PrivateKey)
        require.NoError(t, err)

        a.refreshDto.AccessToken = tokenStr

        cred := entity.Credential{
          Id:            types.Must(types.IdFromString(claims.CredentialId)),
          UserId:        types.Must(types.IdFromString(claims.UserId)),
          AccessTokenId: types.Must(types.IdFromString(claims.ID)),
          Device: entity.Device{
            Name: gofakeit.AppName(),
          },
          RefreshToken: sharedUtil.RandomString(32),
          ExpiresAt:    time.Now().Add(time.Hour * 5),
        }

        mocked.Cred.EXPECT().
          Find(mock.Anything, cred.Id).
          Return(nil, sql.ErrNoRows)
      },
      args: args{
        ctx: context.Background(),
        refreshDto: &dto.RefreshTokenDTO{
          TokenType:   constant.TOKEN_TYPE,
          AccessToken: "",
        },
      },
      want1: status.ErrBadRequest(errors.ErrRefreshTokenNotFound),
    },
    {
      name: "Failed to get user roles",
      setup: func(mocked *credentialMocked, arg any, want any) {
        a := arg.(*args)

        // Generate access token
        claims := generateUserClaims()
        token := jwt.NewWithClaims(mocked.Config.SigningMethod, claims)
        tokenStr, err := token.SignedString(mocked.Config.PrivateKey)
        require.NoError(t, err)

        a.refreshDto.AccessToken = tokenStr

        cred := entity.Credential{
          Id:            types.Must(types.IdFromString(claims.CredentialId)),
          UserId:        types.Must(types.IdFromString(claims.UserId)),
          AccessTokenId: types.Must(types.IdFromString(claims.ID)),
          Device: entity.Device{
            Name: gofakeit.AppName(),
          },
          RefreshToken: sharedUtil.RandomString(32),
          ExpiresAt:    time.Now().Add(time.Hour * 5),
        }

        mocked.Cred.EXPECT().
          Find(mock.Anything, cred.Id).
          Return(&cred, nil)

        mocked.RoleClient.EXPECT().
          GetUserRoles(mock.Anything, cred.UserId).
          Return(nil, dummyErr)
      },
      args: args{
        ctx: context.Background(),
        refreshDto: &dto.RefreshTokenDTO{
          TokenType:   constant.TOKEN_TYPE,
          AccessToken: "",
        },
      },
      want1: status.ErrExternal(dummyErr),
    },
    {
      name: "Failed to update data in repository",
      setup: func(mocked *credentialMocked, arg any, want any) {
        a := arg.(*args)

        // Generate access token
        claims := generateUserClaims()
        token := jwt.NewWithClaims(mocked.Config.SigningMethod, claims)
        tokenStr, err := token.SignedString(mocked.Config.PrivateKey)
        require.NoError(t, err)

        a.refreshDto.AccessToken = tokenStr

        cred := entity.Credential{
          Id:            types.Must(types.IdFromString(claims.CredentialId)),
          UserId:        types.Must(types.IdFromString(claims.UserId)),
          AccessTokenId: types.Must(types.IdFromString(claims.ID)),
          Device: entity.Device{
            Name: gofakeit.AppName(),
          },
          RefreshToken: sharedUtil.RandomString(32),
          ExpiresAt:    time.Now().Add(time.Hour * 5),
        }

        rolesResp := sharedUtil.GenerateMultiple(3, func() dto.RoleResponseDTO {
          return dto.RoleResponseDTO{
            Id:   types.MustCreateId(),
            Role: gofakeit.AnimalType(),
            Permissions: sharedUtil.GenerateMultiple(2, func() dto.Permission {
              return dto.Permission{
                Id:   types.MustCreateId(),
                Code: fmt.Sprintf("%s:%s", gofakeit.AppName(), gofakeit.AnimalType()),
              }
            }),
          }
        })

        mocked.Cred.EXPECT().
          Find(mock.Anything, cred.Id).
          Return(&cred, nil)

        mocked.RoleClient.EXPECT().
          GetUserRoles(mock.Anything, cred.UserId).
          Return(rolesResp, nil)

        mocked.Cred.EXPECT().
          Patch(mock.Anything, mock.Anything).
          Return(sql.ErrNoRows)
      },
      args: args{
        ctx: context.Background(),
        refreshDto: &dto.RefreshTokenDTO{
          TokenType:   constant.TOKEN_TYPE,
          AccessToken: "",
        },
      },
      want1: status.FromRepository(sql.ErrNoRows, optional.Some(status.INTERNAL_SERVER_ERROR)),
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newCredentialMocked(t)

      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, nil)
      }

      c := &credentialService{
        credRepo:  mocked.Cred,
        tokenRepo: mocked.Token,
        userExt:   mocked.UserClient,
        roleExt:   mocked.RoleClient,
        mailExt:   mocked.MailClient,
        config:    mocked.Config,
        tracer:    mocked.Tracer,
      }

      got, got1 := c.RefreshToken(tt.args.ctx, tt.args.refreshDto)

      if !reflect.DeepEqual(got1, tt.want1) {
        t.Errorf("RefreshToken() got1 = %v, want %v", got1, tt.want1)
      }

      if got1.IsError() {
        return
      }

      require.Equal(t, got.TokenType, constant.TOKEN_TYPE)
      require.NotEmpty(t, got.AccessToken)
    })
  }
}

func Test_credentialService_Register(t *testing.T) {
  type args struct {
    ctx         context.Context
    registerDTO *dto.RegisterDTO
  }
  tests := []struct {
    name  string
    setup setupCredTestFunc
    args  args
    want  status.Object
  }{
    {
      name: "Success register",
      setup: func(mocked *credentialMocked, arg any, want any) {
        a := arg.(*args)

        userId := types.MustCreateId()
        mocked.UserClient.EXPECT().
          Create(mock.Anything, a.registerDTO).
          Return(userId, nil)

        mocked.Token.EXPECT().
          Create(mock.Anything, mock.Anything).
          Return(nil)

        mocked.MailClient.EXPECT().
          Send(mock.Anything, mock.Anything).
          Return(nil)
      },
      args: args{
        ctx: context.Background(),
        registerDTO: &dto.RegisterDTO{
          Username:  gofakeit.Username(),
          Email:     types.Email(gofakeit.Email()),
          Password:  types.Password(sharedUtil.RandomString(12)),
          FirstName: gofakeit.FirstName(),
          LastName:  types.NullableString{},
          Bio:       types.NullableString{},
        },
      },
      want: status.Created(),
    },

    {
      name: "Failed to user",
      setup: func(mocked *credentialMocked, arg any, want any) {
        a := arg.(*args)

        mocked.UserClient.EXPECT().
          Create(mock.Anything, a.registerDTO).
          Return(types.NullId(), dummyErr)
      },
      args: args{
        ctx: context.Background(),
        registerDTO: &dto.RegisterDTO{
          Username:  gofakeit.Username(),
          Email:     types.Email(gofakeit.Email()),
          Password:  types.Password(sharedUtil.RandomString(12)),
          FirstName: gofakeit.FirstName(),
          LastName:  types.NullableString{},
          Bio:       types.NullableString{},
        },
      },
      want: status.ErrExternal(dummyErr),
    },
    {
      name: "Failed to create verification token",
      setup: func(mocked *credentialMocked, arg any, want any) {
        a := arg.(*args)

        userId := types.MustCreateId()
        mocked.UserClient.EXPECT().
          Create(mock.Anything, a.registerDTO).
          Return(userId, nil)

        mocked.Token.EXPECT().
          Create(mock.Anything, mock.Anything).
          Return(sql.ErrNoRows)
      },
      args: args{
        ctx: context.Background(),
        registerDTO: &dto.RegisterDTO{
          Username:  gofakeit.Username(),
          Email:     types.Email(gofakeit.Email()),
          Password:  types.Password(sharedUtil.RandomString(12)),
          FirstName: gofakeit.FirstName(),
          LastName:  types.NullableString{},
          Bio:       types.NullableString{},
        },
      },
      want: status.FromRepository(sql.ErrNoRows, status.NullCode),
    },
    {
      name: "Failed to send email", // Likely will not happen
      setup: func(mocked *credentialMocked, arg any, want any) {
        a := arg.(*args)

        userId := types.MustCreateId()
        mocked.UserClient.EXPECT().
          Create(mock.Anything, a.registerDTO).
          Return(userId, nil)

        mocked.Token.EXPECT().
          Create(mock.Anything, mock.Anything).
          Return(nil)

        mocked.MailClient.EXPECT().
          Send(mock.Anything, mock.Anything).
          Return(dummyErr)
      },
      args: args{
        ctx: context.Background(),
        registerDTO: &dto.RegisterDTO{
          Username:  gofakeit.Username(),
          Email:     types.Email(gofakeit.Email()),
          Password:  types.Password(sharedUtil.RandomString(12)),
          FirstName: gofakeit.FirstName(),
          LastName:  types.NullableString{},
          Bio:       types.NullableString{},
        },
      },
      want: status.ErrExternal(dummyErr),
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newCredentialMocked(t)

      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, &tt.want)
      }

      c := &credentialService{
        credRepo:  mocked.Cred,
        tokenRepo: mocked.Token,
        userExt:   mocked.UserClient,
        roleExt:   mocked.RoleClient,
        mailExt:   mocked.MailClient,
        config:    mocked.Config,
        tracer:    mocked.Tracer,
      }

      if got := c.Register(tt.args.ctx, tt.args.registerDTO); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("Register() = %v, want %v", got, tt.want)
      }
    })
  }
}

func Test_credentialService_checkPermission(t *testing.T) {
  type args struct {
    ctx        context.Context
    targetId   types.Id
    permission string
  }
  tests := []struct {
    name    string
    setup   setupCredTestFunc
    args    args
    wantErr bool
  }{
    {
      name: "Success modify other",
      setup: func(mocked *credentialMocked, arg any, want any) {
      },
      args: args{
        ctx:        generateClaimsCtx(constant.AUTHN_LOGOUT_OTHER, constant.AUTHN_GET_OTHER_CREDENTIALS),
        targetId:   types.MustCreateId(),
        permission: constant.AUTHN_PERMISSIONS[constant.AUTHN_LOGOUT_OTHER],
      },
      wantErr: false,
    },
    {
      name: "Failed modify other",
      setup: func(mocked *credentialMocked, arg any, want any) {
      },
      args: args{
        ctx:        generateClaimsCtx(constant.AUTHN_LOGOUT_OTHER),
        targetId:   types.MustCreateId(),
        permission: constant.AUTHN_PERMISSIONS[constant.AUTHN_GET_OTHER_CREDENTIALS],
      },
      wantErr: true,
    },
    {
      name: "Success self modify",
      setup: func(mocked *credentialMocked, arg any, want any) {
        a := arg.(*args)

        claims, err := sharedJwt.GetUserClaimsFromCtx(a.ctx)
        require.NoError(t, err)

        a.targetId = types.Must(types.IdFromString(claims.UserId))
      },
      args: args{
        ctx:        generateClaimsCtx(constant.AUTHN_GET_OTHER_CREDENTIALS),
        targetId:   types.Id{}, // Set on setup
        permission: constant.AUTHN_PERMISSIONS[constant.AUTHN_LOGOUT_OTHER],
      },
      wantErr: false,
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newCredentialMocked(t)

      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, nil)
      }

      c := &credentialService{
        credRepo:  mocked.Cred,
        tokenRepo: mocked.Token,
        userExt:   mocked.UserClient,
        roleExt:   mocked.RoleClient,
        mailExt:   mocked.MailClient,
        config:    mocked.Config,
        tracer:    mocked.Tracer,
      }

      if err := c.checkPermission(tt.args.ctx, tt.args.targetId, tt.args.permission); (err != nil) != tt.wantErr {
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
      return constant.AUTHN_PERMISSIONS[action]
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
    CredentialId: types.MustCreateId().String(),
    UserId:       claimsUserId.String(),
    Username:     gofakeit.Username(),
    Roles:        roles,
  }
}

func generateClaimsCtx(actions ...string) context.Context {
  return context.WithValue(context.Background(), sharedConst.USER_CLAIMS_CONTEXT_KEY, generateUserClaims(generateRole(actions...)))
}
