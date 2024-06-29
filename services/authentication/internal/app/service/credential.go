package service

import (
  "context"
  "database/sql"
  "github.com/golang-jwt/jwt/v5"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/genproto/googleapis/rpc/errdetails"
  "nexa/services/authentication/config"
  "nexa/services/authentication/constant"
  "nexa/services/authentication/internal/domain/dto"
  "nexa/services/authentication/internal/domain/entity"
  "nexa/services/authentication/internal/domain/external"
  "nexa/services/authentication/internal/domain/mapper"
  "nexa/services/authentication/internal/domain/repository"
  "nexa/services/authentication/internal/domain/service"
  "nexa/services/authentication/util"
  "nexa/services/authentication/util/errors"
  "nexa/shared/auth"
  sharedErr "nexa/shared/errors"
  sharedJwt "nexa/shared/jwt"
  "nexa/shared/optional"
  spanUtil "nexa/shared/span"
  "nexa/shared/status"
  "nexa/shared/types"
  sharedUtil "nexa/shared/util"
  "nexa/shared/wrapper"
  "time"
)

func NewCredential(credential repository.ICredential, userExt external.IUserClient, serverConfig *config.Server) service.ICredential {
  return &credentialService{
    credRepo: credential,
    userExt:  userExt,
    cfg:      serverConfig,
    tracer:   util.GetTracer(),
  }
}

type CredentialServiceConfig struct {
  SigningMethod          jwt.SigningMethod
  AccessTokenExpiration  time.Duration
  RefreshTokenExpiration time.Duration
  SecretKey              string
}

type credentialService struct {
  credRepo repository.ICredential

  userExt external.IUserClient
  roleExt external.IRoleClient

  config CredentialServiceConfig
  tracer trace.Tracer
  cfg    *config.Server
}

func (c *credentialService) Login(ctx context.Context, loginDto *dto.LoginDTO) (dto.LoginResponseDTO, status.Object) {
  type RetType = dto.LoginResponseDTO

  ctx, span := c.tracer.Start(ctx, "CredentialService.Login")
  defer span.End()

  if err := sharedUtil.ValidateStructCtx(ctx, loginDto); err != nil {
    spanUtil.RecordError(err, span)
    return RetType{}, status.ErrBadRequest(err)
  }

  // Get user by the email and validate the password
  user, err := c.userExt.Validate(ctx, wrapper.Must(types.EmailFromString(loginDto.Email)), loginDto.Password)
  if err != nil {
    spanUtil.RecordError(err, span)
    return RetType{}, status.ErrExternal(err)
  }

  // Get user roles and permissions
  roles, err := c.roleExt.GetUserRoles(ctx, user.UserId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return RetType{}, status.ErrExternal(err)
  }

  jwtRoles := sharedUtil.CastSliceP(roles, func(from *dto.RoleResponseDTO) sharedJwt.Role {
    return from.ToJWT()
  })

  pairTokens, err := c.config.generatePairTokens(user.Username, user.UserId, jwtRoles)
  if err != nil {
    spanUtil.RecordError(err, span)
    return RetType{}, status.ErrInternal(err)
  }

  // Save credentials
  credential := loginDto.ToDomain(user.UserId, pairTokens.Access.Id, &pairTokens.Refresh, c.config.RefreshTokenExpiration)
  err = c.credRepo.Create(ctx, &credential)
  if err != nil {
    spanUtil.RecordError(err, span)
    return RetType{}, status.FromRepository(err, status.NullCode)
  }

  return RetType{
    TokenType: constant.TOKEN_TYPE,
    Token:     pairTokens.Access.Token,
  }, status.Success()
}

func (c *credentialService) Register(ctx context.Context, dto *dto.RegisterDTO) status.Object {
  ctx, span := c.tracer.Start(ctx, "CredentialService.Register")
  defer span.End()

  if err := sharedUtil.ValidateStructCtx(ctx, dto); err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrBadRequest(err)
  }

  err := c.userExt.Create(ctx, dto)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrExternal(err)
  }
  return status.Created()
}

func (c *credentialService) RefreshToken(ctx context.Context, refreshDto *dto.RefreshTokenDTO) (dto.RefreshTokenResponseDTO, status.Object) {
  type RetType = dto.RefreshTokenResponseDTO

  ctx, span := c.tracer.Start(ctx, "CredentialService.RefreshToken")
  defer span.End()

  if err := sharedUtil.ValidateStructCtx(ctx, refreshDto); err != nil {
    spanUtil.RecordError(err, span)
    return RetType{}, status.ErrBadRequest(err)
  }

  // Check scheme
  if refreshDto.TokenType != constant.TOKEN_TYPE {
    return RetType{}, status.ErrBadRequest(errors.ErrDifferentScheme)
  }

  // Parse token
  token, err := jwt.ParseWithClaims(refreshDto.AccessToken, &sharedJwt.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
    return []byte(c.config.SecretKey), nil
  })
  if err != nil {
    spanUtil.RecordError(err, span)
    return RetType{}, status.ErrBadRequest(err)
  }

  claims := token.Claims.(sharedJwt.UserClaims)
  rtid := wrapper.Must(types.IdFromString(claims.RefreshTokenId))
  cred, err := c.credRepo.Find(ctx, rtid)
  if err != nil {
    spanUtil.RecordError(err, span)
    if err == sql.ErrNoRows {
      return RetType{}, status.ErrBadRequest(errors.ErrRefreshTokenNotFound)
    }
    return RetType{}, status.FromRepository(err, optional.Some(status.BAD_REQUEST_ERROR))
  }

  // Check relation
  if !cred.UserId.EqWithString(claims.UserId) || !cred.AccessTokenId.EqWithString(claims.ID) {
    spanUtil.RecordError(errors.ErrBadRelation, span)
    return RetType{}, status.ErrBadRequest(errors.ErrBadRelation)
  }

  // Get user roles and permissions
  roles, err := c.roleExt.GetUserRoles(ctx, cred.UserId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return RetType{}, status.ErrExternal(err)
  }

  jwtRoles := sharedUtil.CastSliceP(roles, func(from *dto.RoleResponseDTO) sharedJwt.Role {
    return from.ToJWT()
  })

  // Create new access token
  accessToken, err := c.config.generateAccessToken(claims.Username, cred.UserId, cred.Id, jwtRoles)
  if err != nil {
    spanUtil.RecordError(err, span)
    return RetType{}, status.ErrInternal(err)
  }

  // Patch credentials with new access token
  updateCred := refreshDto.ToDomain(cred.Id, accessToken.Id)
  err = c.credRepo.Patch(ctx, &updateCred)
  if err != nil {
    spanUtil.RecordError(err, span)
    return RetType{}, status.FromRepository(err, optional.Some(status.INTERNAL_SERVER_ERROR))
  }

  response := RetType{
    TokenType:   constant.TOKEN_TYPE,
    AccessToken: accessToken.Token,
  }
  return response, status.Updated()
}

func (c *credentialService) GetCredentials(ctx context.Context, userId string) ([]dto.CredentialResponseDTO, status.Object) {
  ctx, span := c.tracer.Start(ctx, "CredentialService.GetCredentials")
  defer span.End()

  id, err := types.IdFromString(userId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, status.ErrBadRequest(sharedErr.GrpcFieldErrors(&errdetails.BadRequest_FieldViolation{
      Field:       "user_id",
      Description: err.Error(),
    }))
  }

  userClaims, err := sharedJwt.GetClaimsFromCtx(ctx)
  // Unauthenticated
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, status.ErrUnAuthenticated(err)
  }

  // Get other user credentials
  if !id.EqWithString(userClaims.UserId) &&
      !auth.ContainsPermissions(userClaims.Roles, constant.CRED_READ_OTHERS) {
    spanUtil.RecordError(sharedErr.ErrUnauthorizedPermission, span)
    return nil, status.ErrUnAuthenticated(sharedErr.ErrUnauthorizedPermission)
  }

  credentials, err := c.credRepo.FindByUserId(ctx, id)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, status.FromRepository(err, status.NullCode)
  }

  response := sharedUtil.CastSliceP(credentials, func(from *entity.Credential) dto.CredentialResponseDTO {
    return mapper.ToCredentialResponseDTO(from)
  })

  return response, status.Success()
}

func (c *credentialService) Logout(ctx context.Context, logoutDTO *dto.LogoutDTO) status.Object {
  ctx, span := c.tracer.Start(ctx, "CredentialService.Logout")
  defer span.End()

  // Input validation
  userId, err := types.IdFromString(logoutDTO.UserId)
  if err != nil {
    err := sharedErr.GrpcFieldErrors(&errdetails.BadRequest_FieldViolation{
      Field:       "user_id",
      Description: err.Error(),
    })
    spanUtil.RecordError(err, span)
    return status.ErrBadRequest(err)
  }

  credIds, ierr := sharedUtil.CastSliceErrs(logoutDTO.CredentialIds, types.IdFromString)
  if !ierr.IsNil() {
    err := sharedErr.GrpcFieldIndexedErrors("cred_ids", ierr)
    spanUtil.RecordError(err, span)
    return status.ErrBadRequest(err)
  }

  // Get claims
  userClaims, err := sharedJwt.GetClaimsFromCtx(ctx)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrUnAuthenticated(err)
  }

  isPrivileged := auth.ContainsPermissions(userClaims.Roles, constant.CRED_DELETE_OTHERS)

  if !userId.EqWithString(userClaims.UserId) {
    if !isPrivileged {
      spanUtil.RecordError(err, span)
      return status.ErrUnAuthorized(sharedErr.ErrUnauthorizedPermission)
    }

    // Delete arbitrary user credentials
    if logoutDTO.UserId == "" {
      err := c.credRepo.Delete(ctx)
      if err != nil {
        spanUtil.RecordError(err, span)
        return status.FromRepository(err, status.NullCode)
      }
      return status.Deleted()
    }
  }

  err = c.credRepo.DeleteByUserId(ctx, wrapper.Must(types.IdFromString(logoutDTO.UserId)), credIds...)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }

  return status.Deleted()
}

func (c *credentialService) LogoutAll(ctx context.Context, userId string) status.Object {
  ctx, span := c.tracer.Start(ctx, "CredentialService.LogoutAll")
  defer span.End()

  // Input validation
  id, err := types.IdFromString(userId)
  if err != nil {
    err := sharedErr.GrpcFieldErrors(&errdetails.BadRequest_FieldViolation{
      Field:       "user_id",
      Description: err.Error(),
    })
    spanUtil.RecordError(err, span)
    return status.ErrBadRequest(err)
  }

  userClaims, err := sharedJwt.GetClaimsFromCtx(ctx)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrUnAuthenticated(err)
  }

  isPrivileged := auth.ContainsPermissions(userClaims.Roles, constant.CRED_DELETE_OTHERS)
  if !id.EqWithString(userClaims.UserId) && !isPrivileged {
    spanUtil.RecordError(err, span)
    return status.ErrUnAuthorized(sharedErr.ErrUnauthorizedPermission)
  }

  err = c.credRepo.DeleteByUserId(ctx, id)
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

  id := types.NewId2()

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
