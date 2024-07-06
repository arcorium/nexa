package service

import (
  "context"
  "database/sql"
  "github.com/golang-jwt/jwt/v5"
  "go.opentelemetry.io/otel/trace"
  "nexa/services/authentication/constant"
  "nexa/services/authentication/internal/domain/dto"
  "nexa/services/authentication/internal/domain/entity"
  "nexa/services/authentication/internal/domain/external"
  "nexa/services/authentication/internal/domain/mapper"
  "nexa/services/authentication/internal/domain/repository"
  "nexa/services/authentication/internal/domain/service"
  "nexa/services/authentication/util"
  "nexa/services/authentication/util/errors"
  sharedErr "nexa/shared/errors"
  sharedJwt "nexa/shared/jwt"
  "nexa/shared/optional"
  "nexa/shared/status"
  "nexa/shared/types"
  sharedUtil "nexa/shared/util"
  authUtil "nexa/shared/util/auth"
  spanUtil "nexa/shared/util/span"
  "time"
)

func NewCredential(credential repository.ICredential, token repository.IToken, roleExt external.IRoleClient, userExt external.IUserClient, mailExt external.IMailClient, config CredentialServiceConfig) service.ICredential {
  return &credentialService{
    credRepo:  credential,
    tokenRepo: token,
    userExt:   userExt,
    roleExt:   roleExt,
    mailExt:   mailExt,
    config:    config,
    tracer:    util.GetTracer(),
  }
}

type CredentialServiceConfig struct {
  SigningMethod          jwt.SigningMethod
  AccessTokenExpiration  time.Duration
  RefreshTokenExpiration time.Duration
  SecretKey              string
}

type credentialService struct {
  credRepo  repository.ICredential
  tokenRepo repository.IToken

  userExt external.IUserClient
  roleExt external.IRoleClient
  mailExt external.IMailClient

  config CredentialServiceConfig
  tracer trace.Tracer
}

func (c *credentialService) checkPermission(ctx context.Context, targetId types.Id, permissions ...string) error {
  // Validate permission
  claims, _ := sharedJwt.GetClaimsFromCtx(ctx)
  if !targetId.EqWithString(claims.UserId) {
    // Need permission to update other users
    if !authUtil.ContainsPermissions(claims.Roles, permissions...) {
      return sharedErr.ErrUnauthorizedPermission
    }
  }
  return nil
}

func (c *credentialService) Login(ctx context.Context, loginDto *dto.LoginDTO) (dto.LoginResponseDTO, status.Object) {
  ctx, span := c.tracer.Start(ctx, "CredentialService.Login")
  defer span.End()

  // Get user by the email and validate the password
  user, err := c.userExt.Validate(ctx, loginDto.Email, loginDto.Password)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.LoginResponseDTO{}, status.ErrExternal(err)
  }

  // Get user roles and permissions
  roles, err := c.roleExt.GetUserRoles(ctx, user.UserId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.LoginResponseDTO{}, status.ErrExternal(err)
  }

  jwtRoles := sharedUtil.CastSliceP(roles, func(from *dto.RoleResponseDTO) sharedJwt.Role {
    return from.ToJWT()
  })

  // Generate token pairs
  pairTokens, err := c.config.generatePairTokens(user.Username, user.UserId, jwtRoles)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.LoginResponseDTO{}, status.ErrInternal(err)
  }

  // Save credentials
  credential := loginDto.ToDomain(user.UserId, pairTokens.Access.Id, &pairTokens.Refresh, c.config.RefreshTokenExpiration)
  err = c.credRepo.Create(ctx, &credential)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.LoginResponseDTO{}, status.FromRepository(err, status.NullCode)
  }

  return dto.LoginResponseDTO{
    TokenType: constant.TOKEN_TYPE,
    Token:     pairTokens.Access.Token,
  }, status.Success()
}

func (c *credentialService) Register(ctx context.Context, registerDTO *dto.RegisterDTO) status.Object {
  ctx, span := c.tracer.Start(ctx, "CredentialService.Register")
  defer span.End()

  // Create user
  userId, err := c.userExt.Create(ctx, registerDTO)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrExternal(err)
  }

  // Create verification token
  token := entity.NewToken(userId, entity.TokenUsageVerification, constant.TOKEN_VERIFICATION_EXPIRY_TIME)
  err = c.tokenRepo.Create(ctx, &token)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }

  // Send email verification
  mailDTO := dto.SendVerificationEmailDTO{
    Recipient: registerDTO.Email,
    Token:     token.Token,
  }
  err = c.mailExt.Send(ctx, &mailDTO)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrExternal(err)
  }

  return status.Created()
}

func (c *credentialService) RefreshToken(ctx context.Context, refreshDto *dto.RefreshTokenDTO) (dto.RefreshTokenResponseDTO, status.Object) {
  ctx, span := c.tracer.Start(ctx, "CredentialService.RefreshToken")
  defer span.End()

  if err := sharedUtil.ValidateStructCtx(ctx, refreshDto); err != nil {
    spanUtil.RecordError(err, span)
    return dto.RefreshTokenResponseDTO{}, status.ErrBadRequest(err)
  }

  // Check scheme
  if refreshDto.TokenType != constant.TOKEN_TYPE {
    return dto.RefreshTokenResponseDTO{}, status.ErrBadRequest(errors.ErrDifferentScheme)
  }

  // Parse token
  token, err := jwt.ParseWithClaims(refreshDto.AccessToken, &sharedJwt.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
    return []byte(c.config.SecretKey), nil
  })
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.RefreshTokenResponseDTO{}, status.ErrBadRequest(err)
  }

  claims := token.Claims.(sharedJwt.UserClaims)
  rtid := types.Must(types.IdFromString(claims.RefreshTokenId))
  cred, err := c.credRepo.Find(ctx, rtid)
  if err != nil {
    spanUtil.RecordError(err, span)
    if err == sql.ErrNoRows {
      return dto.RefreshTokenResponseDTO{}, status.ErrBadRequest(errors.ErrRefreshTokenNotFound)
    }
    return dto.RefreshTokenResponseDTO{}, status.FromRepository(err, optional.Some(status.BAD_REQUEST_ERROR))
  }

  // Check relation
  if !cred.UserId.EqWithString(claims.UserId) || !cred.AccessTokenId.EqWithString(claims.ID) {
    spanUtil.RecordError(errors.ErrBadRelation, span)
    return dto.RefreshTokenResponseDTO{}, status.ErrBadRequest(errors.ErrBadRelation)
  }

  // Get user roles and permissions
  roles, err := c.roleExt.GetUserRoles(ctx, cred.UserId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.RefreshTokenResponseDTO{}, status.ErrExternal(err)
  }

  jwtRoles := sharedUtil.CastSliceP(roles, func(from *dto.RoleResponseDTO) sharedJwt.Role {
    return from.ToJWT()
  })

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
  }
  return response, status.Updated()
}

func (c *credentialService) GetCredentials(ctx context.Context, userId types.Id) ([]dto.CredentialResponseDTO, status.Object) {
  ctx, span := c.tracer.Start(ctx, "CredentialService.GetCredentials")
  defer span.End()

  if err := c.checkPermission(ctx, userId, constant.AUTHN_PERMISSIONS[constant.AUTHN_GET_OTHER_CREDENTIALS]); err != nil {
    spanUtil.RecordError(err, span)
    return nil, status.ErrUnAuthorized(err)
  }

  credentials, err := c.credRepo.FindByUserId(ctx, userId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, status.FromRepository(err, status.NullCode)
  }

  response := sharedUtil.CastSliceP(credentials, mapper.ToCredentialResponseDTO)
  return response, status.Success()
}

func (c *credentialService) Logout(ctx context.Context, logoutDTO *dto.LogoutDTO) status.Object {
  ctx, span := c.tracer.Start(ctx, "CredentialService.Logout")
  defer span.End()

  // Get claims
  userClaims, err := sharedJwt.GetClaimsFromCtx(ctx)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrUnAuthenticated(err)
  }

  isPrivileged := authUtil.ContainsPermissions(userClaims.Roles, constant.AUTHN_PERMISSIONS[constant.AUTHN_LOGOUT_OTHER])

  // Check if userid empty
  if !logoutDTO.UserId.HasValue() {
    if !isPrivileged {
      spanUtil.RecordError(err, span)
      return status.ErrUnAuthorized(sharedErr.ErrUnauthorizedPermission)
    }

    // Delete arbitrary user credentials
    err = c.credRepo.Delete(ctx, logoutDTO.CredentialIds...)
    if err != nil {
      spanUtil.RecordError(err, span)
      return status.FromRepository(err, status.NullCode)
    }
    return status.Deleted()
  }

  // Check if user id is different from claims
  if !logoutDTO.UserId.Value().EqWithString(userClaims.UserId) && !isPrivileged {
    spanUtil.RecordError(err, span)
    return status.ErrUnAuthorized(sharedErr.ErrUnauthorizedPermission)
  }

  err = c.credRepo.DeleteByUserId(ctx, logoutDTO.UserId.RawValue(), logoutDTO.CredentialIds...)
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
  err := c.checkPermission(ctx, userId, constant.AUTHN_PERMISSIONS[constant.AUTHN_LOGOUT_OTHER])
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrUnAuthorized(sharedErr.ErrUnauthorizedPermission)
  }

  err = c.credRepo.DeleteByUserId(ctx, userId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }

  return status.Deleted()
}

func (c *CredentialServiceConfig) generatePairTokens(username string, userId types.Id, roles []sharedJwt.Role) (entity.PairTokens, error) {
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

func (c *CredentialServiceConfig) generateRefreshToken() (entity.JWTToken, error) {
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

func (c *CredentialServiceConfig) generateAccessToken(username string, userId, refreshId types.Id, roles []sharedJwt.Role) (entity.JWTToken, error) {
  ct := time.Now()
  expAt := jwt.NewNumericDate(ct.Add(c.AccessTokenExpiration))

  id := types.MustCreateId()

  accessClaims := sharedJwt.UserClaims{
    RegisteredClaims: jwt.RegisteredClaims{
      Issuer:    constant.CLAIMS_ISSUER,
      ExpiresAt: expAt,
      NotBefore: expAt,
      IssuedAt:  jwt.NewNumericDate(ct),
      ID:        id.String(),
    },
    RefreshTokenId: refreshId.String(),
    UserId:         userId.String(),
    Username:       username,
    Roles:          roles,
  }
  accessToken := jwt.NewWithClaims(c.SigningMethod, accessClaims)
  accessSignedString, err := accessToken.SignedString([]byte(c.SecretKey))
  if err != nil {
    return types.Null[entity.JWTToken](), err
  }

  return entity.JWTToken{
    Id:    id,
    Token: accessSignedString,
  }, nil
}
