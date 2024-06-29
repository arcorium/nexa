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

func NewTag(tag service.ITag) TagHandler {
  return TagHandler{
    tagService: tag,
  }
}

type TagHandler struct {
  mailerv1.UnimplementedTagServiceServer
  tagService service.ITag
}

func (t *TagHandler) Register(server *grpc.Server) {
  mailerv1.RegisterTagServiceServer(server, t)
}

func (t *TagHandler) Find(ctx context.Context, input *common.PagedElementInput) (*mailerv1.FindTagResponse, error) {
  span := trace.SpanFromContext(ctx)

  elementDto := dto.PagedElementDTO{
    Element: input.Element,
    Page:    input.Page,
  }

  result, stat := t.tagService.Find(ctx, &elementDto)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  resp := &mailerv1.FindTagResponse{
    Details: &common.PagedElementOutput{
      Element:       result.Element,
      Page:          result.Page,
      TotalElements: result.TotalElements,
      TotalPages:    result.TotalPages,
    },
    Tags: sharedUtil.CastSliceP(result.Data, mapper.ToProtoTag),
  }

  return resp, nil
}

func (t *TagHandler) FindByIds(ctx context.Context, request *mailerv1.FindTagByIdsRequest) (*mailerv1.FindTagByIdsResponse, error) {
  span := trace.SpanFromContext(ctx)

  ids, ierr := sharedUtil.CastSliceErrs(request.TagIds, types.IdFromString)
  if ierr != nil {
    errs := sharedUtil.CastSlice(ierr, func(from sharedErr.IndexedError) error {
      return from.Err
    })
    err := errors.Join(errs...)
    spanUtil.RecordError(err, span)
    return nil, sharedErr.GrpcFieldIndexedErrors("tag_ids", ierr)
  }

  tags, stat := t.tagService.FindByIds(ctx, ids...)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  resp := &mailerv1.FindTagByIdsResponse{
    Tags: sharedUtil.CastSliceP(tags, mapper.ToProtoTag),
  }

  return resp, nil
}

func (t *TagHandler) FindByName(ctx context.Context, request *mailerv1.FindTagByNameRequest) (*mailerv1.FindTagByNameResponse, error) {
  span := trace.SpanFromContext(ctx)

  tag, stat := t.tagService.FindByName(ctx, request.Name)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  resp := &mailerv1.FindTagByNameResponse{
    Tag: mapper.ToProtoTag(&tag),
  }

  return resp, nil
}

func (t *TagHandler) Create(ctx context.Context, request *mailerv1.CreateTagRequest) (*mailerv1.CreateTagResponse, error) {
  span := trace.SpanFromContext(ctx)

  dto := mapper.ToCreateTagDTO(request)
  err := sharedUtil.ValidateStructCtx(ctx, &dto)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  id, stat := t.tagService.Create(ctx, &dto)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  resp := &mailerv1.CreateTagResponse{
    TagId: id.Underlying().String(),
  }

  return resp, nil
}

func (t *TagHandler) Update(ctx context.Context, request *mailerv1.UpdateTagRequest) (*emptypb.Empty, error) {
  span := trace.SpanFromContext(ctx)

  dto := mapper.ToUpdateTagDTO(request)
  err := sharedUtil.ValidateStructCtx(ctx, &dto)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  stat := t.tagService.Update(ctx, &dto)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  return &emptypb.Empty{}, nil
}

func (t *TagHandler) Remove(ctx context.Context, request *mailerv1.RemoveTagRequest) (*emptypb.Empty, error) {
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

  stat := t.tagService.Remove(ctx, id)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  return &emptypb.Empty{}, nil
}
