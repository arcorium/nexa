package handler

import (
  "context"
  authZv1 "github.com/arcorium/nexa/proto/gen/go/authorization/v1"
  "github.com/arcorium/nexa/proto/gen/go/common"
  sharedDto "github.com/arcorium/nexa/shared/dto"
  sharedErr "github.com/arcorium/nexa/shared/errors"
  "github.com/arcorium/nexa/shared/status"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  "google.golang.org/grpc/codes"
  grpcStat "google.golang.org/grpc/status"
  "google.golang.org/protobuf/types/known/emptypb"
  "nexa/services/authorization/internal/api/grpc/mapper"
  "nexa/services/authorization/internal/domain/dto"
  "nexa/services/authorization/internal/domain/service"
  "nexa/services/authorization/util"
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

  result, stat := r.roleService.GetAll(ctx, &pagedDto)
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

  removeDTO, err := mapper.ToRemoveUsersDTO(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  stat := r.roleService.RemoveUsers(ctx, &removeDTO)
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

func (r *RoleHandler) AppendDefaultRolePermissions(ctx context.Context, request *authZv1.AppendDefaultRolePermissionsRequest) (*emptypb.Empty, error) {
  ctx, span := r.tracer.Start(ctx, "RoleHandler.AppendDefaultRolePermissions")
  defer span.End()

  permIds, ierr := sharedUtil.CastSliceErrs(request.PermissionIds, types.IdFromString)
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return nil, ierr.ToGRPCError("permission_ids")
  }

  if request.Role == authZv1.DefaultRole_DEFAULT_ROLE {
    stat := r.roleService.AppendDefaultRolesPermission(ctx, permIds...)
    return nil, stat.ToGRPCErrorWithSpan(span)
  } else if request.Role == authZv1.DefaultRole_SUPER_ROLE {
    stat := r.roleService.AppendSuperRolesPermission(ctx, permIds...)
    return nil, stat.ToGRPCErrorWithSpan(span)
  }

  err := grpcStat.New(codes.InvalidArgument, "invalid role").Err()
  spanUtil.RecordError(err, span)
  return nil, err
}

func (r *RoleHandler) GetDefault(ctx context.Context, request *authZv1.GetDefaultRoleRequest) (*authZv1.GetDefaultRoleResponse, error) {
  ctx, span := r.tracer.Start(ctx, "RoleHandler.GetDefault")
  defer span.End()

  var res dto.RoleResponseDTO
  var stat status.Object
  if request.Role == authZv1.DefaultRole_DEFAULT_ROLE {
    res, stat = r.roleService.GetDefault(ctx)
  } else if request.Role == authZv1.DefaultRole_SUPER_ROLE {
    res, stat = r.roleService.GetSuper(ctx)
  } else {
    err := grpcStat.New(codes.InvalidArgument, "invalid role").Err()
    spanUtil.RecordError(err, span)
    return nil, err
  }

  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  permission := mapper.ToProtoRolePermission(&res, request.IncludePermissions)
  return &authZv1.GetDefaultRoleResponse{Role: permission}, nil
}

func (r *RoleHandler) SetAsDefault(ctx context.Context, request *authZv1.SetAsDefaultRequest) (*emptypb.Empty, error) {
  ctx, span := r.tracer.Start(ctx, "RoleHandler.SetAsDefault")
  defer span.End()

  userId, err := types.IdFromString(request.UserId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, sharedErr.NewFieldError("user_id", err).ToGrpcError()
  }

  stat := r.roleService.SetUserAsDefault(ctx, userId)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (r *RoleHandler) SetAsSuper(ctx context.Context, request *authZv1.SetAsSuperRequest) (*emptypb.Empty, error) {
  ctx, span := r.tracer.Start(ctx, "RoleHandler.SetAsSuper")
  defer span.End()

  userId, err := types.IdFromString(request.UserId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, sharedErr.NewFieldError("user_id", err).ToGrpcError()
  }

  stat := r.roleService.SetUserAsSuper(ctx, userId)
  return nil, stat.ToGRPCErrorWithSpan(span)
}
