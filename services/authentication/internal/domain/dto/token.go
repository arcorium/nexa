package dto

import (
  entity2 "nexa/services/authentication/internal/domain/entity"
  "nexa/shared/types"
  "nexa/shared/util"
  "time"
)

type TokenRequestDTO struct {
  UsageId string `validate:"required,uuid4"`
}

func (r *TokenRequestDTO) ToEntity(userId types.Id) entity2.Token {
  return entity2.Token{
    Token:  util.RandomString(64),
    UserId: userId,
    Usage: entity2.TokenUsage{
      Id: types.IdFromString(r.UsageId),
    },
    ExpiredAt: time.Now(),
  }
}

type TokenRequestResponseDTO struct {
  Token string
}

type TokenVerifyDTO struct {
  Token   string `validate:"required,len=64"` //TODO: set length of string
  UsageId string `validate:"required,uuid4"`
}
