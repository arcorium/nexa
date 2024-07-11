package repository

import (
  "context"
  "github.com/arcorium/nexa/shared/types"
  "github.com/arcorium/nexa/shared/util/repo"
  "nexa/services/authentication/internal/domain/entity"
)

type IToken interface {
  Create(ctx context.Context, token *entity.Token) error
  Delete(ctx context.Context, token string) error
  DeleteByUserId(ctx context.Context, userId types.Id) error
  Find(ctx context.Context, token string) (entity.Token, error)
  Get(ctx context.Context, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Token], error)
}
