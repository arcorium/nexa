package service

import (
  "context"
  sharedDto "github.com/arcorium/nexa/shared/dto"
  "github.com/arcorium/nexa/shared/status"
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/reaction/internal/domain/dto"
  "nexa/services/reaction/internal/domain/entity"
)

type IReaction interface {
  Like(ctx context.Context, itemType entity.ItemType, itemId types.Id) status.Object
  Dislike(ctx context.Context, itemType entity.ItemType, itemId types.Id) status.Object
  GetItemsReactions(ctx context.Context, itemType entity.ItemType, itemId types.Id, pageDTO sharedDto.PagedElementDTO) (sharedDto.PagedElementResult[dto.ReactionResponseDTO], status.Object)
  GetCounts(ctx context.Context, itemType entity.ItemType, itemIds ...types.Id) ([]dto.CountResponseDTO, status.Object)
  Delete(ctx context.Context, itemType entity.ItemType, itemIds ...types.Id) status.Object
  ClearUsers(ctx context.Context, userId types.Id) status.Object
}
