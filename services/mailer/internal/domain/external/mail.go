package external

import (
  "context"
  "nexa/services/mailer/internal/domain/dto"
  domain "nexa/services/mailer/internal/domain/entity"
)

type IMail interface {
  Send(ctx context.Context, mail *domain.Mail, attachments []dto.FileAttached) error
}
