package external

import (
  "context"
  "nexa/services/mailer/internal/domain/dto"
  domain "nexa/services/mailer/internal/domain/entity"
)

type IMail interface {
  // Send send single email into single recipient
  Send(ctx context.Context, mail *domain.Mail, attachments []dto.FileAttachment) error
  // Close close the connection
  Close(context.Context) error
}
