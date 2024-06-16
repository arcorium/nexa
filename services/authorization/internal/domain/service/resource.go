package service

import (
  "context"
  "nexa/services/authorization/internal/domain/dto"
  sharedDto "nexa/shared/dto"
  "nexa/shared/status"
  "nexa/shared/types"
)

type IResource interface {
  Find(ctx context.Context, id types.Id) (dto.ResourceResponseDTO, status.Object)
  FindAll(ctx context.Context, input *sharedDto.PagedElementDTO) (sharedDto.PagedElementResult[dto.ResourceResponseDTO], status.Object)
  Create(ctx context.Context, createDTO *dto.ResourceCreateDTO) (types.Id, status.Object)
  Update(ctx context.Context, updateDTO *dto.ResourceUpdateDTO) status.Object
  Delete(ctx context.Context, id types.Id) status.Object
}
