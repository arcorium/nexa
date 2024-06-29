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
  sharedDto "nexa/shared/dto"
  spanUtil "nexa/shared/span"
  sharedUtil "nexa/shared/util"
)

func NewRole(role service.IRole) RoleHandler {
  return RoleHandler{roleSvc: role}
}

type RoleHandler struct {
  authZv1.UnimplementedRoleServiceServer
  roleSvc service.IRole
}

func (r *RoleHandler) Register(server *grpc.Server) {
  authZv1.RegisterRoleServiceServer(server, r)
}

func (r *RoleHandler) Create(ctx context.Context, request *authZv1.RoleCreateRequest) (*authZv1.RoleCreateResponse, error) {
  span := trace.SpanFromContext(ctx)

  dtos := mapper.ToRoleCreateDTO(request)
  id, stat := r.roleSvc.Create(ctx, &dtos)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  return &authZv1.RoleCreateResponse{Id: id}, nil
}

func (r *RoleHandler) Update(ctx context.Context, request *authZv1.RoleUpdateRequest) (*emptypb.Empty, error) {
  span := trace.SpanFromContext(ctx)

  dtos := mapper.ToRoleUpdateDTO(request)
  stat := r.roleSvc.Update(ctx, &dtos)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (r *RoleHandler) Delete(ctx context.Context, request *authZv1.RoleDeleteRequest) (*emptypb.Empty, error) {
  span := trace.SpanFromContext(ctx)

  stat := r.roleSvc.Delete(ctx, request.Id)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (r *RoleHandler) GetUser(ctx context.Context, request *authZv1.GetUserRolesRequest) (*authZv1.GetUserRolesResponse, error) {
  span := trace.SpanFromContext(ctx)

  roles, stat := r.roleSvc.FindByUserId(ctx, request.UserId)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  response := sharedUtil.CastSliceP(roles, func(role *dto.RoleResponseDTO) *authZv1.RolePermission {
    return mapper.ToProtoRolePermission(role, request.IncludePermission)
  })

  return &authZv1.GetUserRolesResponse{RolePermissions: response}, nil
}

func (r *RoleHandler) Find(ctx context.Context, request *authZv1.RoleFindRequest) (*authZv1.RoleFindResponse, error) {
  span := trace.SpanFromContext(ctx)

  roles, stat := r.roleSvc.Find(ctx, request.RoleIds...)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  response := sharedUtil.CastSliceP(roles, mapper.ToRoleResponse)
  return &authZv1.RoleFindResponse{Roles: response}, nil
}

func (r *RoleHandler) FindAll(ctx context.Context, input *common.PagedElementInput) (*authZv1.RoleFindAllResponse, error) {
  span := trace.SpanFromContext(ctx)

  pagedDto := sharedDto.PagedElementDTO{
    Element: input.Element,
    Page:    input.Page,
  }

  result, stat := r.roleSvc.FindAll(ctx, &pagedDto)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  response := &authZv1.RoleFindAllResponse{
    Details: &common.PagedElementOutput{
      Element:       result.Element,
      Page:          result.Page,
      TotalElements: result.TotalElements,
      TotalPages:    result.TotalPages,
    },
    Roles: sharedUtil.CastSliceP(result.Data, mapper.ToRoleResponse),
  }
  return response, nil
}

func (r *RoleHandler) AddUser(ctx context.Context, request *authZv1.AddUserRolesRequest) (*emptypb.Empty, error) {
  span := trace.SpanFromContext(ctx)

  dtos := mapper.ToAddUsersDTO(request)
  stat := r.roleSvc.AddUsers(ctx, &dtos)

  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (r *RoleHandler) RemoveUser(ctx context.Context, request *authZv1.RemoveUserRolesRequest) (*emptypb.Empty, error) {
  span := trace.SpanFromContext(ctx)

  dtos := mapper.ToRemoveUsersDTO(request)
  stat := r.roleSvc.RemoveUsers(ctx, &dtos)

  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (r *RoleHandler) AppendPermissions(ctx context.Context, request *authZv1.RoleAppendPermissionsRequest) (*emptypb.Empty, error) {
  span := trace.SpanFromContext(ctx)

  dtos := mapper.ToAddRolePermissionsDTO(request)
  stat := r.roleSvc.AddPermissions(ctx, &dtos)

  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (r *RoleHandler) RemovePermissions(ctx context.Context, request *authZv1.RoleRemovePermissionsRequest) (*emptypb.Empty, error) {
  span := trace.SpanFromContext(ctx)

  dtos := mapper.ToRemoveRolePermissionsDTO(request)
  stat := r.roleSvc.RemovePermissions(ctx, &dtos)

  return nil, stat.ToGRPCErrorWithSpan(span)
}
