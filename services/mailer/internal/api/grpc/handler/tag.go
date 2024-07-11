package handler

import (
  "context"
  "github.com/arcorium/nexa/proto/gen/go/common"
  mailerv1 "github.com/arcorium/nexa/proto/gen/go/mailer/v1"
  "github.com/arcorium/nexa/shared/dto"
  sharedErr "github.com/arcorium/nexa/shared/errors"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  "google.golang.org/protobuf/types/known/emptypb"
  "nexa/services/mailer/internal/api/grpc/mapper"
  "nexa/services/mailer/internal/domain/service"
  "nexa/services/mailer/util"
)

func NewTag(tag service.ITag) TagHandler {
  return TagHandler{
    tagService: tag,
    tracer:     util.GetTracer(),
  }
}

type TagHandler struct {
  mailerv1.UnimplementedTagServiceServer
  tagService service.ITag

  tracer trace.Tracer
}

func (t *TagHandler) Register(server *grpc.Server) {
  mailerv1.RegisterTagServiceServer(server, t)
}

func (t *TagHandler) Find(ctx context.Context, input *common.PagedElementInput) (*mailerv1.FindTagResponse, error) {
  ctx, span := t.tracer.Start(ctx, "TagHandler.GetAll")
  defer span.End()

  elementDto := dto.PagedElementDTO{
    Element: input.Element,
    Page:    input.Page,
  }

  result, stat := t.tagService.GetAll(ctx, &elementDto)
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
  ctx, span := t.tracer.Start(ctx, "TagHandler.FindByIds")
  defer span.End()

  // Validate
  tagIds, ierr := sharedUtil.CastSliceErrs(request.TagIds, types.IdFromString)
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return nil, ierr.ToGRPCError("tag_id")
  }

  tags, stat := t.tagService.FindByIds(ctx, tagIds...)
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
  ctx, span := t.tracer.Start(ctx, "TagHandler.FindByName")
  defer span.End()

  // Validate
  if err := sharedUtil.StringEmptyValidates(types.NewField("name", request.Name)); !err.IsNil() {
    spanUtil.RecordError(err, span)
    return nil, err.ToGRPCError()
  }

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
  ctx, span := t.tracer.Start(ctx, "TagHandler.Create")
  defer span.End()

  createDto, err := mapper.ToCreateTagDTO(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  id, stat := t.tagService.Create(ctx, &createDto)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  resp := &mailerv1.CreateTagResponse{
    TagId: id.String(),
  }

  return resp, nil
}

func (t *TagHandler) Update(ctx context.Context, request *mailerv1.UpdateTagRequest) (*emptypb.Empty, error) {
  ctx, span := t.tracer.Start(ctx, "TagHandler.Update")
  defer span.End()

  updateDto, err := mapper.ToUpdateTagDTO(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  stat := t.tagService.Update(ctx, &updateDto)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  return &emptypb.Empty{}, nil
}

func (t *TagHandler) Remove(ctx context.Context, request *mailerv1.RemoveTagRequest) (*emptypb.Empty, error) {
  ctx, span := t.tracer.Start(ctx, "TagHandler.Remove")
  defer span.End()

  id, err := types.IdFromString(request.TagId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, sharedErr.NewFieldError("tag_id", err).ToGrpcError()
  }

  stat := t.tagService.Remove(ctx, id)
  return nil, stat.ToGRPCErrorWithSpan(span)
}
