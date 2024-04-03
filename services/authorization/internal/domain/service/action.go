package service

import (
	"context"
	"nexa/services/authorization/internal/domain/dto"
	sharedDto "nexa/shared/dto"
	"nexa/shared/status"
	"nexa/shared/types"
)

type IAction interface {
	Find(ctx context.Context, id types.Id) (dto.ActionResponseDTO, status.Object)
	FindAll(ctx context.Context, input *sharedDto.PagedElementDTO) (sharedDto.PagedElementResult[dto.ActionResponseDTO], status.Object)
	Create(ctx context.Context, input *dto.ActionCreateDTO) (types.Id, status.Object)
	Update(ctx context.Context, input *dto.ActionUpdateDTO) status.Object
	Delete(ctx context.Context, id types.Id) status.Object
}
