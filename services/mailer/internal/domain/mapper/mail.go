package mapper

import (
  "github.com/arcorium/nexa/shared/util"
  "nexa/services/mailer/internal/domain/dto"
  domain "nexa/services/mailer/internal/domain/entity"
)

func ToMailResponseDTO(mail *domain.Mail) dto.MailResponseDTO {
  return dto.MailResponseDTO{
    Id:          mail.Id,
    Subject:     mail.Subject,
    Recipient:   mail.Recipient,
    Sender:      mail.Sender,
    Status:      mail.Status,
    SentAt:      mail.SentAt,
    DeliveredAt: mail.DeliveredAt,
    Tags:        util.CastSliceP(mail.Tags, ToTagResponseDTO),
  }
}
