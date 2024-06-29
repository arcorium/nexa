package service

import (
  "context"
  "nexa/services/authorization/internal/domain/dto"
  sharedDto "nexa/shared/dto"
  "nexa/shared/status"
)

type IRole interface {
  Find(ctx context.Context, ids ...string) ([]dto.RoleResponseDTO, status.Object)
  FindByUserId(ctx context.Context, userId string) ([]dto.RoleResponseDTO, status.Object)
  FindAll(ctx context.Context, input *sharedDto.PagedElementDTO) (sharedDto.PagedElementResult[dto.RoleResponseDTO], status.Object)
  Create(ctx context.Context, createDTO *dto.RoleCreateDTO) (string, status.Object)
  Update(ctx context.Context, updateDTO *dto.RoleUpdateDTO) status.Object
  Delete(ctx context.Context, id string) status.Object
  AddPermissions(ctx context.Context, permissionDTO *dto.RoleAddPermissionsDTO) status.Object
  RemovePermissions(ctx context.Context, permissionsDTO *dto.RoleRemovePermissionsDTO) status.Object
  AddUsers(ctx context.Context, usersDTO *dto.RoleAddUsersDTO) status.Object
  RemoveUsers(ctx context.Context, usersDTO *dto.RoleRemoveUsersDTO) status.Object
}
