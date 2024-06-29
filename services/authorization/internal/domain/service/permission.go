package service

import (
  "context"
  "nexa/services/authorization/internal/domain/dto"
  sharedDto "nexa/shared/dto"
  "nexa/shared/status"
)

type IPermission interface {
  //CheckUserPermission(ctx context.Context, permissionDTO *dto.CheckUserPermissionDTO) status.Object
  Create(ctx context.Context, createDTO *dto.PermissionCreateDTO) (string, status.Object)
  Find(ctx context.Context, permIds ...string) ([]dto.PermissionResponseDTO, status.Object)
  FindByRoles(ctx context.Context, roleId ...string) ([]dto.PermissionResponseDTO, status.Object)
  FindAll(ctx context.Context, input *sharedDto.PagedElementDTO) (sharedDto.PagedElementResult[dto.PermissionResponseDTO], status.Object)
  Delete(ctx context.Context, id string) status.Object
}
