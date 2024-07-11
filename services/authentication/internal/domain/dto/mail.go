package dto

import "github.com/arcorium/nexa/shared/types"

type SendVerificationEmailDTO struct {
  Recipient types.Email
  Token     string
}
