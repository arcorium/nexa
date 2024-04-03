package repository

import (
	"context"
	"nexa/services/authorization/shared/domain/entity"
	"nexa/shared/types"
	"nexa/shared/util/repo"
)

type IRole interface {
	FindByIds(ctx context.Context, id ...types.Id) ([]entity.Role, error)
	FindByUserId(ctx context.Context, userId types.Id) ([]entity.Role, error)
	FindAll(ctx context.Context, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Role], error)
	Create(ctx context.Context, role *entity.Role) error
	Patch(ctx context.Context, role *entity.Role) error
	Delete(ctx context.Context, id types.Id) error
	// AddPermissions add permission into role
	AddPermissions(ctx context.Context, roleId types.Id, permissionIds ...types.Id) error
	// RemovePermissions remove permission from role
	RemovePermissions(ctx context.Context, roleId types.Id, permissionIds ...types.Id) error
	// AddUser add roles into user
	AddUser(ctx context.Context, userId types.Id, roleIds ...types.Id) error
	// RemoveUser remove roles from user
	RemoveUser(ctx context.Context, userId types.Id, roleIds ...types.Id) error
}
