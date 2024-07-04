package service

import (
  "context"
  "nexa/services/authorization/internal/domain/dto"
  sharedDto "nexa/shared/dto"
  "nexa/shared/status"
  "nexa/shared/types"
)

type IPermission interface {
  Create(ctx context.Context, createDTO *dto.PermissionCreateDTO) (types.Id, status.Object)
  Find(ctx context.Context, permIds ...types.Id) ([]dto.PermissionResponseDTO, status.Object)
  FindByRoles(ctx context.Context, roleId ...types.Id) ([]dto.PermissionResponseDTO, status.Object)
  FindAll(ctx context.Context, input *sharedDto.PagedElementDTO) (sharedDto.PagedElementResult[dto.PermissionResponseDTO], status.Object)
  Delete(ctx context.Context, permId types.Id) status.Object
}
