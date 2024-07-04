package dto

import (
  "nexa/services/user/constant"
  "nexa/shared/types"
  "time"
)

type TokenPurpose uint8

const (
  EmailVerificationToken TokenPurpose = iota
  ForgotPasswordToken
)

func NewEmailVerificationToken(userId types.Id) TokenGenerationDTO {
  return TokenGenerationDTO{
    UserId:  userId,
    Purpose: EmailVerificationToken,
    TTL:     constant.EMAIL_VERIFICAITON_TOKEN_TTL,
  }
}

func NewForgotPasswordToken(userId types.Id) TokenGenerationDTO {
  return TokenGenerationDTO{
    UserId:  userId,
    Purpose: ForgotPasswordToken,
    TTL:     constant.FORGOT_PASSWORD_TOKEN_TTL,
  }
}

type TokenGenerationDTO struct {
  UserId  types.Id
  Purpose TokenPurpose
  TTL     time.Duration
}

type TokenVerificationDTO struct {
  Token   string
  Purpose TokenPurpose
}

type TokenResponseDTO struct {
  Token     string
  Purpose   TokenPurpose
  ExpiredAt time.Time
}
