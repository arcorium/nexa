package auth

import (
  "fmt"
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
    return false
  })
}

func Encode(resource, action string) string {
  return fmt.Sprintf("%s:%s", resource, action)
}

func FullEncode(resource string, actions ...string) map[string]string {
  result := make(map[string]string, len(actions))
  for _, action := range actions {
    result[action] = Encode(resource, action)
  }
  return result
}
