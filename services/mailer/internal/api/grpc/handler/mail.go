package handler

import (
  "context"
  "errors"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/genproto/googleapis/rpc/errdetails"
  "google.golang.org/grpc"
  "google.golang.org/protobuf/types/known/emptypb"
  "nexa/proto/gen/go/common"
  "nexa/proto/gen/go/mailer/v1"
  "nexa/services/mailer/internal/api/grpc/mapper"
  "nexa/services/mailer/internal/domain/service"
  "nexa/shared/dto"
  sharedErr "nexa/shared/errors"
  spanUtil "nexa/shared/span"
  "nexa/shared/types"
  sharedUtil "nexa/shared/util"
)

func NewMail(mail service.IMail) MailHandler {
  return MailHandler{
    mailService: mail,
  }
}

type MailHandler struct {
  mailerv1.UnimplementedMailerServiceServer
  mailService service.IMail
}

func (m *MailHandler) Register(server *grpc.Server) {
  mailerv1.RegisterMailerServiceServer(server, m)
}

func (m *MailHandler) Find(ctx context.Context, input *common.PagedElementInput) (*mailerv1.FindResponse, error) {
  span := trace.SpanFromContext(ctx)

  elementDTO := dto.PagedElementDTO{
    Element: input.Element,
    Page:    input.Page,
  }

  result, stat := m.mailService.Find(ctx, &elementDTO)
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
  span := trace.SpanFromContext(ctx)

  ids, ierr := sharedUtil.CastSliceErrs(request.MailIds, types.IdFromString)
  if ierr != nil {
    errs := sharedUtil.CastSlice(ierr, func(from sharedErr.IndexedError) error {
      return from.Err
    })
    err := errors.Join(errs...)
    spanUtil.RecordError(err, span)
    return nil, sharedErr.GrpcFieldIndexedErrors("mail_ids", ierr)
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
  span := trace.SpanFromContext(ctx)

  id, err := types.IdFromString(request.TagId)
  if err != nil {
    if errors.Is(err, types.ErrMalformedUUID) {
      return nil, sharedErr.GrpcFieldErrors(&errdetails.BadRequest_FieldViolation{
        Field:       "tag_id",
        Description: err.Error(),
      })
    }
    return nil, err
  }

  mails, stat := m.mailService.FindByTag(ctx, id)
  if stat.IsError() {
    spanUtil.RecordError(err, span)
    return nil, stat.ToGRPCError()
  }

  resp := &mailerv1.FindMailByTagResponse{
    Mails: sharedUtil.CastSliceP(mails, mapper.ToProtoMail),
  }

  return resp, nil
}

func (m *MailHandler) Send(ctx context.Context, request *mailerv1.SendMailRequest) (*mailerv1.SendMailResponse, error) {
  span := trace.SpanFromContext(ctx)

  dtos := mapper.ToSendMailDTO(request)
  ids, stat := m.mailService.Send(ctx, &dtos)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  resp := &mailerv1.SendMailResponse{
    MailIds: sharedUtil.CastSlice(ids, func(id types.Id) string { return id.Underlying().String() }),
  }
  return resp, nil
}

func (m *MailHandler) Update(ctx context.Context, request *mailerv1.UpdateMailRequest) (*emptypb.Empty, error) {
  span := trace.SpanFromContext(ctx)

  updateDto := mapper.ToUpdateMailDTO(request)
  err := sharedUtil.ValidateStructCtx(ctx, &updateDto)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  stat := m.mailService.Update(ctx, &updateDto)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  return &emptypb.Empty{}, nil
}

func (m *MailHandler) Remove(ctx context.Context, request *mailerv1.RemoveMailRequest) (*emptypb.Empty, error) {
  span := trace.SpanFromContext(ctx)

  id, err := types.IdFromString(request.MailId)
  if err != nil {
    if errors.Is(err, types.ErrMalformedUUID) {
      return nil, sharedErr.GrpcFieldErrors(&errdetails.BadRequest_FieldViolation{
        Field:       "mail_id",
        Description: err.Error(),
      })
    }
    return nil, err
  }

  stat := m.mailService.Remove(ctx, id)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  return &emptypb.Empty{}, nil
}
