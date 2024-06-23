package dto

import (
  "nexa/services/user/constant"
  "time"
)

type TokenPurpose uint8

const (
  EmailVerificationToken TokenPurpose = iota
  ForgotPasswordToken
)

func NewEmailVerificationToken() TokenGenerationDTO {
  return TokenGenerationDTO{
    Purpose: EmailVerificationToken,
    TTL:     constant.EMAIL_VERIFICAITON_TOKEN_TTL,
  }
}

func NewForgotPasswordToken() TokenGenerationDTO {
  return TokenGenerationDTO{
    Purpose: ForgotPasswordToken,
    TTL:     constant.FORGOT_PASSWORD_TOKEN_TTL,
  }
}

type TokenGenerationDTO struct {
  Purpose TokenPurpose  `json:"purpose"`
  TTL     time.Duration `json:"ttl"`
}

type TokenResponseDTO struct {
  Token string `json:"token"`
}
