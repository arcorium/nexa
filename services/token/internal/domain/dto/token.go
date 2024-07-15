package dto

import (
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/token/internal/domain/entity"
  "time"
)

type TokenCreateDTO struct {
  UserId types.Id
  Type   entity.TokenType
  Length uint32
  Usage  entity.TokenUsage
}

func (c *TokenCreateDTO) ToDomain(expiryTime time.Duration) entity.Token {
  return entity.NewToken(c.UserId, expiryTime, c.Usage, c.Type, uint16(c.Length))
}

type TokenVerifyDTO struct {
  Token         string `validate:"required"`
  ExpectedUsage entity.TokenUsage
}

type TokenAuthVerifyDTO struct {
  TokenVerifyDTO
  ExpectedUserId types.Id
}

type TokenResponseDTO struct {
  Token     string
  UserId    types.Id
  Usage     entity.TokenUsage
  ExpiredAt time.Time
}
