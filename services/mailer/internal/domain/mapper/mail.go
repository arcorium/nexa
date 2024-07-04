package mapper

import (
  "nexa/services/mailer/internal/domain/dto"
  domain "nexa/services/mailer/internal/domain/entity"
  "nexa/shared/util"
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
