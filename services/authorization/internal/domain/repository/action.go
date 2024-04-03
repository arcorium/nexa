package repository

import (
	"context"
	"nexa/services/authorization/shared/domain/entity"
	"nexa/shared/types"
	"nexa/shared/util/repo"
)

type IAction interface {
	FindById(ctx context.Context, id types.Id) (entity.Action, error)
	FindAll(ctx context.Context, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Action], error)
	Create(ctx context.Context, action *entity.Action) error
	Patch(ctx context.Context, action *entity.Action) error
	DeleteById(ctx context.Context, id types.Id) error
}
