package external

import (
  "context"
  "nexa/services/authentication/internal/domain/dto"
)

type IMailClient interface {
  Send(ctx context.Context, dto *dto.SendVerificationEmailDTO) error
}
