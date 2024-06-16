package dto

import (
  "nexa/services/authentication/shared/domain/entity"
  "nexa/shared/types"
  "nexa/shared/util"
  "time"
)

type TokenRequestDTO struct {
  UsageId string `validate:"required,uuid4"`
}

func (r *TokenRequestDTO) ToEntity(userId types.Id) entity.Token {
  return entity.Token{
    Token:  util.RandomString(64),
    UserId: userId,
    Usage: entity.TokenUsage{
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
