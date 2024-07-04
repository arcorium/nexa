package interceptor

import (
  "context"
  "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
  sharedJwt "nexa/shared/jwt"
)

func AuthSelector(_ context.Context, callMeta interceptors.CallMeta) bool {
  return true // all need auth
}

func Auth(claims *sharedJwt.UserClaims, fullMethod string) bool {
  // Doesn't need permissions
  return true
}
