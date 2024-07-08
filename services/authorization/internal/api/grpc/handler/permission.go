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
  "nexa/services/authorization/util"
  sharedDto "nexa/shared/dto"
  sharedErr "nexa/shared/errors"
  "nexa/shared/types"
  sharedUtil "nexa/shared/util"
  spanUtil "nexa/shared/util/span"
)

func NewPermission(permission service.IPermission) PermissionHandler {
  return PermissionHandler{
    permService: permission,
    tracer:      util.GetTracer(),
  }
}

type PermissionHandler struct {
  authZv1.UnimplementedPermissionServiceServer

  permService service.IPermission
  tracer      trace.Tracer
}

func (p *PermissionHandler) Register(server *grpc.Server) {
  authZv1.RegisterPermissionServiceServer(server, p)
}

func (p *PermissionHandler) Create(ctx context.Context, request *authZv1.CreatePermissionRequest) (*authZv1.CreatePermissionResponse, error) {
  ctx, span := p.tracer.Start(ctx, "PermissionHandler.Create")
  defer span.End()

  dtos, err := mapper.ToCreatePermissionDTO(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  id, stat := p.permService.Create(ctx, &dtos)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  return &authZv1.CreatePermissionResponse{
    PermissionId: id.String(),
  }, nil
}

func (p *PermissionHandler) Find(ctx context.Context, request *authZv1.FindPermissionRequest) (*authZv1.FindPermissionResponse, error) {
  ctx, span := p.tracer.Start(ctx, "PermissionHandler.Find")
  defer span.End()

  permIds, ierr := sharedUtil.CastSliceErrs(request.PermIds, types.IdFromString)
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return nil, ierr.ToGRPCError("perm_ids")
  }

  permissions, stat := p.permService.Find(ctx, permIds...)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  response := sharedUtil.CastSliceP(permissions, mapper.ToProtoPermission)
  return &authZv1.FindPermissionResponse{Permission: response}, nil
}

func (p *PermissionHandler) FindByRoles(ctx context.Context, request *authZv1.FindPermissionsByRoleRequest) (*authZv1.FindPermissionByRoleResponse, error) {
  ctx, span := p.tracer.Start(ctx, "PermissionHandler.FindByRole")
  defer span.End()

  roleIds, ierr := sharedUtil.CastSliceErrs(request.RoleIds, types.IdFromString)
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return nil, ierr.ToGRPCError("role_ids")
  }

  perms, stat := p.permService.FindByRoles(ctx, roleIds...)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  return &authZv1.FindPermissionByRoleResponse{
    Permissions: sharedUtil.CastSliceP(perms, mapper.ToProtoPermission),
  }, nil
}

func (p *PermissionHandler) FindAll(ctx context.Context, input *common.PagedElementInput) (*authZv1.FindAllPermissionResponse, error) {
  ctx, span := p.tracer.Start(ctx, "PermissionHandler.Get")
  defer span.End()

  pagedDto := sharedDto.PagedElementDTO{
    Element: input.Element,
    Page:    input.Page,
  }

  result, stat := p.permService.FindAll(ctx, &pagedDto)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  resp := &authZv1.FindAllPermissionResponse{
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

func (p *PermissionHandler) Delete(ctx context.Context, request *authZv1.DeletePermissionRequest) (*emptypb.Empty, error) {
  span := trace.SpanFromContext(ctx)

  permId, err := types.IdFromString(request.Id)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, sharedErr.NewFieldError("id", err).ToGrpcError()
  }

  stat := p.permService.Delete(ctx, permId)
  return nil, stat.ToGRPCErrorWithSpan(span)
}
