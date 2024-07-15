package interceptor

import (
  "context"
  authNv1 "github.com/arcorium/nexa/proto/gen/go/authentication/v1"
  sharedJwt "github.com/arcorium/nexa/shared/jwt"
  "github.com/arcorium/nexa/shared/logger"
  authUtil "github.com/arcorium/nexa/shared/util/auth"
  "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
  "nexa/services/authentication/constant"
  "slices"
)

var publicApi = []string{
  authNv1.AuthenticationService_Login_FullMethodName,
  authNv1.AuthenticationService_Register_FullMethodName,
  authNv1.AuthenticationService_RefreshToken_FullMethodName,
  authNv1.UserService_ForgotPassword_FullMethodName,
  authNv1.UserService_VerifyEmail_FullMethodName, //
  authNv1.UserService_ResetPasswordByToken_FullMethodName,
}

func AuthSkipSelector(_ context.Context, meta interceptors.CallMeta) bool {
  return slices.Contains(publicApi, meta.FullMethod())
}

func PermissionCheck(claims *sharedJwt.UserClaims, meta interceptors.CallMeta) bool {
  switch meta.FullMethod() {
  // Authentication
  case authNv1.AuthenticationService_GetCredentials_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.AUTHN_PERMISSIONS[constant.AUTHN_GET_CREDENTIAL])
  case authNv1.AuthenticationService_Logout_FullMethodName:
    fallthrough
  case authNv1.AuthenticationService_LogoutAll_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.AUTHN_PERMISSIONS[constant.AUTHN_LOGOUT_USER])
  // User
  case authNv1.UserService_Create_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.AUTHN_PERMISSIONS[constant.AUTHN_CREATE_USER])
  case authNv1.UserService_Update_FullMethodName:
    fallthrough
  case authNv1.UserService_UpdatePassword_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.AUTHN_PERMISSIONS[constant.AUTHN_UPDATE_USER])
  case authNv1.UserService_Find_FullMethodName:
    fallthrough
  case authNv1.UserService_FindByIds_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.AUTHN_PERMISSIONS[constant.AUTHN_GET_USER])
  case authNv1.UserService_Banned_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.AUTHN_PERMISSIONS[constant.AUTHN_BANNED])
  case authNv1.UserService_Delete_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.AUTHN_PERMISSIONS[constant.AUTHN_DELETE_USER])
  case authNv1.UserService_ResetPassword_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.AUTHN_PERMISSIONS[constant.AUTHN_UPDATE_USER_ARB])
  // Profile
  case authNv1.ProfileService_Update_FullMethodName:
    fallthrough
  case authNv1.ProfileService_UpdateAvatar_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.AUTHN_PERMISSIONS[constant.AUTHN_UPDATE_USER])
  default:
    logger.Warnf("Unknown permission method: %s", meta.FullMethod())
  }

  return true
}
