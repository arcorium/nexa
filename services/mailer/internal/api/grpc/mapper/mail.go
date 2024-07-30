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
  var fieldErrs []sharedErr.FieldError

  recipientEmails, ierr := sharedUtil.CastSliceErrs(request.Recipients, types.EmailFromString)
  if !ierr.IsNil() {
    fieldErrs = append(fieldErrs, sharedErr.NewFieldError("recipients", ierr))
  }

  tagIds, ierr := sharedUtil.CastSliceErrs(request.TagIds, types.IdFromString)
  if !ierr.IsNil() && !ierr.IsEmptySlice() {
    fieldErrs = append(fieldErrs, sharedErr.NewFieldError("tag_ids", ierr))
  }

  bodyType, err := util.ToDomainBodyType(request.BodyType)
  if err != nil {
    fieldErrs = append(fieldErrs, sharedErr.NewFieldError("body_type", err))
  }

  var senderEmail *types.Email = nil
  if request.Sender != nil {
    email, err := types.EmailFromString(*request.Sender)
    if err != nil {
      fieldErrs = append(fieldErrs, sharedErr.NewFieldError("sender", err))
    }
    senderEmail = &email // dangling and escaped
  }

  fileIds, ierr := sharedUtil.CastSliceErrs(request.AttachmentFileIds, types.IdFromString)
  if !ierr.IsNil() && !ierr.IsEmptySlice() {
    fieldErrs = append(fieldErrs, sharedErr.NewFieldError("attachment_file_ids", ierr))
  }

  if len(fieldErrs) > 0 {
    return dto.SendMailDTO{}, sharedErr.GrpcFieldErrors2(fieldErrs...)
  }

  dtos := dto.SendMailDTO{
    Subject:           request.Subject,
    Recipients:        recipientEmails,
    Sender:            types.NewNullable(senderEmail),
    BodyType:          bodyType,
    Body:              request.Body,
    TagIds:            tagIds,
    AttachmentFileIds: fileIds,
  }

  //err = sharedUtil.ValidateStruct(&dtos)
  return dtos, err
}

func ToUpdateMailDTO(request *mailerv1.UpdateMailRequest) (dto.UpdateMailDTO, error) {
  var fieldErrs []sharedErr.FieldError
  // Mapping and Validation
  mailId, err := types.IdFromString(request.MailId)
  if err != nil {
    fieldErrs = append(fieldErrs, sharedErr.NewFieldError("mail_id", err))
  }

  appendTagIds, ierr := sharedUtil.CastSliceErrs(request.AddedTagIds, types.IdFromString)
  if !ierr.IsNil() && !ierr.IsEmptySlice() {
    fieldErrs = append(fieldErrs, sharedErr.NewFieldError("added_tag_ids", ierr))
  }

  removedTagIds, ierr := sharedUtil.CastSliceErrs(request.RemovedTagIds, types.IdFromString)
  if !ierr.IsNil() && !ierr.IsEmptySlice() {
    fieldErrs = append(fieldErrs, sharedErr.NewFieldError("removed_tag_ids", ierr))
  }

  if len(fieldErrs) > 0 {
    return dto.UpdateMailDTO{}, sharedErr.GrpcFieldErrors2(fieldErrs...)
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
