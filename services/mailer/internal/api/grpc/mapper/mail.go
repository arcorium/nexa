package mapper

import (
  mailerv1 "nexa/proto/gen/go/mailer/v1"
  "nexa/services/mailer/internal/domain/dto"
  "nexa/shared/util"
  "nexa/shared/wrapper"
)

func ToSendMailDTO(request *mailerv1.SendMailRequest) dto.SendMailDTO {
  return dto.SendMailDTO{
    Subject:    request.Subject,
    Recipients: request.Recipients,
    Sender:     wrapper.NewNullable(request.Sender),
    BodyType:   uint8(request.BodyType),
    Body:       request.Body,
    TagIds:     request.TagIds,
    Attachments: util.CastSlice(request.Attachments, func(from *mailerv1.FileAttachment) dto.FileAttachment {
      return dto.FileAttachment{
        Filename: from.Filename,
        Data:     from.Data,
      }
    }),
  }
}

func ToUpdateMailDTO(request *mailerv1.UpdateMailRequest) dto.UpdateMailDTO {
  return dto.UpdateMailDTO{
    Id:            request.MailId,
    AddedTagIds:   request.AddedTagIds,
    RemovedTagIds: request.RemovedTagIds,
  }
}

func ToProtoMail(dto *dto.MailResponseDTO) *mailerv1.Mail {
  return &mailerv1.Mail{
    Id:        dto.Id,
    Subject:   dto.Subject,
    Recipient: dto.Recipient,
    Sender:    dto.Sender,
    Status:    dto.Status,
    Tags:      util.CastSliceP(dto.Tags, ToProtoTag),
  }
}
