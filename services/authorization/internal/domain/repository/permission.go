package repository

import (
  "context"
  "github.com/arcorium/nexa/shared/types"
  "github.com/arcorium/nexa/shared/util/repo"
  "nexa/services/authorization/internal/domain/entity"
)

type IPermission interface {
  Create(ctx context.Context, permission *entity.Permission) error
  Creates(ctx context.Context, permissions ...entity.Permission) error
  FindByIds(ctx context.Context, ids ...types.Id) ([]entity.Permission, error)
  FindByRoleIds(ctx context.Context, roleIds ...types.Id) ([]entity.Permission, error)
  Get(ctx context.Context, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Permission], error)
  Delete(ctx context.Context, id types.Id) error
}
