package dto

import "github.com/arcorium/nexa/shared/types"

type SendEmailVerificationDTO struct {
  Recipient types.Email
  Token     string
}

type SendForgotPasswordDTO struct {
  Recipient types.Email
  Token     string
}
