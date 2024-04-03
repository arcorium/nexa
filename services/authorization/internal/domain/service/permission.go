package service

import (
	"context"
	"nexa/services/authorization/internal/domain/dto"
	sharedDto "nexa/shared/dto"
	"nexa/shared/status"
	"nexa/shared/types"
)

type IPermission interface {
	CheckUserPermission(ctx context.Context, permissionDTO *dto.CheckUserPermissionDTO) status.Object
	Find(ctx context.Context, id types.Id) (dto.PermissionResponseDTO, status.Object)
	FindAll(ctx context.Context, input *sharedDto.PagedElementDTO) (sharedDto.PagedElementResult[dto.PermissionResponseDTO], status.Object)
	Create(ctx context.Context, createDTO *dto.PermissionCreateDTO) (types.Id, status.Object)
	Delete(ctx context.Context, id types.Id) status.Object
}
