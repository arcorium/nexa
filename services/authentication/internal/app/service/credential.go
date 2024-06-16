package service

import (
  "context"
  "github.com/golang-jwt/jwt/v5"
  "nexa/services/authentication/config"
  "nexa/services/authentication/internal/domain/dto"
  "nexa/services/authentication/internal/domain/external"
  "nexa/services/authentication/internal/domain/mapper"
  "nexa/services/authentication/internal/domain/repository"
  "nexa/services/authentication/internal/domain/service"
  "nexa/services/authentication/shared/common"
  "nexa/services/authentication/shared/domain/entity"
  "nexa/services/authentication/shared/errors"
  appUtil "nexa/services/authentication/util"
  "nexa/shared/optional"
  "nexa/shared/status"
  "nexa/shared/types"
  "nexa/shared/util"
  "nexa/shared/variadic"
  "time"
)

func NewCredential(credential repository.ICredential, userExt external.IUserClient, serverConfig *config.ServerConfig) service.ICredential {
  return &credentialService{
    repo:    credential,
    userExt: userExt,
    cfg:     serverConfig,
  }
}

type credentialService struct {
  repo    repository.ICredential
  userExt external.IUserClient

  cfg *config.ServerConfig
}

func (c *credentialService) Login(ctx context.Context, dto *dto.LoginDTO) (string, status.Object) {
  // Get user by the email and validate the password
  user, err := c.userExt.ValidateUser(ctx, types.EmailFromString(dto.Email), dto.Password)
  if err != nil {
    return "", status.ErrExternal(err)
  }
  // Create pair tokens
  pairTokens, stats := common.GeneratePairTokens(c.cfg.SigningMethod(), c.cfg.JWTTokenExpiration, c.cfg.SecretKey(), user.Id, user.Username)
  if stats.IsError() {
    return "", stats
  }

  // Create credentials
  credential := dto.ToEntity(user.Id, pairTokens.AccessToken.Id, pairTokens.RefreshToken.Id, "TEST", pairTokens.RefreshToken.String)
  err = c.repo.Create(ctx, &credential)
  if err != nil {
    return "", status.FromRepository(err, status.NullCode)
  }
  return pairTokens.AccessToken.String, status.Success()
}

func (c *credentialService) Register(ctx context.Context, dto *dto.RegisterDTO) status.Object {
  err := c.userExt.RegisterUser(ctx, dto)
  if err != nil {
    return status.ErrExternal(err)
  }
  return status.Created()
}

func (c *credentialService) RefreshToken(ctx context.Context, dto *dto.RefreshTokenDTO) (string, status.Object) {
  // Get refresh token from repo
  var accessClaims common.AccessTokenClaims
  _, err := jwt.ParseWithClaims(dto.AccessToken, &accessClaims, func(token *jwt.Token) (interface{}, error) {
    return c.cfg.SecretKey(), nil
  })
  if err != nil {
    return "", status.ErrBadRequest(err)
  }

  credential, err := c.repo.Find(ctx, types.IdFromString(accessClaims.ID))
  if err != nil {
    return "", status.FromRepository(err, optional.Some(status.BAD_REQUEST_ERROR))
  }
  if accessClaims.UserId != credential.UserId {
    return "", status.ErrBadRequest(nil)
  }
  // Extract refresh token
  var refreshClaims common.RefreshTokenClaims
  _, err = jwt.ParseWithClaims(credential.RefreshToken, &refreshClaims, c.cfg.KeyFunc())
  if err != nil {
    return "", status.ErrInternal(err)
  }
  // Check expiration time
  if refreshClaims.ExpiresAt.After(time.Now()) {
    // Remove from repo
    err := c.repo.Delete(ctx, types.IdFromString(accessClaims.ID))
    if err != nil {
      return "", status.ErrInternal(err)
    }
    return "", status.ErrUnAuthorized(errors.ErrRefreshTokenExpired)
  }

  // Create new access token
  accessId, newAccessToken, stats := common.GenerateAccessToken(c.cfg.SigningMethod(), c.cfg.JWTTokenExpiration, c.cfg.SecretKey(), accessClaims.UserId, accessClaims.Username)
  if stats.IsError() {
    return "", stats
  }
  // Update credential
  updateCred := entity.Credential{
    Id:            types.IdFromString(refreshClaims.ID),
    AccessTokenId: accessId,
  }

  err = c.repo.Patch(ctx, &updateCred)
  if err != nil {
    return "", status.FromRepository(err, optional.Some(status.INTERNAL_SERVER_ERROR))
  }
  return newAccessToken, status.Success()
}

func (c *credentialService) GetCurrentCredentials(ctx context.Context) ([]dto.CredentialResponseDTO, status.Object) {
  claims := appUtil.GetUserClaims(ctx)
  credentials, err := c.repo.FindByUserId(ctx, claims.UserId)
  if err != nil {
    return nil, status.FromRepository(err, status.NullCode)
  }
  responses := util.CastSlice(credentials, mapper.ToCredentialResponseDTO)
  return responses, status.Success()
}

func (c *credentialService) Logout(ctx context.Context, credIds ...types.Id) status.Object {
  var err error
  ids := variadic.New(credIds...)
  if !ids.HasValue() {
    // Logout current
    claims := appUtil.GetUserClaims(ctx)

    err = c.repo.Delete(ctx, types.IdFromString(claims.ID))
  } else {
    // Logout the ids
    err = c.repo.Delete(ctx, credIds...)
  }
  if err != nil {
    return status.FromRepository(err, status.NullCode)
  }
  return status.Success()
}

func (c *credentialService) LogoutAll(ctx context.Context) status.Object {
  // Logout current
  claims := appUtil.GetUserClaims(ctx)
  err := c.repo.DeleteByUserId(ctx, claims.UserId)
  if err != nil {
    return status.FromRepository(err, status.NullCode)
  }
  return status.Success()
}
