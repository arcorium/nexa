package dto

import (
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/authentication/internal/domain/entity"
  "time"
)

type LoginDTO struct {
  Email      types.Email
  Password   types.Password
  DeviceName string `validate:"required"`
}

func (d *LoginDTO) ToDomain(userId, accessTokenId types.Id, refreshToken *entity.JWTToken, expiryTime time.Duration) entity.Credential {
  return entity.Credential{
    Id:            refreshToken.Id,
    UserId:        userId,
    AccessTokenId: accessTokenId,
    Device:        entity.Device{Name: d.DeviceName},
    RefreshToken:  refreshToken.Token,
    ExpiresAt:     time.Now().UTC().Add(expiryTime),
  }
}

type RegisterDTO struct {
  Username  string `validate:"required,gte=6"`
  Email     types.Email
  Password  types.Password
  FirstName string `validate:"required"`
  LastName  types.NullableString
  Bio       types.NullableString
}

func (d *RegisterDTO) ToDomain() (entity.User, entity.Profile, error) {
  user, err := entity.NewUser(d.Username, d.Email, d.Password)
  if err != nil {
    return entity.User{}, entity.Profile{}, err
  }

  profile, err := entity.NewProfile(user.Id, d.FirstName)
  if err != nil {
    return entity.User{}, entity.Profile{}, err
  }

  types.SetOnNonNull(&profile.LastName, d.LastName)
  types.SetOnNonNull(&profile.Bio, d.Bio)

  return user, profile, nil
}

type LoginResponseDTO struct {
  TokenType  string
  Token      string
  ExpiryTime time.Duration
}

type RefreshTokenDTO struct {
  TokenType   string `validate:"required"`
  AccessToken string `validate:"required,jwt"`
}

func (r *RefreshTokenDTO) ToDomain(credId, accessTokenId types.Id) entity.Credential {
  return entity.Credential{
    Id:            credId,
    AccessTokenId: accessTokenId,
  }
}

type RefreshTokenResponseDTO struct {
  TokenType   string
  AccessToken string
  ExpiryTime  time.Duration
}

type LogoutDTO struct {
  UserId        types.Id
  CredentialIds []types.Id
}
