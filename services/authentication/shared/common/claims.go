package common

import (
  "github.com/golang-jwt/jwt/v5"
  "nexa/shared/types"
)

type AccessTokenClaims struct {
  UserId   types.Id
  Username string

  jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
  jwt.RegisteredClaims
}
