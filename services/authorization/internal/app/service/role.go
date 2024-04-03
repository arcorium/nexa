package service

import (
	"context"
	"nexa/services/authorization/internal/domain/dto"
	"nexa/services/authorization/internal/domain/mapper"
	"nexa/services/authorization/internal/domain/repository"
	"nexa/services/authorization/internal/domain/service"
	sharedDto "nexa/shared/dto"
	"nexa/shared/status"
	"nexa/shared/types"
	"nexa/shared/util"
)

func NewRole(role repository.IRole) service.IRole {
	return &roleService{repo: role}
}

type roleService struct {
	repo repository.IRole
}

func (r *roleService) Find(ctx context.Context, ids ...types.Id) ([]dto.RoleResponseDTO, status.Object) {
	roles, err := r.repo.FindByIds(ctx, ids...)

	return util.CastSlice(roles, mapper.ToRoleResponseDTO), status.FromRepository(err, status.NullCode)
}

func (r *roleService) FindByUserId(ctx context.Context, userId types.Id) ([]dto.RoleResponseDTO, status.Object) {
	roles, err := r.repo.FindByUserId(ctx, userId)
	if err != nil {
		return nil, status.FromRepository(err, status.NullCode)
	}
	return util.CastSlice(roles, mapper.ToRoleResponseDTO), status.Success()
}

func (r *roleService) FindAll(ctx context.Context, input *sharedDto.PagedElementDTO) (sharedDto.PagedElementResult[dto.RoleResponseDTO], status.Object) {
	result, err := r.repo.FindAll(ctx, input.ToQueryParam())
	rolesDTO := util.CastSlice(result.Data, mapper.ToRoleResponseDTO)

	return sharedDto.NewPagedElementOutput2(rolesDTO, input, result.Total), status.FromRepository(err, status.NullCode)
}

func (r *roleService) Create(ctx context.Context, createDTO *dto.RoleCreateDTO) (types.Id, status.Object) {
	role := createDTO.ToDomain()
	err := r.repo.Create(ctx, &role)
	return role.Id, status.FromRepository(err, status.NullCode)
}

func (r *roleService) Update(ctx context.Context, updateDTO *dto.RoleUpdateDTO) status.Object {
	role := updateDTO.ToDomain()
	err := r.repo.Patch(ctx, &role)
	return status.FromRepository(err, status.NullCode)
}

func (r *roleService) Delete(ctx context.Context, id types.Id) status.Object {
	err := r.repo.Delete(ctx, id)
	return status.FromRepository(err, status.NullCode)
}

func (r *roleService) AddPermissions(ctx context.Context, permissionsDTO *dto.RoleAddPermissionsDTO) status.Object {
	roleId := types.IdFromString(permissionsDTO.RoleId)
	permissionIds := util.CastSlice(permissionsDTO.PermissionIds, func(from *string) types.Id {
		return types.IdFromString(*from)
	})
	err := r.repo.AddPermissions(ctx, roleId, permissionIds...)
	return status.FromRepository(err, status.NullCode)
}

func (r *roleService) RemovePermissions(ctx context.Context, permissionsDTO *dto.RoleRemovePermissionsDTO) status.Object {
	roleId := types.IdFromString(permissionsDTO.RoleId)
	permissionIds := util.CastSlice(permissionsDTO.PermissionIds, func(from *string) types.Id {
		return types.IdFromString(*from)
	})

	err := r.repo.RemovePermissions(ctx, roleId, permissionIds...)
	return status.FromRepository(err, status.NullCode)
}

func (r *roleService) AddUsers(ctx context.Context, usersDTO *dto.RoleAddUsersDTO) status.Object {
	userId := types.IdFromString(usersDTO.UserId)
	roleIds := util.CastSlice(usersDTO.RoleIds, func(from *string) types.Id {
		return types.IdFromString(*from)
	})

	err := r.repo.AddUser(ctx, userId, roleIds...)
	return status.FromRepository(err, status.NullCode)
}

func (r *roleService) RemoveUsers(ctx context.Context, usersDTO *dto.RoleRemoveUsersDTO) status.Object {
	userId := types.IdFromString(usersDTO.UserId)
	roleIds := util.CastSlice(usersDTO.RoleIds, func(from *string) types.Id {
		return types.IdFromString(*from)
	})

	err := r.repo.RemoveUser(ctx, userId, roleIds...)
	return status.FromRepository(err, status.NullCode)
}
