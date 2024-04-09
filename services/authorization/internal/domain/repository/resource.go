package repository

import (
  "context"
  "nexa/services/authorization/internal/domain/entity"
  "nexa/shared/types"
  "nexa/shared/util/repo"
)

type IResource interface {
  FindById(ctx context.Context, id types.Id) (entity.Resource, error)
  FindAll(ctx context.Context, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Resource], error)
  Create(ctx context.Context, resource *entity.Resource) error
  Patch(ctx context.Context, resource *entity.Resource) error
  DeleteById(ctx context.Context, id types.Id) error
}
