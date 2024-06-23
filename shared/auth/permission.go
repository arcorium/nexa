package auth

import (
  "slices"
  "strings"
)

func ContainsPermissions(permissions []string, expected string) bool {
  return slices.ContainsFunc(permissions, func(s string) bool {
    return strings.EqualFold(s, expected)
  })
}
