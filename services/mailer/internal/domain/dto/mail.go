package dto

import (
  "nexa/services/mailer/constant"
  domain "nexa/services/mailer/internal/domain/entity"
  "nexa/shared/types"
  sharedUtil "nexa/shared/util"
  "nexa/shared/wrapper"
  "time"
)

type FileAttachment struct {
  Filename string `validate:"required"`
  Data     []byte `validate:"required"`
}

type SendMailDTO struct {
  Subject     string        `validate:"required"`
  Recipients  []types.Email `validate:"required"`
  Sender      wrapper.Nullable[types.Email]
  BodyType    domain.MailBodyType ` validate:"required"`
  Body        string
  TagIds      []types.Id `validate:"required"`
  Attachments []FileAttachment
}

func (m *SendMailDTO) ToDomain() ([]domain.Mail, error) {
  tags := sharedUtil.CastSlice(m.TagIds, func(tagId types.Id) domain.Tag {
    return domain.Tag{
      Id: tagId,
    }
  })

  mails := make([]domain.Mail, len(m.Recipients))
  for _, email := range m.Recipients {
    mailId, err := types.NewId()
    if err != nil {
      return nil, err
    }

    mail := domain.Mail{
      Id:        mailId,
      Subject:   m.Subject,
      Recipient: email,
      Sender:    m.Sender.ValueOr(constant.SERVICE_MAIL_SENDER),
      BodyType:  m.BodyType,
      Body:      m.Body,
      Status:    domain.StatusPending,
      SentAt:    time.Now(),
      Tags:      tags,
    }

    mails = append(mails, mail)
  }

  return mails, nil
}

type UpdateMailDTO struct {
  Id            types.Id
  AddedTagIds   []types.Id
  RemovedTagIds []types.Id
}

type MailResponseDTO struct {
  Id          types.Id
  Subject     string
  Recipient   types.Email
  Sender      types.Email
  Status      domain.Status
  SentAt      time.Time
  DeliveredAt time.Time
  Tags        []TagResponseDTO
}
