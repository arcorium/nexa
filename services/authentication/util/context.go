package util

import (
  "context"
  "nexa/services/authentication/constant"
  "nexa/services/authentication/shared/common"
)

// GetUserClaims get user common.AccessTokenClaims from context and will panic when the context doesn't have the value
func GetUserClaims(ctx context.Context) *common.AccessTokenClaims {
  claims, ok := ctx.Value(constant.USER_CLAIMS_KEY).(*common.AccessTokenClaims)
  // TODO: Remove it to raise performance
  if !ok {
    panic("something wrong! Unauthenticated user should not access here")
  }
  return claims
}
