package external

import (
	"nexa/services/authorization/shared/domain/entity"
	"nexa/shared/types"
)

type IAuthorizationClient interface {
	FindUserRoles(userId types.Id) ([]entity.Role, error)
	HasPermission(userId types.Id) error
}
