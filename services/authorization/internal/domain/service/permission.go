package service

import (
  "context"
  sharedDto "github.com/arcorium/nexa/shared/dto"
  "github.com/arcorium/nexa/shared/status"
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/authorization/internal/domain/dto"
)

type IPermission interface {
  Create(ctx context.Context, createDTO *dto.PermissionCreateDTO) (types.Id, status.Object)
  Find(ctx context.Context, permIds ...types.Id) ([]dto.PermissionResponseDTO, status.Object)
  FindByRoles(ctx context.Context, roleId ...types.Id) ([]dto.PermissionResponseDTO, status.Object)
  GetAll(ctx context.Context, input *sharedDto.PagedElementDTO) (sharedDto.PagedElementResult[dto.PermissionResponseDTO], status.Object)
  Delete(ctx context.Context, permId types.Id) status.Object
}
