package auth

import (
  "nexa/shared/jwt"
  "slices"
  "strings"
)

func ContainsPermissions(roles []jwt.Role, expected ...string) bool {
  return slices.ContainsFunc(roles, func(role jwt.Role) bool {
    return HasPermissions(role.Permissions, expected...)
  })
}

func HasPermissions(permissions []string, expected ...string) bool {
  return slices.ContainsFunc(permissions, func(permission string) bool {
    for _, perm := range expected {
      if strings.EqualFold(permission, perm) {
        return true
      }
    }
  })
}
