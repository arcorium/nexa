package external

import (
  "context"
  "nexa/services/user/internal/domain/dto"
)

type IMailerClient interface {
  SendEmailVerification(ctx context.Context, dto *dto.SendEmailVerificationDTO) error
  SendForgotPassword(ctx context.Context, passwordDTO *dto.SendForgotPasswordDTO) error
}
