package util

import (
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/authorization/constant"
)

// GetPermission get encoded permission from action
func GetPermission(action types.Action) string {
  return constant.AUTHZ_PERMISSIONS[action]
}
