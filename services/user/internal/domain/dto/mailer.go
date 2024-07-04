package dto

import "nexa/shared/types"

type SendEmailVerificationDTO struct {
  Recipient types.Email
  Token     string
}

type SendForgotPasswordDTO struct {
  Recipient types.Email
  Token     string
}
