package dto

import (
  domain "nexa/services/authentication/internal/domain/entity"
  "nexa/shared/types"
  "time"
)

type LoginDTO struct {
  Email      types.Email
  Password   types.Password
  DeviceName string `validate:"required"`
}

func (d *LoginDTO) ToDomain(userId, accessTokenId types.Id, refreshToken *domain.JWTToken, expiryTime time.Duration) domain.Credential {
  return domain.Credential{
    Id:            refreshToken.Id,
    UserId:        userId,
    AccessTokenId: accessTokenId,
    Device:        domain.Device{Name: d.DeviceName},
    RefreshToken:  refreshToken.Token,
    ExpiresAt:     time.Now().UTC().Add(expiryTime),
  }
}

type LoginResponseDTO struct {
  TokenType string
  Token     string
}

type RegisterDTO struct {
  Username  string `validate:"required,gte=6"`
  Email     types.Email
  Password  types.Password
  FirstName string `validate:"required"`
  LastName  types.NullableString
  Bio       types.NullableString
}

type RefreshTokenDTO struct {
  TokenType   string `validate:"required"`
  AccessToken string `validate:"required,jwt"`
}

func (r *RefreshTokenDTO) ToDomain(credId, accessTokenId types.Id) domain.Credential {
  return domain.Credential{
    Id:            credId,
    AccessTokenId: accessTokenId,
  }
}

type RefreshTokenResponseDTO struct {
  TokenType   string
  AccessToken string
}

type CredentialResponseDTO struct {
  Id     types.Id
  Device string
}

type LogoutDTO struct {
  UserId        types.Id
  CredentialIds []types.Id
}
