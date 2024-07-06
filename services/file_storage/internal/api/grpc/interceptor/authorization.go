package interceptor

import (
  "context"
  "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
  sharedJwt "nexa/shared/jwt"
)

func AuthSelector(_ context.Context, _ interceptors.CallMeta) bool {
  return true // all need auth
}

func Auth(_ *sharedJwt.UserClaims, _ interceptors.CallMeta) bool {
  // Doesn't need permissions
  return true
}
