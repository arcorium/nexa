package interceptor

import (
  "context"
  commentv1 "github.com/arcorium/nexa/proto/gen/go/comment/v1"
  sharedJwt "github.com/arcorium/nexa/shared/jwt"
  "github.com/arcorium/nexa/shared/logger"
  authUtil "github.com/arcorium/nexa/shared/util/auth"
  "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
  "nexa/services/comment/constant"
)

func AuthSelector(_ context.Context, meta interceptors.CallMeta) bool {
  return meta.FullMethod() != commentv1.CommentService_IsExist_FullMethodName
}

func PermissionCheck(claims *sharedJwt.UserClaims, meta interceptors.CallMeta) bool {
  switch meta.FullMethod() {
  case commentv1.CommentService_Create_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.COMMENT_PERMISSIONS[constant.COMMENT_CREATE])
  case commentv1.CommentService_Edit_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.COMMENT_PERMISSIONS[constant.COMMENT_UPDATE])
  case commentv1.CommentService_Delete_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.COMMENT_PERMISSIONS[constant.COMMENT_DELETE])
  case commentv1.CommentService_GetPosts_FullMethodName:
    fallthrough
  case commentv1.CommentService_GetReplies_FullMethodName:
    fallthrough
  case commentv1.CommentService_GetCounts_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.COMMENT_PERMISSIONS[constant.COMMENT_GET])
  case commentv1.CommentService_ClearPosts_FullMethodName:
    fallthrough
  case commentv1.CommentService_ClearUsers_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.COMMENT_PERMISSIONS[constant.COMMENT_DELETE])
  default:
    logger.Warnf("Unknown permission method: %s", meta.FullMethod())
  }

  return true
}
