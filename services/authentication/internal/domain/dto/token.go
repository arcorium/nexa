package dto

import (
  "github.com/arcorium/nexa/shared/types"
  "time"
)

type TokenUsage uint8

const (
  TokenUsageEmailVerification TokenUsage = iota
  TokenUsageResetPassword
  TokenUsageLogin
  TokenUsageGeneral
)

type TokenGenerationDTO struct {
  UserId types.Id
  Usage  TokenUsage
}

type TokenVerificationDTO struct {
  Token   string
  Purpose TokenUsage
}

type TokenResponseDTO struct {
  Token     string
  UserId    types.Id
  Usage     TokenUsage
  ExpiredAt time.Time
}
