package service

import (
  "context"
  sharedDto "github.com/arcorium/nexa/shared/dto"
  "github.com/arcorium/nexa/shared/status"
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/mailer/internal/domain/dto"
)

type ITag interface {
  GetAll(ctx context.Context, dto *sharedDto.PagedElementDTO) (sharedDto.PagedElementResult[dto.TagResponseDTO], status.Object)
  FindByIds(ctx context.Context, ids ...types.Id) ([]dto.TagResponseDTO, status.Object)
  FindByName(ctx context.Context, name string) (dto.TagResponseDTO, status.Object)
  Create(ctx context.Context, dto *dto.CreateTagDTO) (types.Id, status.Object)
  Update(ctx context.Context, dto *dto.UpdateTagDTO) status.Object
  Remove(ctx context.Context, id types.Id) status.Object
}
