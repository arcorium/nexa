package handler

import (
  "context"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  "google.golang.org/protobuf/types/known/emptypb"
  "nexa/proto/gen/go/common"
  "nexa/proto/gen/go/mailer/v1"
  "nexa/services/mailer/internal/api/grpc/mapper"
  "nexa/services/mailer/internal/domain/service"
  "nexa/services/mailer/util"
  "nexa/shared/dto"
  sharedErr "nexa/shared/errors"
  "nexa/shared/types"
  sharedUtil "nexa/shared/util"
  spanUtil "nexa/shared/util/span"
)

func NewMail(mail service.IMail) MailHandler {
  return MailHandler{
    mailService: mail,
    tracer:      util.GetTracer(),
  }
}

type MailHandler struct {
  mailerv1.UnimplementedMailerServiceServer
  mailService service.IMail

  tracer trace.Tracer
}

func (m *MailHandler) Register(server *grpc.Server) {
  mailerv1.RegisterMailerServiceServer(server, m)
}

func (m *MailHandler) Find(ctx context.Context, input *common.PagedElementInput) (*mailerv1.FindResponse, error) {
  ctx, span := m.tracer.Start(ctx, "MailerHandler.GetAll")
  defer span.End()

  elementDTO := dto.PagedElementDTO{
    Element: input.Element,
    Page:    input.Page,
  }

  result, stat := m.mailService.GetAll(ctx, &elementDTO)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  resp := &mailerv1.FindResponse{
    Details: &common.PagedElementOutput{
      Element:       result.Element,
      Page:          result.Page,
      TotalElements: result.TotalElements,
      TotalPages:    result.Page,
    },
    Mails: sharedUtil.CastSliceP(result.Data, mapper.ToProtoMail),
  }

  return resp, nil
}

func (m *MailHandler) FindByIds(ctx context.Context, request *mailerv1.FindMailByIdsRequest) (*mailerv1.FindMailByIdsResponse, error) {
  ctx, span := m.tracer.Start(ctx, "MailerHandler.GetAll")
  defer span.End()

  // Input type validation
  ids, ierr := sharedUtil.CastSliceErrs(request.MailIds, types.IdFromString)
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return nil, ierr.ToGRPCError("mail_ids")
  }

  mails, stat := m.mailService.FindByIds(ctx, ids...)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  resp := &mailerv1.FindMailByIdsResponse{
    Mails: sharedUtil.CastSliceP(mails, mapper.ToProtoMail),
  }
  return resp, nil
}

func (m *MailHandler) FindByTag(ctx context.Context, request *mailerv1.FindMailByTagRequest) (*mailerv1.FindMailByTagResponse, error) {
  ctx, span := m.tracer.Start(ctx, "MailerHandler.FindByTag")
  defer span.End()

  // Input type validation
  tagId, err := types.IdFromString(request.TagId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, sharedErr.NewFieldError("tag_id", err).ToGrpcError()
  }

  mails, stat := m.mailService.FindByTag(ctx, tagId)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  resp := &mailerv1.FindMailByTagResponse{
    Mails: sharedUtil.CastSliceP(mails, mapper.ToProtoMail),
  }

  return resp, nil
}

func (m *MailHandler) Send(ctx context.Context, request *mailerv1.SendMailRequest) (*mailerv1.SendMailResponse, error) {
  ctx, span := m.tracer.Start(ctx, "MailerHandler.Send")
  defer span.End()

  dtos, err := mapper.ToSendMailDTO(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  ids, stat := m.mailService.Send(ctx, &dtos)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  resp := &mailerv1.SendMailResponse{
    MailIds: sharedUtil.CastSlice(ids, sharedUtil.ToString[types.Id]),
  }
  return resp, nil
}

func (m *MailHandler) Update(ctx context.Context, request *mailerv1.UpdateMailRequest) (*emptypb.Empty, error) {
  ctx, span := m.tracer.Start(ctx, "MailerHandler.Update")
  defer span.End()

  updateDto, err := mapper.ToUpdateMailDTO(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  stat := m.mailService.Update(ctx, &updateDto)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (m *MailHandler) Remove(ctx context.Context, request *mailerv1.RemoveMailRequest) (*emptypb.Empty, error) {
  ctx, span := m.tracer.Start(ctx, "MailerHandler.Remove")
  defer span.End()

  mailId, err := types.IdFromString(request.MailId)
  if err != nil {
    spanUtil.RecordError(err, span)
  }

  stat := m.mailService.Remove(ctx, mailId)
  return nil, stat.ToGRPCErrorWithSpan(span)
}
