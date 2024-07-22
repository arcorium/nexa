package repository

import (
  "context"
  "github.com/arcorium/nexa/shared/optional"
  "github.com/arcorium/nexa/shared/types"
  "github.com/arcorium/nexa/shared/util/repo"
  "nexa/services/reaction/internal/domain/entity"
)

type IReaction interface {
  Delsert(ctx context.Context, reaction *entity.Reaction) error
  FindByItemId(ctx context.Context, itemType entity.ItemType, itemId types.Id, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Reaction], error)
  FindByUserId(ctx context.Context, userId types.Id, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Reaction], error)
  GetCounts(ctx context.Context, itemType entity.ItemType, itemIds ...types.Id) ([]entity.Count, error)
  DeleteByItemId(ctx context.Context, itemType entity.ItemType, itemIds ...types.Id) error
  DeleteByUserId(ctx context.Context, userId types.Id, items optional.Object[types.Pair[entity.ItemType, []types.Id]]) error
}
