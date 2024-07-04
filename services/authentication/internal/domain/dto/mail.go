package dto

import "nexa/shared/types"

type SendVerificationEmailDTO struct {
  Recipient types.Email
  Token     string
}
