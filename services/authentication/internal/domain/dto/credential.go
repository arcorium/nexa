package dto

import (
  "nexa/services/authentication/shared/domain/entity"
  "nexa/services/authentication/shared/domain/valueobject"
  "nexa/shared/types"
  "nexa/shared/wrapper"
)

type LoginDTO struct {
  Email    string `validate:"required,email"`
  Password string `validate:"required"`
}

func (d *LoginDTO) ToEntity(userId, accessTokenId, refreshTokenId types.Id, deviceName string, refreshToken string) entity.Credential {
  return entity.Credential{
    Id:            refreshTokenId,
    UserId:        userId,
    AccessTokenId: accessTokenId,
    Device:        valueobject.Device{Name: deviceName},
    RefreshToken:  refreshToken,
  }
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
  AccessToken string `validate:"required"`
}

type CredentialResponseDTO struct {
  Id     string
  Device string
}
