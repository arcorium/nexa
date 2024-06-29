package jwt

import (
  "context"
  "errors"
  "github.com/golang-jwt/jwt/v5"
  sharedConst "nexa/shared/constant"
)

type Role struct {
  Id          string   `json:"id"`
  Role        string   `json:"role"`
  Permissions []string `json:"perms"`
}

type UserClaims struct {
  jwt.RegisteredClaims
  RefreshTokenId string `json:"rtid"`
  UserId         string `json:"uid"`
  Username       string `json:"username"`
  Roles          []Role `json:"roles"`
}

func GetClaimsFromCtx(ctx context.Context) (*UserClaims, error) {
  value := ctx.Value(sharedConst.CLAIMS_CONTEXT_KEY)
  val, ok := value.(*UserClaims)
  if !ok {
    return nil, errors.New("claims not found")
  }
  return val, nil
}
