package service

import (
  "context"
  sharedDto "github.com/arcorium/nexa/shared/dto"
  "github.com/arcorium/nexa/shared/status"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "go.opentelemetry.io/otel/trace"
  "nexa/services/authorization/constant"
  "nexa/services/authorization/internal/domain/dto"
  "nexa/services/authorization/internal/domain/mapper"
  "nexa/services/authorization/internal/domain/repository"
  "nexa/services/authorization/internal/domain/service"
  "nexa/services/authorization/util"
)

func NewRole(role repository.IRole) service.IRole {
  return &roleService{
    roleRepo: role,
    tracer:   util.GetTracer(),
  }
}

type roleService struct {
  roleRepo repository.IRole
  tracer   trace.Tracer
}

func (r *roleService) FindByIds(ctx context.Context, roleIds ...types.Id) ([]dto.RoleResponseDTO, status.Object) {
  ctx, span := r.tracer.Start(ctx, "RoleService.FindByIds")
  defer span.End()

  roles, err := r.roleRepo.FindByIds(ctx, roleIds...)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, status.FromRepository(err, status.NullCode)
  }

  roleResponses := sharedUtil.CastSliceP(roles, mapper.ToRoleResponseDTO)
  return roleResponses, status.Success()
}

func (r *roleService) FindByUserId(ctx context.Context, userId types.Id) ([]dto.RoleResponseDTO, status.Object) {
  ctx, span := r.tracer.Start(ctx, "RoleService.FindByUserId")
  defer span.End()

  roles, err := r.roleRepo.FindByUserId(ctx, userId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, status.FromRepository(err, status.NullCode)
  }

  roleResponses := sharedUtil.CastSliceP(roles, mapper.ToRoleResponseDTO)
  return roleResponses, status.Success()
}

func (r *roleService) GetAll(ctx context.Context, input *sharedDto.PagedElementDTO) (sharedDto.PagedElementResult[dto.RoleResponseDTO], status.Object) {
  ctx, span := r.tracer.Start(ctx, "RoleService.GetAll")
  defer span.End()

  result, err := r.roleRepo.Get(ctx, input.ToQueryParam())
  if err != nil {
    spanUtil.RecordError(err, span)
    return sharedDto.PagedElementResult[dto.RoleResponseDTO]{}, status.FromRepository(err, status.NullCode)
  }

  rolesDTO := sharedUtil.CastSliceP(result.Data, mapper.ToRoleResponseDTO)
  return sharedDto.NewPagedElementResult2(rolesDTO, input, result.Total), status.Success()
}

func (r *roleService) Create(ctx context.Context, createDTO *dto.RoleCreateDTO) (types.Id, status.Object) {
  ctx, span := r.tracer.Start(ctx, "RoleService.Create")
  defer span.End()

  // Map
  role, err := createDTO.ToDomain()
  if err != nil {
    spanUtil.RecordError(err, span)
    return types.NullId(), status.ErrInternal(err)
  }

  err = r.roleRepo.Create(ctx, &role)
  if err != nil {
    spanUtil.RecordError(err, span)
    return types.NullId(), status.FromRepositoryExist(err)
  }
  return role.Id, status.Created()
}

func (r *roleService) Update(ctx context.Context, updateDTO *dto.RoleUpdateDTO) status.Object {
  ctx, span := r.tracer.Start(ctx, "RoleService.Update")
  defer span.End()

  role := updateDTO.ToDomain()

  err := r.roleRepo.Patch(ctx, &role)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }
  return status.Updated()
}

func (r *roleService) Delete(ctx context.Context, roleId types.Id) status.Object {
  ctx, span := r.tracer.Start(ctx, "RoleService.Delete")
  defer span.End()

  err := r.roleRepo.Delete(ctx, roleId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }

  return status.Deleted()
}

func (r *roleService) AddPermissions(ctx context.Context, permissionsDTO *dto.ModifyRolesPermissionsDTO) status.Object {
  ctx, span := r.tracer.Start(ctx, "RoleService.AddPermissions")
  defer span.End()

  err := r.roleRepo.AddPermissions(ctx, permissionsDTO.RoleId, permissionsDTO.PermissionIds...)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepositoryExist(err)
  }

  return status.Created()
}

func (r *roleService) RemovePermissions(ctx context.Context, permissionsDTO *dto.ModifyRolesPermissionsDTO) status.Object {
  ctx, span := r.tracer.Start(ctx, "RoleService.RemovePermissions")
  defer span.End()

  err := r.roleRepo.RemovePermissions(ctx, permissionsDTO.RoleId, permissionsDTO.PermissionIds...)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }

  return status.Deleted()
}

func (r *roleService) AddUsers(ctx context.Context, usersDTO *dto.ModifyUserRolesDTO) status.Object {
  ctx, span := r.tracer.Start(ctx, "RoleService.AddUsers")
  defer span.End()

  err := r.roleRepo.AddUser(ctx, usersDTO.UserId, usersDTO.RoleIds...)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepositoryExist(err)
  }

  return status.Created()
}

func (r *roleService) RemoveUsers(ctx context.Context, usersDTO *dto.ModifyUserRolesDTO) status.Object {
  ctx, span := r.tracer.Start(ctx, "RoleService.RemoveUsers")
  defer span.End()

  err := r.roleRepo.RemoveUser(ctx, usersDTO.UserId, usersDTO.RoleIds...)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }

  return status.Deleted()
}

func (r *roleService) AppendSuperRolesPermission(ctx context.Context, permIds ...types.Id) status.Object {
  ctx, span := r.tracer.Start(ctx, "PermissionService.AppendSuperRolesPermission")
  defer span.End()

  // Find role by names
  role, err := r.roleRepo.FindByName(ctx, constant.DEFAULT_SUPER_ROLE_NAME)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }

  // Append permission into it
  err = r.roleRepo.AddPermissions(ctx, role.Id, permIds...)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepositoryExist(err)
  }

  return status.Created()
}

func (r *roleService) SetUserAsSuper(ctx context.Context, userId types.Id) status.Object {
  ctx, span := r.tracer.Start(ctx, "PermissionService.SetUserAsSuper")
  defer span.End()

  // Get super roles
  role, err := r.roleRepo.FindByName(ctx, constant.DEFAULT_SUPER_ROLE_NAME)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }

  // Set the user with the roles
  err = r.roleRepo.AddUser(ctx, userId, role.Id)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }

  return status.Success()
}
