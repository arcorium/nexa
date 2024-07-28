package service

import (
  "context"
  "database/sql"
  "errors"
  sharedDto "github.com/arcorium/nexa/shared/dto"
  sharedErr "github.com/arcorium/nexa/shared/errors"
  sharedJwt "github.com/arcorium/nexa/shared/jwt"
  "github.com/arcorium/nexa/shared/status"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  authUtil "github.com/arcorium/nexa/shared/util/auth"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "go.opentelemetry.io/otel/trace"
  "nexa/services/authorization/constant"
  "nexa/services/authorization/internal/domain/dto"
  "nexa/services/authorization/internal/domain/mapper"
  "nexa/services/authorization/internal/domain/repository"
  "nexa/services/authorization/internal/domain/service"
  "nexa/services/authorization/util"
  "nexa/services/authorization/util/errs"
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

func (r *roleService) checkPermission(ctx context.Context, targetId types.Id, permission string) error {
  // Validate permission
  claims := types.Must(sharedJwt.GetUserClaimsFromCtx(ctx))
  if !targetId.EqWithString(claims.UserId) {
    // Need permission to update other users
    if !authUtil.ContainsPermission(claims.Roles, permission) {
      return sharedErr.ErrUnauthorizedPermission
    }
  }
  return nil
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
  if err != nil && !errors.Is(err, sql.ErrNoRows) {
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

  var err error
  if permissionsDTO.PermissionIds == nil {
    err = r.roleRepo.ClearPermission(ctx, permissionsDTO.RoleId)
  } else {
    err = r.roleRepo.RemovePermissions(ctx, permissionsDTO.RoleId, permissionsDTO.PermissionIds...)
  }
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

  if err := r.checkPermission(ctx, usersDTO.UserId, constant.AUTHZ_PERMISSIONS[constant.AUTHZ_DELETE_USER_ROLE_ARB]); err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrUnAuthorized(err)
  }

  var err error
  if usersDTO.RoleIds == nil {
    err = r.roleRepo.ClearUser(ctx, usersDTO.UserId)
  } else {
    err = r.roleRepo.RemoveUser(ctx, usersDTO.UserId, usersDTO.RoleIds...)
  }
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }

  return status.Deleted()
}

func (r *roleService) getByName(ctx context.Context, roleName string) (dto.RoleResponseDTO, status.Object) {
  span := trace.SpanFromContext(ctx)

  role, err := r.roleRepo.FindByName(ctx, roleName)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.RoleResponseDTO{}, status.FromRepositoryOverride(err, types.NewPair(status.INTERNAL_SERVER_ERROR, errs.ErrDefaultRoleNotSeeded))
  }

  responseDTO := mapper.ToRoleResponseDTO(&role)
  return responseDTO, status.Success()
}

func (r *roleService) GetDefault(ctx context.Context) (dto.RoleResponseDTO, status.Object) {
  ctx, span := r.tracer.Start(ctx, "RoleService.GetDefault")
  defer span.End()

  return r.getByName(ctx, constant.DEFAULT_ROLE_NAME)
}

func (r *roleService) GetSuper(ctx context.Context) (dto.RoleResponseDTO, status.Object) {
  ctx, span := r.tracer.Start(ctx, "RoleService.GetSuper")
  defer span.End()

  return r.getByName(ctx, constant.SUPER_ROLE_NAME)
}

func (r *roleService) appendDefaultRolesPermissions(ctx context.Context, roleName string, permIds ...types.Id) status.Object {
  span := trace.SpanFromContext(ctx)

  // Find role by names
  role, err := r.roleRepo.FindByName(ctx, roleName)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepositoryOverride(err, types.NewPair(status.INTERNAL_SERVER_ERROR, errs.ErrDefaultRoleNotSeeded))
  }

  // Append permission into it
  err = r.roleRepo.AddPermissions(ctx, role.Id, permIds...)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepositoryExist(err)
  }

  return status.Created()
}

func (r *roleService) AppendSuperRolesPermission(ctx context.Context, permIds ...types.Id) status.Object {
  ctx, span := r.tracer.Start(ctx, "RoleService.AppendSuperRolesPermission")
  defer span.End()

  return r.appendDefaultRolesPermissions(ctx, constant.SUPER_ROLE_NAME, permIds...)
}

func (r *roleService) AppendDefaultRolesPermission(ctx context.Context, permIds ...types.Id) status.Object {
  ctx, span := r.tracer.Start(ctx, "RoleService.AppendDefaultRolesPermission")
  defer span.End()

  return r.appendDefaultRolesPermissions(ctx, constant.DEFAULT_ROLE_NAME, permIds...)
}

func (r *roleService) setUserRolesByName(ctx context.Context, userId types.Id, roleName string) status.Object {
  span := trace.SpanFromContext(ctx)

  // Get super roles
  role, err := r.roleRepo.FindByName(ctx, roleName)
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

func (r *roleService) SetUserAsSuper(ctx context.Context, userId types.Id) status.Object {
  ctx, span := r.tracer.Start(ctx, "RoleService.SetUserAsSuper")
  defer span.End()

  return r.setUserRolesByName(ctx, userId, constant.SUPER_ROLE_NAME)
}

func (r *roleService) SetUserAsDefault(ctx context.Context, userId types.Id) status.Object {
  ctx, span := r.tracer.Start(ctx, "RoleService.SetUserAsDefault")
  defer span.End()

  return r.setUserRolesByName(ctx, userId, constant.DEFAULT_ROLE_NAME)
}
