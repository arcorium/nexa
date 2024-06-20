package mapper

import (
  "nexa/services/mailer/constant"
  "nexa/services/mailer/internal/domain/dto"
  domain "nexa/services/mailer/internal/domain/entity"
  "nexa/shared/types"
  "nexa/shared/util"
  "nexa/shared/wrapper"
)

func MapSendMailDTO(mailDTO *dto.SendMailDTO) []domain.Mail {
  recipientEmails, _ := util.CastSliceErrs(mailDTO.Recipients, types.EmailFromString)

  mails := make([]domain.Mail, len(mailDTO.Recipients))

  for _, email := range recipientEmails {
    mail := domain.Mail{
      Id:        types.NewId2(),
      Subject:   mailDTO.Subject,
      Recipient: email,
      BodyType:  domain.MailBodyType(mailDTO.BodyType),
      Body:      mailDTO.Body,
      Status:    domain.StatusPending,
      Tags: util.CastSlice(mailDTO.TagIds, func(tagId string) domain.Tag {
        return domain.Tag{
          Id: wrapper.DropError(types.IdFromString(tagId)),
        }
      }),
    }

    if !mailDTO.Sender.HasValue() {
      // Set service default as sender
      mail.Sender = constant.SERVICE_MAIL_SENDER
    } else {
      senderEmail, _ := types.EmailFromString(mailDTO.Sender.RawValue())
      mail.Sender = senderEmail
    }

    mails = append(mails, mail)
  }
  return mails
}

func ToMailResponseDTO(mail *domain.Mail) dto.MailResponseDTO {
  return dto.MailResponseDTO{
    Id:        mail.Id.Underlying().String(),
    Subject:   mail.Subject,
    Recipient: mail.Recipient.Underlying(),
    Sender:    mail.Sender.Underlying(),
    Status:    mail.Status.String(),
    Tags:      util.CastSliceP(mail.Tags, ToResponseDTO),
  }
}
