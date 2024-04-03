package external

import (
	"context"
	"nexa/services/authorization/shared/domain/entity"
	"nexa/shared/types"
)

type IAuthorizationClient interface {
	// FindUserRoles get user roles
	FindUserRoles(ctx context.Context, userId types.Id) ([]entity.Role, error)
	// HasPermission check is the user has permission to access the resource with those actions
	HasPermission(ctx context.Context, userId types.Id, resource string, actions ...string) error
}
