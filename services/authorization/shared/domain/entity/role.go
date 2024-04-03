package entity

import (
	"nexa/shared/types"
)

type Role struct {
	Id          types.Id
	Name        string
	Description string

	Permissions []Permission
}

func (r *Role) HasPermissions(permissions ...Permission) bool {
	// TODO: Implement
	return false
}
