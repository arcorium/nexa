package mapper

import (
  mailerv1 "github.com/arcorium/nexa/proto/gen/go/mailer/v1"
  sharedErr "github.com/arcorium/nexa/shared/errors"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  "google.golang.org/protobuf/types/known/timestamppb"
  "nexa/services/mailer/internal/domain/dto"
  "nexa/services/mailer/util"
)

func ToSendMailDTO(request *mailerv1.SendMailRequest) (dto.SendMailDTO, error) {
  recipientEmails, ierr := sharedUtil.CastSliceErrs(request.Recipients, types.EmailFromString)
  if !ierr.IsNil() {
    return dto.SendMailDTO{}, ierr.ToGRPCError("recipients")
  }

  tagIds, ierr := sharedUtil.CastSliceErrs(request.TagIds, types.IdFromString)
  if !ierr.IsNil() {
    return dto.SendMailDTO{}, ierr.ToGRPCError("tag_ids")
  }

  bodyType, err := util.ToDomainBodyType(request.BodyType)
  if err != nil {
    return dto.SendMailDTO{}, sharedErr.NewFieldError("body_type", err).ToGrpcError()
  }

  var senderEmail *types.Email = nil
  if request.Sender != nil {
    email, err := types.EmailFromString(*request.Sender)
    if err != nil {
      return dto.SendMailDTO{}, sharedErr.NewFieldError("sender", err).ToGrpcError()
    }
    senderEmail = &email // dangling and escaped
  }

  dtos := dto.SendMailDTO{
    Subject:    request.Subject,
    Recipients: recipientEmails,
    Sender:     types.NewNullable(senderEmail),
    BodyType:   bodyType,
    Body:       request.Body,
    TagIds:     tagIds,
    Attachments: sharedUtil.CastSlice(request.Attachments, func(attachment *mailerv1.FileAttachment) dto.FileAttachment {
      return dto.FileAttachment{
        Filename: attachment.Filename,
        Data:     attachment.Data,
      }
    }),
  }

  err = sharedUtil.ValidateStruct(&dtos)
  return dtos, err
}

func ToUpdateMailDTO(request *mailerv1.UpdateMailRequest) (dto.UpdateMailDTO, error) {
  // Mapping and Validation
  mailId, err := types.IdFromString(request.MailId)
  if err != nil {
    return dto.UpdateMailDTO{}, sharedErr.NewFieldError("mail_id", err).ToGrpcError()
  }

  appendTagIds, ierr := sharedUtil.CastSliceErrs(request.AddedTagIds, types.IdFromString)
  if !ierr.IsNil() {
    return dto.UpdateMailDTO{}, ierr.ToGRPCError("added_tag_ids")
  }

  removedTagIds, ierr := sharedUtil.CastSliceErrs(request.RemovedTagIds, types.IdFromString)
  if !ierr.IsNil() {
    return dto.UpdateMailDTO{}, ierr.ToGRPCError("removed_tag_ids")
  }

  return dto.UpdateMailDTO{
    Id:            mailId,
    AddedTagIds:   appendTagIds,
    RemovedTagIds: removedTagIds,
  }, nil
}

func ToProtoMail(dto *dto.MailResponseDTO) *mailerv1.Mail {
  return &mailerv1.Mail{
    Id:          dto.Id.String(),
    Subject:     dto.Subject,
    Recipient:   dto.Recipient.String(),
    Sender:      dto.Sender.String(),
    Status:      dto.Status.String(),
    SentAt:      timestamppb.New(dto.SentAt),
    DeliveredAt: timestamppb.New(dto.DeliveredAt),
    Tags:        sharedUtil.CastSliceP(dto.Tags, ToProtoTag),
  }
}
