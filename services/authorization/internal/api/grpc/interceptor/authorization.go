package interceptor

import (
  "context"
  authZv1 "github.com/arcorium/nexa/proto/gen/go/authorization/v1"
  "github.com/arcorium/nexa/shared/grpc/interceptor/authz"
  sharedJwt "github.com/arcorium/nexa/shared/jwt"
  "github.com/arcorium/nexa/shared/logger"
  authUtil "github.com/arcorium/nexa/shared/util/auth"
  "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
  "nexa/services/authorization/constant"
  "slices"
)

var privateApi = []string{
  authZv1.RoleService_AppendDefaultRolePermissions_FullMethodName,
  authZv1.RoleService_SetAsSuper_FullMethodName,
  authZv1.PermissionService_Seed_FullMethodName,
}

var publicApi = []string{
  authZv1.RoleService_GetUsers_FullMethodName,
  authZv1.RoleService_GetDefault_FullMethodName,
  authZv1.PermissionService_FindByRoles_FullMethodName,
  authZv1.PermissionService_FindAll_FullMethodName,
}

func CombinationSelector(_ context.Context, meta interceptors.CallMeta) authz.CombinationType {
  if slices.Contains(privateApi, meta.FullMethod()) {
    return authz.Private
  }
  if slices.Contains(publicApi, meta.FullMethod()) {
    return authz.Public
  }
  return authz.Authorized
}

func UserCheckPermission(claims *sharedJwt.UserClaims, meta interceptors.CallMeta) bool {
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
  case authZv1.PermissionService_Create_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.AUTHZ_PERMISSIONS[constant.AUTHZ_CREATE_PERMISSION])
  case authZv1.PermissionService_Delete_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.AUTHZ_PERMISSIONS[constant.AUTHZ_DELETE_PERMISSION])
  default:
    logger.Warnf("Unknown method: %s", meta.FullMethod())
  }

  return true
}
