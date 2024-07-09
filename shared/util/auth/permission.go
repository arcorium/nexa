package auth

import (
  "fmt"
  "nexa/shared/jwt"
  "slices"
  "strings"
)

func ContainsPermission(roles []jwt.Role, expected string) bool {
  return slices.ContainsFunc(roles, func(role jwt.Role) bool {
    return HasPermission(role.Permissions, expected)
  })
}

func ContainsOneOfPermission(roles []jwt.Role, expected ...string) bool {
  for _, exp := range expected {
    for _, role := range roles {
      if HasPermission(role.Permissions, exp) {
        return true
      }
    }
  }
  return false
}

// HasPermission Check if it contains single expected permissions (code)
func HasPermission(permissions []string, expected string) bool {
  return slices.ContainsFunc(permissions, func(permission string) bool {
    return strings.EqualFold(permission, expected)
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
