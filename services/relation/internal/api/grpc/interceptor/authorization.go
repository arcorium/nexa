package interceptor

import (
  "context"
  relationv1 "github.com/arcorium/nexa/proto/gen/go/relation/v1"
  sharedJwt "github.com/arcorium/nexa/shared/jwt"
  "github.com/arcorium/nexa/shared/logger"
  authUtil "github.com/arcorium/nexa/shared/util/auth"
  "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
  "nexa/services/relation/constant"
)

func AuthSelector(_ context.Context, meta interceptors.CallMeta) bool {
  return true
}

func PermissionCheck(claims *sharedJwt.UserClaims, meta interceptors.CallMeta) bool {
  switch meta.FullMethod() {
  // Block
  case relationv1.BlockService_Block_FullMethodName:
    fallthrough
  case relationv1.BlockService_Unblock_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.RELATION_PERMISSIONS[constant.RELATION_CREATE_BLOCK])
  case relationv1.BlockService_GetBlocked_FullMethodName:
    fallthrough
  case relationv1.BlockService_GetUsersCount_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.RELATION_PERMISSIONS[constant.RELATION_GET_BLOCK])
  case relationv1.BlockService_ClearUsers_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.RELATION_PERMISSIONS[constant.RELATION_DELETE_BLOCK])
    // Follow
  case relationv1.FollowService_Follow_FullMethodName:
    fallthrough
  case relationv1.FollowService_Unfollow_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.RELATION_PERMISSIONS[constant.RELATION_CREATE_FOLLOW])
  case relationv1.FollowService_GetFollowers_FullMethodName:
    fallthrough
  case relationv1.FollowService_GetFollowees_FullMethodName:
    fallthrough
  case relationv1.FollowService_GetRelation_FullMethodName:
    fallthrough
  case relationv1.FollowService_GetUsersCount_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.RELATION_PERMISSIONS[constant.RELATION_GET_FOLLOW])
  case relationv1.FollowService_ClearUsers_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.RELATION_PERMISSIONS[constant.RELATION_DELETE_FOLLOW])
  default:
    logger.Warnf("Unknown method: %s", meta.FullMethod())
  }

  return true
}
