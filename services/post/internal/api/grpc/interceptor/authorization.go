package interceptor

import (
  "context"
  postv1 "github.com/arcorium/nexa/proto/gen/go/post/v1"
  sharedJwt "github.com/arcorium/nexa/shared/jwt"
  "github.com/arcorium/nexa/shared/logger"
  authUtil "github.com/arcorium/nexa/shared/util/auth"
  "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
  "nexa/services/post/constant"
)

func AuthSelector(_ context.Context, meta interceptors.CallMeta) bool {
  return true // All endpoint needs authorization
}

func PermissionCheck(claims *sharedJwt.UserClaims, meta interceptors.CallMeta) bool {
  switch meta.FullMethod() {
  case postv1.PostService_Find_FullMethodName:
    fallthrough
  case postv1.PostService_FindEdited_FullMethodName:
    fallthrough
  case postv1.PostService_FindById_FullMethodName:
    fallthrough
  case postv1.PostService_FindUsers_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.POST_PERMISSIONS[constant.POST_GET])
  case postv1.PostService_Create_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.POST_PERMISSIONS[constant.POST_CREATE])
  case postv1.PostService_UpdateVisibility_FullMethodName:
    fallthrough
  case postv1.PostService_Edit_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.POST_PERMISSIONS[constant.POST_UPDATE])
  case postv1.PostService_Bookmark_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.POST_PERMISSIONS[constant.POST_CREATE])
  case postv1.PostService_GetBookmarked_FullMethodName:
    fallthrough
  case postv1.PostService_Delete_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.POST_PERMISSIONS[constant.POST_DELETE])
  default:
    logger.Warnf("Method %s undefined", meta.FullMethod())
  }
  return true
}
