package jwt

import (
  "context"
  "errors"
  "github.com/golang-jwt/jwt/v5"
  sharedConst "nexa/shared/constant"
)

type UserClaims struct {
  jwt.RegisteredClaims
  UserId      string   `json:"user_id"`
  Name        string   `json:"name"`
  Roles       []string `json:"roles"`
  Permissions []string `json:"perms"`
}

func GetClaimsFromCtx(ctx context.Context) (*UserClaims, error) {
  value := ctx.Value(sharedConst.CLAIMS_CONTEXT_KEY)
  val, ok := value.(*UserClaims)
  if !ok {
    return nil, errors.New("claims not found")
  }
  return val, nil
}
