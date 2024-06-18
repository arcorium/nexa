package service

import (
  "context"
  "nexa/services/mailer/internal/domain/dto"
  sharedDto "nexa/shared/dto"
  "nexa/shared/status"
  "nexa/shared/types"
)

type ITag interface {
  Find(ctx context.Context, dto *sharedDto.PagedElementDTO) (*sharedDto.PagedElementResult[dto.TagResponseDTO], status.Object)
  FindByIds(ctx context.Context, ids ...types.Id) ([]dto.TagResponseDTO, status.Object)
  FindByName(ctx context.Context, name string) (dto.TagResponseDTO, status.Object)
  Create(ctx context.Context, dto *dto.CreateTagDTO) (types.Id, status.Object)
  Update(ctx context.Context, dto *dto.UpdateTagDTO) status.Object
  Remove(ctx context.Context, id types.Id) status.Object
}
