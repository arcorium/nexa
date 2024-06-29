package dto

import (
  domain "nexa/services/authentication/internal/domain/entity"
  sharedJwt "nexa/shared/jwt"
  "nexa/shared/types"
  sharedUtil "nexa/shared/util"
  "time"
)

type TokenCreateDTO struct {
  UserId string `validate:"required,uuid4"`
  Usage  uint8  `validate:"required,usage_enum"`
}

func (c *TokenCreateDTO) ToDomain(expiryTime time.Duration) (domain.Token, error) {
  type Null = domain.Token
  if err := sharedUtil.ValidateStruct(c); err != nil {
    return Null{}, err
  }

  userId, err := types.IdFromString(c.UserId)
  if err != nil {
    return Null{}, err
  }

  return domain.Token{
    Token:     sharedJwt.GenerateRefreshToken(),
    UserId:    userId,
    Usage:     domain.TokenUsage(c.Usage),
    ExpiredAt: time.Now().UTC().Add(expiryTime),
  }, nil
}

type TokenVerifyDTO struct {
  Token string `validate:"required"`
  Usage uint8  `validate:"required,usage_enum"`
}

type TokenResponseDTO struct {
  Token     string
  Usage     uint8
  ExpiredAt time.Time
}
