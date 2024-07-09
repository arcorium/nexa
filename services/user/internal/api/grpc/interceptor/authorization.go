package interceptor

import (
  "context"
  userv1 "github.com/arcorium/nexa/proto/gen/go/user/v1"
  sharedJwt "github.com/arcorium/nexa/shared/jwt"
  "github.com/arcorium/nexa/shared/logger"
  authUtil "github.com/arcorium/nexa/shared/util/auth"
  "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
  "nexa/services/user/constant"
  "slices"
)

func AuthSkipSelector(_ context.Context, meta interceptors.CallMeta) bool {
  return slices.Contains([]string{
    userv1.UserService_Create_FullMethodName,
    userv1.UserService_Find_FullMethodName,
    userv1.UserService_FindByIds_FullMethodName,
    userv1.UserService_Validate_FullMethodName,
    userv1.UserService_ResetPassword_FullMethodName,
    userv1.UserService_ForgotPassword_FullMethodName,
    userv1.UserService_VerifyEmail_FullMethodName, // TODO: Need rework
    userv1.ProfileService_Find_FullMethodName,
  }, meta.FullMethod())
}

func PermissionCheck(claims *sharedJwt.UserClaims, meta interceptors.CallMeta) bool {
  switch meta.FullMethod() {
  case userv1.UserService_Update_FullMethodName:
    fallthrough
  case userv1.UserService_UpdatePassword_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.USER_PERMISSIONS[constant.USER_UPDATE])
  case userv1.UserService_VerifyEmail_FullMethodName:
    break
  case userv1.UserService_Find_FullMethodName:
    fallthrough
  case userv1.UserService_FindByIds_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.USER_PERMISSIONS[constant.USER_GET])
  case userv1.UserService_Banned_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.USER_PERMISSIONS[constant.USER_BANNED])
  case userv1.UserService_Delete_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.USER_PERMISSIONS[constant.USER_DELETE])
  case userv1.ProfileService_Find_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.USER_PERMISSIONS[constant.PROFILE_GET])
  default:
    logger.Warnf("Unknown permission method: %s", meta.FullMethod())
  }

  return true
}
