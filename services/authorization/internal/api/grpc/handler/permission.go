package handler

import (
  "context"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  "google.golang.org/protobuf/types/known/emptypb"
  authZv1 "nexa/proto/gen/go/authorization/v1"
  common "nexa/proto/gen/go/common"
  "nexa/services/authorization/internal/api/grpc/mapper"
  "nexa/services/authorization/internal/domain/service"
  sharedDto "nexa/shared/dto"
  spanUtil "nexa/shared/span"
  sharedUtil "nexa/shared/util"
)

func NewPermission(permission service.IPermission) PermissionHandler {
  return PermissionHandler{permSvc: permission}
}

type PermissionHandler struct {
  authZv1.UnimplementedPermissionServiceServer

  permSvc service.IPermission
}

func (p *PermissionHandler) Register(server *grpc.Server) {
  authZv1.RegisterPermissionServiceServer(server, p)
}

func (p *PermissionHandler) Create(ctx context.Context, request *authZv1.PermissionCreateRequest) (*authZv1.PermissionCreateResponse, error) {
  span := trace.SpanFromContext(ctx)

  dtos := mapper.ToPermissionCreateDTO(request)
  id, stat := p.permSvc.Create(ctx, &dtos)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  return &authZv1.PermissionCreateResponse{Id: id.String()}, nil
}

func (p *PermissionHandler) Find(ctx context.Context, request *authZv1.PermissionFindRequest) (*authZv1.FindPermissionResponse, error) {
  span := trace.SpanFromContext(ctx)

  permissions, stat := p.permSvc.Find(ctx, request.PermIds...)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  response := sharedUtil.CastSliceP(permissions, mapper.ToProtoPermission)
  return &authZv1.FindPermissionResponse{Permission: response}, nil
}

func (p *PermissionHandler) FindAll(ctx context.Context, input *common.PagedElementInput) (*authZv1.PermissionFindAllResponse, error) {
  span := trace.SpanFromContext(ctx)

  pagedDto := sharedDto.PagedElementDTO{
    Element: input.Element,
    Page:    input.Page,
  }

  result, stat := p.permSvc.FindAll(ctx, &pagedDto)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  resp := &authZv1.PermissionFindAllResponse{
    Details: &common.PagedElementOutput{
      Element:       result.Element,
      Page:          result.Page,
      TotalElements: result.TotalElements,
      TotalPages:    result.TotalPages,
    },
    Permissions: sharedUtil.CastSliceP(result.Data, mapper.ToProtoPermission),
  }

  return resp, nil
}

func (p *PermissionHandler) Delete(ctx context.Context, request *authZv1.PermissionDeleteRequest) (*emptypb.Empty, error) {
  span := trace.SpanFromContext(ctx)

  stat := p.permSvc.Delete(ctx, request.Id)
  return nil, stat.ToGRPCErrorWithSpan(span)
}
