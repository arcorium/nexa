package interceptor

import (
  "context"
  "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
  userv1 "nexa/proto/gen/go/user/v1"
  "nexa/services/user/constant"
  sharedJwt "nexa/shared/jwt"
  "nexa/shared/logger"
  authUtil "nexa/shared/util/auth"
  "slices"
)

func AuthSelector(_ context.Context, meta interceptors.CallMeta) bool {
  return !slices.Contains([]string{
    userv1.UserService_Create_FullMethodName,
    userv1.UserService_Validate_FullMethodName,
    userv1.UserService_ResetPassword_FullMethodName,
    userv1.UserService_ForgotPassword_FullMethodName,
    userv1.UserService_VerifyEmail_FullMethodName,
    userv1.ProfileService_Update_FullMethodName,
    userv1.ProfileService_UpdateAvatar_FullMethodName,
  }, meta.FullMethod())
}

func PermissionCheck(claims *sharedJwt.UserClaims, meta interceptors.CallMeta) bool {
  switch meta.FullMethod() {
  case userv1.UserService_Update_FullMethodName:
    fallthrough
  case userv1.UserService_VerifyEmail_FullMethodName:
    fallthrough
  case userv1.UserService_UpdatePassword_FullMethodName:
    return authUtil.ContainsPermissions(claims.Roles, constant.USER_UPDATE, constant.USER_UPDATE_OTHER)
  case userv1.UserService_Find_FullMethodName:
    fallthrough
  case userv1.UserService_FindByIds_FullMethodName:
    return authUtil.ContainsPermissions(claims.Roles, constant.USER_GET)
  case userv1.UserService_Banned_FullMethodName:
    return authUtil.ContainsPermissions(claims.Roles, constant.USER_BANNED)
  case userv1.UserService_Delete_FullMethodName:
    return authUtil.ContainsPermissions(claims.Roles, constant.USER_DELETE, constant.USER_DELETE_OTHER)
  case userv1.ProfileService_Find_FullMethodName:
    return authUtil.ContainsPermissions(claims.Roles, constant.USER_GET_PROFILE, constant.USER_GET_PROFILE_OTHER)
  }

  logger.Warn("Unknown method")
  return true
}
