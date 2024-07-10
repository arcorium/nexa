package interceptor

import (
  "context"
  "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
  authNv1 "nexa/proto/gen/go/authentication/v1"
  sharedJwt "nexa/shared/jwt"
  "strings"
)

func AuthSkipSelector(_ context.Context, callMeta interceptors.CallMeta) bool {
  return callMeta.Service == authNv1.TokenService_ServiceDesc.ServiceName ||
      strings.EqualFold(callMeta.FullMethod(), authNv1.CredentialService_Login_FullMethodName) ||
      strings.EqualFold(callMeta.FullMethod(), authNv1.CredentialService_Register_FullMethodName) ||
      strings.EqualFold(callMeta.FullMethod(), authNv1.CredentialService_RefreshToken_FullMethodName)
}

func PermissionCheck(_ *sharedJwt.UserClaims, _ interceptors.CallMeta) bool {
  // Permission check is happen on app layer
  return true
}