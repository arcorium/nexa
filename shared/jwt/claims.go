package jwt

import (
  "context"
  "errors"
  sharedConst "github.com/arcorium/nexa/shared/constant"
  "github.com/golang-jwt/jwt/v5"
)

var DefaultSigningMethod = jwt.SigningMethodRS256

type Role struct {
  Id          string   `json:"id"`
  Role        string   `json:"role"`
  Permissions []string `json:"perms"`
}

type UserClaims struct {
  jwt.RegisteredClaims
  CredentialId string `json:"cid"`
  UserId       string `json:"uid"`
  Username     string `json:"username"`
  Roles        []Role `json:"roles"`
}

func GetUserClaimsFromCtx(ctx context.Context) (*UserClaims, error) {
  value := ctx.Value(sharedConst.USER_CLAIMS_CONTEXT_KEY)
  val, ok := value.(*UserClaims)
  if !ok {
    return nil, errors.New("claims not found")
  }
  return val, nil
}

// PrivateClaims Claims used for temporary token. It is used to call services protected API
type PrivateClaims struct {
  jwt.RegisteredClaims
}

func GetTempClaimsFromCtx(ctx context.Context) (*UserClaims, error) {
  value := ctx.Value(sharedConst.USER_CLAIMS_CONTEXT_KEY)
  val, ok := value.(*UserClaims)
  if !ok {
    return nil, errors.New("claims not found")
  }
  return val, nil
}
