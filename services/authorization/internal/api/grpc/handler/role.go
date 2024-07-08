package handler

import (
  "context"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  "google.golang.org/protobuf/types/known/emptypb"
  authZv1 "nexa/proto/gen/go/authorization/v1"
  "nexa/proto/gen/go/common"
  "nexa/services/authorization/internal/api/grpc/mapper"
  "nexa/services/authorization/internal/domain/dto"
  "nexa/services/authorization/internal/domain/service"
  "nexa/services/authorization/util"
  sharedDto "nexa/shared/dto"
  sharedErr "nexa/shared/errors"
  "nexa/shared/types"
  sharedUtil "nexa/shared/util"
  spanUtil "nexa/shared/util/span"
)

func NewRole(role service.IRole) RoleHandler {
  return RoleHandler{
    roleService: role,
    tracer:      util.GetTracer(),
  }
}

type RoleHandler struct {
  authZv1.UnimplementedRoleServiceServer
  roleService service.IRole

  tracer trace.Tracer
}

func (r *RoleHandler) Register(server *grpc.Server) {
  authZv1.RegisterRoleServiceServer(server, r)
}

func (r *RoleHandler) Create(ctx context.Context, request *authZv1.CreateRoleRequest) (*authZv1.RoleCreateResponse, error) {
  ctx, span := r.tracer.Start(ctx, "RoleHandler.Create")
  defer span.End()

  dtos, err := mapper.ToRoleCreateDTO(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  id, stat := r.roleService.Create(ctx, &dtos)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  return &authZv1.RoleCreateResponse{Id: id.String()}, nil
}

func (r *RoleHandler) Update(ctx context.Context, request *authZv1.UpdateRoleRequest) (*emptypb.Empty, error) {
  ctx, span := r.tracer.Start(ctx, "RoleHandler.Update")
  defer span.End()

  dtos, err := mapper.ToRoleUpdateDTO(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  stat := r.roleService.Update(ctx, &dtos)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (r *RoleHandler) Delete(ctx context.Context, request *authZv1.DeleteRoleRequest) (*emptypb.Empty, error) {
  ctx, span := r.tracer.Start(ctx, "RoleHandler.Delete")
  defer span.End()

  roleId, err := types.IdFromString(request.Id)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, sharedErr.NewFieldError("id", err).ToGrpcError()
  }

  stat := r.roleService.Delete(ctx, roleId)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

// GetUsers get user's roles
func (r *RoleHandler) GetUsers(ctx context.Context, request *authZv1.GetUserRolesRequest) (*authZv1.GetUserRolesResponse, error) {
  ctx, span := r.tracer.Start(ctx, "RoleHandler.GetUsers")
  defer span.End()

  userId, err := types.IdFromString(request.UserId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, sharedErr.NewFieldError("user_id", err).ToGrpcError()
  }

  roles, stat := r.roleService.FindByUserId(ctx, userId)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  response := sharedUtil.CastSliceP(roles, func(role *dto.RoleResponseDTO) *authZv1.RolePermission {
    return mapper.ToProtoRolePermission(role, request.IncludePermission)
  })

  return &authZv1.GetUserRolesResponse{RolePermissions: response}, nil
}

func (r *RoleHandler) Find(ctx context.Context, request *authZv1.FindRoleRequest) (*authZv1.FindRoleResponse, error) {
  ctx, span := r.tracer.Start(ctx, "RoleHandler.Find")
  defer span.End()

  roleIds, ierr := sharedUtil.CastSliceErrs(request.RoleIds, types.IdFromString)
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return nil, ierr.ToGRPCError("role_ids")
  }

  roles, stat := r.roleService.FindByIds(ctx, roleIds...)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  response := sharedUtil.CastSliceP(roles, mapper.ToProtoRole)
  return &authZv1.FindRoleResponse{Roles: response}, nil
}

func (r *RoleHandler) FindAll(ctx context.Context, input *common.PagedElementInput) (*authZv1.FindAllRolesResponse, error) {
  ctx, span := r.tracer.Start(ctx, "RoleHandler.Get")
  defer span.End()

  pagedDto := sharedDto.PagedElementDTO{
    Element: input.Element,
    Page:    input.Page,
  }

  result, stat := r.roleService.FindAll(ctx, &pagedDto)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  response := &authZv1.FindAllRolesResponse{
    Details: &common.PagedElementOutput{
      Element:       result.Element,
      Page:          result.Page,
      TotalElements: result.TotalElements,
      TotalPages:    result.TotalPages,
    },
    Roles: sharedUtil.CastSliceP(result.Data, mapper.ToProtoRole),
  }
  return response, nil
}

func (r *RoleHandler) AddUser(ctx context.Context, request *authZv1.AddUserRolesRequest) (*emptypb.Empty, error) {
  ctx, span := r.tracer.Start(ctx, "RoleHandler.AddUser")
  defer span.End()

  dtos, err := mapper.ToAddUsersDTO(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  stat := r.roleService.AddUsers(ctx, &dtos)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (r *RoleHandler) RemoveUser(ctx context.Context, request *authZv1.RemoveUserRolesRequest) (*emptypb.Empty, error) {
  ctx, span := r.tracer.Start(ctx, "RoleHandler.RemoveUser")
  defer span.End()

  dtos, err := mapper.ToRemoveUsersDTO(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  stat := r.roleService.RemoveUsers(ctx, &dtos)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

// AppendPermissions append permissions into role
func (r *RoleHandler) AppendPermissions(ctx context.Context, request *authZv1.AppendRolePermissionsRequest) (*emptypb.Empty, error) {
  ctx, span := r.tracer.Start(ctx, "RoleHandler.AppendPermissions")
  defer span.End()

  dtos, err := mapper.ToAddRolePermissionsDTO(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  stat := r.roleService.AddPermissions(ctx, &dtos)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (r *RoleHandler) RemovePermissions(ctx context.Context, request *authZv1.RemoveRolePermissionsRequest) (*emptypb.Empty, error) {
  ctx, span := r.tracer.Start(ctx, "RoleHandler.RemovePermissions")
  defer span.End()

  dtos, err := mapper.ToRemoveRolePermissionsDTO(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  stat := r.roleService.RemovePermissions(ctx, &dtos)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (p *RoleHandler) AppendSuperAdminPermissions(ctx context.Context, request *authZv1.AppendSuperAdminPermissionsRequest) (*emptypb.Empty, error) {
  ctx, span := p.tracer.Start(ctx, "RoleHandler.AppendSuperAdminPermissions")
  defer span.End()

  permIds, ierr := sharedUtil.CastSliceErrs(request.PermissionIds, types.IdFromString)
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return nil, ierr.ToGRPCError("permission_ids")
  }

  stat := p.roleService.AppendSuperRolesPermission(ctx, permIds...)
  return nil, stat.ToGRPCErrorWithSpan(span)
}
