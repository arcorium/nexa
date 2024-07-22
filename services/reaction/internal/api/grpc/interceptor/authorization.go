package interceptor

import (
  "context"
  reactionv1 "github.com/arcorium/nexa/proto/gen/go/reaction/v1"
  sharedJwt "github.com/arcorium/nexa/shared/jwt"
  "github.com/arcorium/nexa/shared/logger"
  authUtil "github.com/arcorium/nexa/shared/util/auth"
  "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
  "nexa/services/reaction/constant"
)

func AuthSelector(_ context.Context, meta interceptors.CallMeta) bool {
  return true
}

func PermissionCheck(claims *sharedJwt.UserClaims, meta interceptors.CallMeta) bool {
  switch meta.FullMethod() {
  case reactionv1.ReactionService_Like_FullMethodName:
    fallthrough
  case reactionv1.ReactionService_Dislike_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.REACTION_PERMISSIONS[constant.REACTION_CREATE])
  case reactionv1.ReactionService_GetItems_FullMethodName:
    fallthrough
  case reactionv1.ReactionService_GetCount_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.REACTION_PERMISSIONS[constant.REACTION_GET])
  case reactionv1.ReactionService_DeleteItems_FullMethodName:
    fallthrough
  case reactionv1.ReactionService_ClearUsers_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.REACTION_PERMISSIONS[constant.REACTION_DELETE])
  default:
    logger.Warnf("Unknown method: %s", meta.FullMethod())
  }

  return true
}
