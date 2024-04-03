package repository

import (
	"context"
	"nexa/services/authorization/shared/domain/entity"
	"nexa/shared/types"
	"nexa/shared/util/repo"
)

type IPermission interface {
	FindById(ctx context.Context, id types.Id) (entity.Permission, error)
	FindByUserId(ctx context.Context, userId types.Id) ([]entity.Permission, error)
	FindAll(ctx context.Context, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Permission], error)
	Create(ctx context.Context, permission *entity.Permission) error
	Delete(ctx context.Context, id types.Id) error
}
