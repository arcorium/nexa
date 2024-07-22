package repository

import (
  "context"
  "github.com/arcorium/nexa/shared/types"
  "github.com/arcorium/nexa/shared/util/repo"
  "nexa/services/relation/internal/domain/entity"
)

type IBlock interface {
  Delsert(ctx context.Context, block *entity.Block) error
  Delete(ctx context.Context, block *entity.Block) error
  Create(ctx context.Context, block *entity.Block) error
  GetBlocked(ctx context.Context, userId types.Id, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Block], error)
  GetCounts(ctx context.Context, userIds ...types.Id) ([]entity.BlockCount, error)
  IsBlocked(ctx context.Context, blockerId, targetId types.Id) (bool, error)
  DeleteByUserId(ctx context.Context, deleteBlocker bool, userId types.Id) error
}
