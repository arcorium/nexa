package repository

import (
  "context"
  "github.com/arcorium/nexa/shared/types"
  "github.com/arcorium/nexa/shared/util/repo"
  domain "nexa/services/mailer/internal/domain/entity"
)

type ITag interface {
  Get(ctx context.Context, query repo.QueryParameter) (repo.PaginatedResult[domain.Tag], error)
  FindByIds(ctx context.Context, ids ...types.Id) ([]domain.Tag, error)
  FindByName(ctx context.Context, name string) (*domain.Tag, error)
  Create(ctx context.Context, tag *domain.Tag) error
  Patch(ctx context.Context, tag *domain.PatchedTag) error
  Remove(ctx context.Context, id types.Id) error
}
