package repository

import (
  "context"
  "nexa/services/authentication/internal/domain/entity"
  "nexa/shared/types"
  "nexa/shared/util/repo"
)

type ITokenUsage interface {
  Create(ctx context.Context, usage *entity.TokenUsage) (types.Id, error)
  Patch(ctx context.Context, usage *entity.TokenUsage) error
  Delete(ctx context.Context, id types.Id) error
  Find(ctx context.Context, id types.Id) (entity.TokenUsage, error)
  FindAll(ctx context.Context, parameter repo.QueryParameter) (repo.PaginatedResult[entity.TokenUsage], error)
}
