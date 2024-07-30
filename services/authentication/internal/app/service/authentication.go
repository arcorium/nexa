package service

import (
  "context"
  "crypto/rsa"
  "database/sql"
  "errors"
  "fmt"
  sharedErr "github.com/arcorium/nexa/shared/errors"
  sharedJwt "github.com/arcorium/nexa/shared/jwt"
  "github.com/arcorium/nexa/shared/optional"
  "github.com/arcorium/nexa/shared/status"
  "github.com/arcorium/nexa/shared/types"
  "github.com/arcorium/nexa/shared/uow"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  authUtil "github.com/arcorium/nexa/shared/util/auth"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "github.com/golang-jwt/jwt/v5"
  "go.opentelemetry.io/otel/trace"
  "nexa/services/authentication/constant"
  userUow "nexa/services/authentication/internal/app/uow"
  "nexa/services/authentication/internal/domain/dto"
  "nexa/services/authentication/internal/domain/entity"
  "nexa/services/authentication/internal/domain/external"
  "nexa/services/authentication/internal/domain/mapper"
  "nexa/services/authentication/internal/domain/repository"
  "nexa/services/authentication/internal/domain/service"
  "nexa/services/authentication/util"
  "nexa/services/authentication/util/errs"
  "time"
)

func NewCredential(credential repository.ICredential, unit uow.IUnitOfWork[userUow.UserStorage], config CredentialConfig) service.IAuthentication {
  return &credentialService{
    credRepo: credential,
    unit:     unit,
    config:   config,
    tracer:   util.GetTracer(),
  }
}

type CredentialConfig struct {
  TokenClient external.ITokenClient
  MailClient  external.IMailerClient
  RoleClient  external.IRoleClient

  SigningMethod          jwt.SigningMethod
  AccessTokenExpiration  time.Duration
  RefreshTokenExpiration time.Duration
  PrivateKey             *rsa.PrivateKey
  PublicKey              *rsa.PublicKey
}

type credentialService struct {
  credRepo repository.ICredential
  unit     uow.IUnitOfWork[userUow.UserStorage]

  config CredentialConfig
  tracer trace.Tracer
}

func (c *credentialService) checkPermission(ctx context.Context, targetId types.Id, permissions string) error {
  // Validate permission
  claims, _ := sharedJwt.GetUserClaimsFromCtx(ctx)
  if !targetId.EqWithString(claims.UserId) {
    // Need permission to update other users
    if !authUtil.ContainsPermission(claims.Roles, permissions) {
      return sharedErr.ErrUnauthorizedPermission
    }
  }
  return nil
}

func (c *credentialService) getUserRoles(ctx context.Context, userId types.Id) ([]sharedJwt.Role, error) {
  roles, err := c.config.RoleClient.GetUserRoles(ctx, userId)
  if err != nil {
    return nil, err
  }

  jwtRoles := sharedUtil.CastSliceP(roles, func(from *dto.RoleResponseDTO) sharedJwt.Role {
    return from.ToJWT()
  })
  return jwtRoles, nil
}

func (c *credentialService) getTokenClaims(tokenStr string) (*sharedJwt.UserClaims, error) {
  token, err := jwt.ParseWithClaims(tokenStr, &sharedJwt.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
    return c.config.PublicKey, nil
  })

  // Allow expired token
  if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
    return nil, errs.ErrMalformedToken
  }

  claims, ok := token.Claims.(*sharedJwt.UserClaims)
  if !ok {
    return nil, errs.ErrMalformedToken
  }
  return claims, nil
}

func (c *credentialService) Register(ctx context.Context, registerDTO *dto.RegisterDTO) status.Object {
  ctx, span := c.tracer.Start(ctx, "CredentialService.Register")
  defer span.End()

  // Create user
  user, profile, err := registerDTO.ToDomain()
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrInternal(err)
  }

  isRep := true
  err = c.unit.DoTx(ctx, func(ctx context.Context, storage userUow.UserStorage) error {
    err := storage.User().Create(ctx, &user)
    if err != nil {
      return err
    }

    err = storage.Profile().Create(ctx, &profile)
    if err != nil {
      return err
    }

    // Add user as default roles
    err = c.config.RoleClient.SetUserAsDefault(ctx, user.Id)
    if err != nil {
      isRep = false
    }
    return err
  })
  if err != nil {
    spanUtil.RecordError(err, span)
    if isRep {
      return status.FromRepository2(err, status.Null, status.Null)
    }
    return status.ErrExternal(err)
  }

  // Create verification token
  genDTO := dto.TokenGenerationDTO{
    UserId: user.Id,
    Usage:  dto.TokenUsageEmailVerification,
  }

  result, err := c.config.TokenClient.Generate(ctx, &genDTO)
  if err != nil {
    // Ignore email verification when there is an error on generating token
    span.RecordError(err)
    return status.Created()
  }

  // Send email verification
  mailDTO := dto.SendEmailVerificationDTO{
    Recipient: registerDTO.Email,
    Token:     result.Token,
  }
  err = c.config.MailClient.SendEmailVerification(ctx, &mailDTO)
  if err != nil {
    span.RecordError(err)
  }

  return status.Created()
}

func (c *credentialService) Login(ctx context.Context, loginDto *dto.LoginDTO) (dto.LoginResponseDTO, status.Object) {
  ctx, span := c.tracer.Start(ctx, "CredentialService.Login")
  defer span.End()

  // Get user by the email and validate the password
  repos := c.unit.Repositories()
  users, err := repos.User().FindByEmails(ctx, loginDto.Email)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.LoginResponseDTO{}, status.FromRepository(err, status.NullCode)
  }
  // Validate user password
  if err = users[0].ValidatePassword(loginDto.Password); err != nil {
    spanUtil.RecordError(err, span)
    return dto.LoginResponseDTO{}, status.ErrBadRequest(err)
  }
  // Check if user banned
  if users[0].IsBanned() {
    err = fmt.Errorf("user is currently banned until: %s", users[0].BannedUntil)
    spanUtil.RecordError(err, span)
    return dto.LoginResponseDTO{}, status.ErrBadRequest(err)
  }

  // Get user roles and permission
  jwtRoles, err := c.getUserRoles(ctx, users[0].Id)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.LoginResponseDTO{}, status.ErrExternal(err)
  }

  // Generate token pairs
  pairTokens, err := c.config.generatePairTokens(users[0].Username, users[0].Id, jwtRoles)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.LoginResponseDTO{}, status.ErrInternal(err)
  }

  // Save credentials
  credential := loginDto.ToDomain(users[0].Id, pairTokens.Access.Id, &pairTokens.Refresh, c.config.RefreshTokenExpiration)
  err = c.credRepo.Create(ctx, &credential)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.LoginResponseDTO{}, status.FromRepositoryExist(err)
  }

  return dto.LoginResponseDTO{
    TokenType:  constant.TOKEN_TYPE,
    Token:      pairTokens.Access.Token,
    ExpiryTime: c.config.AccessTokenExpiration,
  }, status.Success()
}

func (c *credentialService) RefreshToken(ctx context.Context, refreshDto *dto.RefreshTokenDTO) (dto.RefreshTokenResponseDTO, status.Object) {
  ctx, span := c.tracer.Start(ctx, "CredentialService.RefreshToken")
  defer span.End()

  // Check scheme
  if refreshDto.TokenType != constant.TOKEN_TYPE {
    spanUtil.RecordError(errs.ErrDifferentScheme, span)
    return dto.RefreshTokenResponseDTO{}, status.ErrBadRequest(errs.ErrDifferentScheme)
  }

  // Parse token
  claims, err := c.getTokenClaims(refreshDto.AccessToken)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.RefreshTokenResponseDTO{}, status.ErrBadRequest(err)
  }

  credId, err := types.IdFromString(claims.CredentialId)
  if err != nil {
    spanUtil.RecordError(errs.ErrMalformedToken, span)
    return dto.RefreshTokenResponseDTO{}, status.ErrBadRequest(errs.ErrMalformedToken)
  }

  cred, err := c.credRepo.Find(ctx, credId)
  if err != nil {
    spanUtil.RecordError(err, span)
    if err == sql.ErrNoRows {
      return dto.RefreshTokenResponseDTO{}, status.ErrBadRequest(errs.ErrRefreshTokenNotFound)
    }
    return dto.RefreshTokenResponseDTO{}, status.FromRepository(err, optional.Some(status.BAD_REQUEST_ERROR))
  }

  // Check relation
  if !cred.UserId.EqWithString(claims.UserId) || !cred.AccessTokenId.EqWithString(claims.ID) {
    spanUtil.RecordError(errs.ErrBadRelation, span)
    return dto.RefreshTokenResponseDTO{}, status.ErrBadRequest(errs.ErrBadRelation)
  }

  // Check if user still exist
  repos := c.unit.Repositories()
  _, err = repos.User().FindByIds(ctx, cred.UserId)
  if err != nil {
    // Delete refresh token when the user doesn't exist (could be deleted)
    _ = c.credRepo.Delete(ctx, credId)
    spanUtil.RecordError(err, span)
    return dto.RefreshTokenResponseDTO{}, status.ErrBadRequest(errs.ErrTokenBelongToNothing)
  }

  // Get user roles and permission
  jwtRoles, err := c.getUserRoles(ctx, cred.UserId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.RefreshTokenResponseDTO{}, status.ErrExternal(err)
  }

  // Create new access token
  accessToken, err := c.config.generateAccessToken(claims.Username, cred.UserId, cred.Id, jwtRoles)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.RefreshTokenResponseDTO{}, status.ErrInternal(err)
  }

  // Patch credentials with new access token
  updateCred := refreshDto.ToDomain(cred.Id, accessToken.Id)
  err = c.credRepo.Patch(ctx, &updateCred)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.RefreshTokenResponseDTO{}, status.FromRepository(err, optional.Some(status.INTERNAL_SERVER_ERROR))
  }

  response := dto.RefreshTokenResponseDTO{
    TokenType:   constant.TOKEN_TYPE,
    AccessToken: accessToken.Token,
    ExpiryTime:  c.config.AccessTokenExpiration,
  }
  return response, status.Updated()
}

func (c *credentialService) GetCredentials(ctx context.Context, userId types.Id) ([]dto.CredentialResponseDTO, status.Object) {
  ctx, span := c.tracer.Start(ctx, "CredentialService.GetCredentials")
  defer span.End()

  if err := c.checkPermission(ctx, userId, constant.AUTHN_PERMISSIONS[constant.AUTHN_GET_CREDENTIAL_ARB]); err != nil {
    spanUtil.RecordError(err, span)
    return nil, status.ErrUnAuthorized(err)
  }

  credentials, err := c.credRepo.FindByUserId(ctx, userId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, status.FromRepository(err, status.NullCode)
  }

  claims := types.Must(sharedJwt.GetUserClaimsFromCtx(ctx))
  response := sharedUtil.CastSliceP(credentials, func(from *entity.Credential) dto.CredentialResponseDTO {
    return mapper.ToCredentialResponseDTO(from, optional.Some(claims.CredentialId))
  })
  return response, status.Success()
}

func (c *credentialService) logoutOther(ctx context.Context, credentialIds []types.Id) status.Object {
  ctx, span := c.tracer.Start(ctx, "CredentialService.LogoutOther")
  defer span.End()

  // Delete arbitrary user credentials
  err := c.credRepo.Delete(ctx, credentialIds...)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }
  return status.Deleted()
}

func (c *credentialService) Logout(ctx context.Context, logoutDTO *dto.LogoutDTO) status.Object {
  ctx, span := c.tracer.Start(ctx, "CredentialService.Logout")
  defer span.End()

  // Get claims
  userClaims, err := sharedJwt.GetUserClaimsFromCtx(ctx)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrUnAuthenticated(err)
  }

  // Check if dto user id is the same with claims user id
  if !logoutDTO.UserId.EqWithString(userClaims.UserId) {
    // Check permission needed
    if !authUtil.ContainsPermission(userClaims.Roles, constant.AUTHN_PERMISSIONS[constant.AUTHN_LOGOUT_USER_ARB]) {
      spanUtil.RecordError(sharedErr.ErrUnauthorizedPermission, span)
      return status.ErrUnAuthorized(sharedErr.ErrUnauthorizedPermission)
    }
    // Logout other user
    return c.logoutOther(ctx, logoutDTO.CredentialIds)
  }

  err = c.credRepo.DeleteByUserId(ctx, logoutDTO.UserId, logoutDTO.CredentialIds...)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }

  return status.Deleted()
}

func (c *credentialService) LogoutAll(ctx context.Context, userId types.Id) status.Object {
  ctx, span := c.tracer.Start(ctx, "CredentialService.LogoutAll")
  defer span.End()

  // Permission check if needed
  err := c.checkPermission(ctx, userId, constant.AUTHN_PERMISSIONS[constant.AUTHN_LOGOUT_USER_ARB])
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrUnAuthorized(err)
  }

  // Delete all user credentials
  err = c.credRepo.DeleteByUserId(ctx, userId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }

  return status.Deleted()
}

func (c *CredentialConfig) generatePairTokens(username string, userId types.Id, roles []sharedJwt.Role) (entity.PairTokens, error) {
  type Return = entity.PairTokens

  refreshToken, err := c.generateRefreshToken()
  if err != nil {
    return Return{}, err
  }

  accessToken, err := c.generateAccessToken(username, userId, refreshToken.Id, roles)
  if err != nil {
    return Return{}, err
  }

  return Return{
    Access:  accessToken,
    Refresh: refreshToken,
  }, nil
}

func (c *CredentialConfig) generateRefreshToken() (entity.JWTToken, error) {
  rtid, err := types.NewId()
  if err != nil {
    return types.Null[entity.JWTToken](), err
  }

  refreshToken := sharedJwt.GenerateRefreshToken()
  return entity.JWTToken{
    Id:    rtid,
    Token: refreshToken,
  }, nil
}

func (c *CredentialConfig) generateAccessToken(username string, userId, refreshId types.Id, roles []sharedJwt.Role) (entity.JWTToken, error) {
  ct := time.Now()
  expAt := jwt.NewNumericDate(ct.Add(c.AccessTokenExpiration))

  id := types.MustCreateId()

  accessClaims := sharedJwt.UserClaims{
    RegisteredClaims: jwt.RegisteredClaims{
      Issuer:    constant.CLAIMS_ISSUER,
      ExpiresAt: expAt,
      NotBefore: jwt.NewNumericDate(ct),
      IssuedAt:  jwt.NewNumericDate(ct),
      ID:        id.String(),
    },
    CredentialId: refreshId.String(),
    UserId:       userId.String(),
    Username:     username,
    Roles:        roles,
  }
  accessToken := jwt.NewWithClaims(c.SigningMethod, accessClaims)
  accessSignedString, err := accessToken.SignedString(c.PrivateKey)
  if err != nil {
    return types.Null[entity.JWTToken](), err
  }

  return entity.JWTToken{
    Id:    id,
    Token: accessSignedString,
  }, nil
}
