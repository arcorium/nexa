package dto

import (
  "nexa/services/authentication/internal/domain/entity"
  "nexa/shared/types"
  "time"
)

type TokenCreateDTO struct {
  UserId types.Id
  Usage  entity.TokenUsage
}

func (c *TokenCreateDTO) ToDomain(expiryTime time.Duration) entity.Token {
  return entity.NewToken(c.UserId, c.Usage, expiryTime)
}

type TokenVerifyDTO struct {
  Token string            `validate:"required"`
  Usage entity.TokenUsage `validate:"required,usage_enum"`
}

type TokenResponseDTO struct {
  Token     string
  Usage     entity.TokenUsage
  ExpiredAt time.Time
}
