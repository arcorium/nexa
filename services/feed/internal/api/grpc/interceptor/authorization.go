package interceptor

import (
  "context"
  feedv1 "github.com/arcorium/nexa/proto/gen/go/feed/v1"
  sharedJwt "github.com/arcorium/nexa/shared/jwt"
  "github.com/arcorium/nexa/shared/logger"
  authUtil "github.com/arcorium/nexa/shared/util/auth"
  "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
  "nexa/services/feed/constant"
)

func AuthSelector(_ context.Context, meta interceptors.CallMeta) bool {
  return true // All endpoint needs authorization
}

func PermissionCheck(claims *sharedJwt.UserClaims, meta interceptors.CallMeta) bool {
  switch meta.FullMethod() {
  case feedv1.FeedService_GetUserFeed_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.FEED_PERMISSIONS[constant.FEED_GET])
  default:
    logger.Warnf("Method %s undefined", meta.FullMethod())
  }
  return true
}
