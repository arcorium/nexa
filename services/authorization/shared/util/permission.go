package util

import (
	"nexa/services/authorization/shared/domain/entity"
	"slices"
)

func HasPermission(haystack []entity.Permission, needle []entity.Permission) bool {
	return slices.ContainsFunc(haystack, func(permission entity.Permission) bool {
		return slices.ContainsFunc(needle, func(perm entity.Permission) bool {
			return perm.String() == permission.String()
		})
	})
}
