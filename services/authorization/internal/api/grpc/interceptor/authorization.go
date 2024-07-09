package interceptor

import (
  "context"
  "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
  authZv1 "nexa/proto/gen/go/authorization/v1"
  "nexa/services/authorization/constant"
  sharedJwt "nexa/shared/jwt"
  "nexa/shared/logger"
  authUtil "nexa/shared/util/auth"
  "strings"
)

func AuthSkipSelector(_ context.Context, callMeta interceptors.CallMeta) bool {
  return callMeta.Service == authZv1.PermissionService_ServiceDesc.ServiceName ||
      strings.EqualFold(callMeta.FullMethod(), authZv1.RoleService_GetUsers_FullMethodName) ||
      strings.EqualFold(callMeta.FullMethod(), authZv1.RoleService_AppendSuperRolePermissions_FullMethodName)
}

func Auth(claims *sharedJwt.UserClaims, meta interceptors.CallMeta) bool {
  switch meta.FullMethod() {
  case authZv1.RoleService_Create_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.AUTHZ_PERMISSIONS[constant.AUTHZ_CREATE_ROLE])
  case authZv1.RoleService_Update_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.AUTHZ_PERMISSIONS[constant.AUTHZ_UPDATE_ROLE])
  case authZv1.RoleService_Delete_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.AUTHZ_PERMISSIONS[constant.AUTHZ_DELETE_ROLE])
  case authZv1.RoleService_Find_FullMethodName:
    fallthrough
  case authZv1.RoleService_FindAll_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.AUTHZ_PERMISSIONS[constant.AUTHZ_READ_ROLE])
  case authZv1.RoleService_AddUser_FullMethodName:
    fallthrough
  case authZv1.RoleService_RemoveUser_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.AUTHZ_PERMISSIONS[constant.AUTHZ_MODIFY_USER_ROLE])
  case authZv1.RoleService_AppendPermissions_FullMethodName:
    fallthrough
  case authZv1.RoleService_RemovePermissions_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.AUTHZ_PERMISSIONS[constant.AUTHZ_MODIFY_USER_ROLE])
  }

  logger.Warn("Unknown method")
  return true
}
