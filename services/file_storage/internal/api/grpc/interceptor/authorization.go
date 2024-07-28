package interceptor

import (
  "context"
  storagev1 "github.com/arcorium/nexa/proto/gen/go/file_storage/v1"
  sharedJwt "github.com/arcorium/nexa/shared/jwt"
  "github.com/arcorium/nexa/shared/logger"
  authUtil "github.com/arcorium/nexa/shared/util/auth"
  "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
  "nexa/services/file_storage/constant"
)

func AuthSelector(_ context.Context, meta interceptors.CallMeta) bool {
  return true
}

func PermissionCheck(claims *sharedJwt.UserClaims, meta interceptors.CallMeta) bool {
  switch meta.FullMethod() {
  case storagev1.FileStorageService_Find_FullMethodName:
    //return authUtil.ContainsPermission(claims.Roles, constant.FILE_STORAGE_PERMISSIONS[constant.FILE_GET])
    return true // Allow to get public file
  case storagev1.FileStorageService_FindMetadata_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.FILE_STORAGE_PERMISSIONS[constant.FILE_GET_METADATA])
  case storagev1.FileStorageService_Store_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.FILE_STORAGE_PERMISSIONS[constant.FILE_STORE])
  case storagev1.FileStorageService_Update_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.FILE_STORAGE_PERMISSIONS[constant.FILE_UPDATE])
  case storagev1.FileStorageService_Delete_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.FILE_STORAGE_PERMISSIONS[constant.FILE_DELETE])
  default:
    logger.Warnf("Unknown permission method: %s", meta.FullMethod())
  }

  return true
}
