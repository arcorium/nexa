package service

import (
  "context"
  "nexa/services/authorization/internal/domain/dto"
  sharedDto "nexa/shared/dto"
  "nexa/shared/status"
  "nexa/shared/types"
)

type IRole interface {
  FindByIds(ctx context.Context, ids ...types.Id) ([]dto.RoleResponseDTO, status.Object)
  FindByUserId(ctx context.Context, userId types.Id) ([]dto.RoleResponseDTO, status.Object)
  FindAll(ctx context.Context, input *sharedDto.PagedElementDTO) (sharedDto.PagedElementResult[dto.RoleResponseDTO], status.Object)
  Create(ctx context.Context, createDTO *dto.RoleCreateDTO) (types.Id, status.Object)
  Update(ctx context.Context, updateDTO *dto.RoleUpdateDTO) status.Object
  Delete(ctx context.Context, roleId types.Id) status.Object
  // AddPermissions append permissions into role
  AddPermissions(ctx context.Context, modifyDTO *dto.ModifyRolesPermissionsDTO) status.Object
  // RemovePermissions remove permissions from role
  RemovePermissions(ctx context.Context, modifyDTO *dto.ModifyRolesPermissionsDTO) status.Object
  // AddUsers add roles into user
  AddUsers(ctx context.Context, usersDTO *dto.ModifyUserRolesDTO) status.Object
  // RemoveUsers remove roles from user
  RemoveUsers(ctx context.Context, usersDTO *dto.ModifyUserRolesDTO) status.Object

  AppendSuperRolesPermission(ctx context.Context, permIds ...types.Id) status.Object
}
