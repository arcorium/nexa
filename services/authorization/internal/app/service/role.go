package service

import (
  "context"
  "go.opentelemetry.io/otel/trace"
  "nexa/services/authorization/internal/domain/dto"
  "nexa/services/authorization/internal/domain/mapper"
  "nexa/services/authorization/internal/domain/repository"
  "nexa/services/authorization/internal/domain/service"
  "nexa/services/authorization/util"
  sharedDto "nexa/shared/dto"
  spanUtil "nexa/shared/span"
  "nexa/shared/status"
  "nexa/shared/types"
  sharedUtil "nexa/shared/util"
)

func NewRole(role repository.IRole) service.IRole {
  return &roleService{
    repo:   role,
    tracer: util.GetTracer(),
  }
}

type roleService struct {
  repo   repository.IRole
  tracer trace.Tracer
}

func (r *roleService) Find(ctx context.Context, ids ...string) ([]dto.RoleResponseDTO, status.Object) {
  ctx, span := r.tracer.Start(ctx, "RoleService.Find")
  defer span.End()

  roleIds, ierr := sharedUtil.CastSliceErrs(ids, func(id string) (types.Id, error) {
    return types.IdFromString(id)
  })
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return nil, status.ErrBadRequest(ierr)
  }

  roles, err := r.repo.FindByIds(ctx, roleIds...)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, status.FromRepository(err, status.NullCode)
  }

  roleResponses := sharedUtil.CastSliceP(roles, mapper.ToRoleResponseDTO)
  return roleResponses, status.Success()
}

func (r *roleService) FindByUserId(ctx context.Context, userId string) ([]dto.RoleResponseDTO, status.Object) {
  ctx, span := r.tracer.Start(ctx, "RoleService.FindByUserId")
  defer span.End()

  id, err := types.IdFromString(userId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, status.ErrBadRequest(err)
  }

  roles, err := r.repo.FindByUserId(ctx, id)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, status.FromRepository(err, status.NullCode)
  }

  roleResponses := sharedUtil.CastSliceP(roles, mapper.ToRoleResponseDTO)
  return roleResponses, status.Success()
}

func (r *roleService) FindAll(ctx context.Context, input *sharedDto.PagedElementDTO) (sharedDto.PagedElementResult[dto.RoleResponseDTO], status.Object) {
  ctx, span := r.tracer.Start(ctx, "RoleService.FindAll")
  defer span.End()

  result, err := r.repo.FindAll(ctx, input.ToQueryParam())
  if err != nil {
    spanUtil.RecordError(err, span)
    return types.Null[sharedDto.PagedElementResult[dto.RoleResponseDTO]](), status.FromRepository(err, status.NullCode)
  }

  rolesDTO := sharedUtil.CastSliceP(result.Data, mapper.ToRoleResponseDTO)
  return sharedDto.NewPagedElementResult2(rolesDTO, input, result.Total), status.FromRepository(err, status.NullCode)
}

func (r *roleService) Create(ctx context.Context, createDTO *dto.RoleCreateDTO) (string, status.Object) {
  ctx, span := r.tracer.Start(ctx, "RoleService.Create")
  defer span.End()

  // Map and validate
  role, err := createDTO.ToDomain()
  if err != nil {
    spanUtil.RecordError(err, span)
    return "", status.ErrBadRequest(err)
  }

  err = r.repo.Create(ctx, &role)
  if err != nil {
    spanUtil.RecordError(err, span)
    return "", status.FromRepository(err, status.NullCode)
  }
  return role.Id.String(), status.Created()
}

func (r *roleService) Update(ctx context.Context, updateDTO *dto.RoleUpdateDTO) status.Object {
  ctx, span := r.tracer.Start(ctx, "RoleService.Update")
  defer span.End()

  role, err := updateDTO.ToDomain()
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrBadRequest(err)
  }

  err = r.repo.Patch(ctx, &role)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }
  return status.Updated()
}

func (r *roleService) Delete(ctx context.Context, roleId string) status.Object {
  ctx, span := r.tracer.Start(ctx, "RoleService.Delete")
  defer span.End()

  id, err := types.IdFromString(roleId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrBadRequest(err)
  }

  err = r.repo.Delete(ctx, id)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }

  return status.Deleted()
}

func (r *roleService) AddPermissions(ctx context.Context, permissionsDTO *dto.RoleAddPermissionsDTO) status.Object {
  ctx, span := r.tracer.Start(ctx, "RoleService.AddPermissions")
  defer span.End()

  // Input validation
  roleId, err := types.IdFromString(permissionsDTO.RoleId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrBadRequest(err)
  }

  permissionIds, ierr := sharedUtil.CastSliceErrs(permissionsDTO.PermissionIds, func(permId string) (types.Id, error) {
    return types.IdFromString(permId)
  })
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return status.ErrBadRequest(ierr)
  }

  err = r.repo.AddPermissions(ctx, roleId, permissionIds...)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }

  return status.Created()
}

func (r *roleService) RemovePermissions(ctx context.Context, permissionsDTO *dto.RoleRemovePermissionsDTO) status.Object {
  ctx, span := r.tracer.Start(ctx, "RoleService.RemovePermissions")
  defer span.End()

  // Input validation
  roleId, err := types.IdFromString(permissionsDTO.RoleId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrBadRequest(err)
  }

  permissionIds, ierr := sharedUtil.CastSliceErrs(permissionsDTO.PermissionIds, func(permId string) (types.Id, error) {
    return types.IdFromString(permId)
  })
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return status.ErrBadRequest(ierr)
  }

  err = r.repo.RemovePermissions(ctx, roleId, permissionIds...)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }

  return status.Deleted()
}

func (r *roleService) AddUsers(ctx context.Context, usersDTO *dto.RoleAddUsersDTO) status.Object {
  ctx, span := r.tracer.Start(ctx, "RoleService.AddUsers")
  defer span.End()

  // Input validation
  userId, err := types.IdFromString(usersDTO.UserId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrBadRequest(err)
  }

  roleIds, ierr := sharedUtil.CastSliceErrs(usersDTO.RoleIds, func(roleId string) (types.Id, error) {
    return types.IdFromString(roleId)
  })
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return status.ErrBadRequest(ierr)
  }

  err = r.repo.AddUser(ctx, userId, roleIds...)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }

  return status.Created()
}

func (r *roleService) RemoveUsers(ctx context.Context, usersDTO *dto.RoleRemoveUsersDTO) status.Object {
  ctx, span := r.tracer.Start(ctx, "RoleService.RemoveUsers")
  defer span.End()

  // Input validation
  userId, err := types.IdFromString(usersDTO.UserId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrBadRequest(err)
  }

  roleIds, ierr := sharedUtil.CastSliceErrs(usersDTO.RoleIds, func(roleId string) (types.Id, error) {
    return types.IdFromString(roleId)
  })
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return status.ErrBadRequest(ierr)
  }

  err = r.repo.RemoveUser(ctx, userId, roleIds...)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }

  return status.Deleted()
}
