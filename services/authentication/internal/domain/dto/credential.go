package dto

import (
  domain "nexa/services/authentication/internal/domain/entity"
  "nexa/shared/types"
  "nexa/shared/wrapper"
  "time"
)

type LoginDTO struct {
  Email      string `json:"email" validate:"required,email"`
  Password   string `json:"password" validate:"required"`
  DeviceName string `json:"device_name" validate:"required"`
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
  TokenType string `json:"token_type"`
  Token     string `json:"token"`
}

type RegisterDTO struct {
  Username  string `validate:"required,gte=6"`
  Email     string `validate:"required,email"`
  Password  string `validate:"required,gte=6"`
  FirstName string `validate:"required"`
  LastName  wrapper.NullableString
  Bio       wrapper.NullableString
}

type RefreshTokenDTO struct {
  TokenType   string `json:"token_type" validate:"required"`
  AccessToken string `json:"access_token" validate:"required"`
}

func (r *RefreshTokenDTO) ToDomain(refreshTokenId, accessTokenId types.Id) domain.Credential {
  return domain.Credential{
    Id:            refreshTokenId,
    AccessTokenId: accessTokenId,
  }
}

type RefreshTokenResponseDTO struct {
  TokenType   string `json:"token_type"`
  AccessToken string `json:"token"`
}

type CredentialResponseDTO struct {
  Id     string `json:"id"`
  Device string `json:"device"`
}

type LogoutDTO struct {
  UserId        string   `validate:"required,uuid4"`
  CredentialIds []string `validate:"required,dive,uuid4"`
}
