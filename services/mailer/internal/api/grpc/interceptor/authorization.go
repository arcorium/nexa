package interceptor

import (
  "context"
  "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
  mailerv1 "nexa/proto/gen/go/mailer/v1"
  "nexa/services/mailer/constant"
  sharedJwt "nexa/shared/jwt"
  "nexa/shared/logger"
  authUtil "nexa/shared/util/auth"
)

func AuthSelector(_ context.Context, callMeta interceptors.CallMeta) bool {
  return callMeta.FullMethod() != mailerv1.MailerService_Send_FullMethodName
}

func Auth(claims *sharedJwt.UserClaims, meta interceptors.CallMeta) bool {
  switch meta.FullMethod() {
  case mailerv1.MailerService_Find_FullMethodName:
    fallthrough
  case mailerv1.MailerService_FindByIds_FullMethodName:
    fallthrough
  case mailerv1.MailerService_FindByTag_FullMethodName:
    return authUtil.ContainsPermissions(claims.Roles, constant.MAILER_PERMISSIONS[constant.MAIL_READ])
  case mailerv1.MailerService_Update_FullMethodName:
    return authUtil.ContainsPermissions(claims.Roles, constant.MAILER_PERMISSIONS[constant.MAIL_UPDATE])
  case mailerv1.MailerService_Remove_FullMethodName:
    return authUtil.ContainsPermissions(claims.Roles, constant.MAILER_PERMISSIONS[constant.MAIL_DELETE])
  case mailerv1.TagService_Find_FullMethodName:
    fallthrough
  case mailerv1.TagService_FindByIds_FullMethodName:
    fallthrough
  case mailerv1.TagService_FindByName_FullMethodName:
    return authUtil.ContainsPermissions(claims.Roles, constant.MAILER_PERMISSIONS[constant.MAIL_READ_TAG])
  case mailerv1.TagService_Create_FullMethodName:
    return authUtil.ContainsPermissions(claims.Roles, constant.MAILER_PERMISSIONS[constant.MAIL_CREATE_TAG])
  case mailerv1.TagService_Update_FullMethodName:
    return authUtil.ContainsPermissions(claims.Roles, constant.MAILER_PERMISSIONS[constant.MAIL_UPDATE_TAG])
  case mailerv1.TagService_Remove_FullMethodName:
    return authUtil.ContainsPermissions(claims.Roles, constant.MAILER_PERMISSIONS[constant.MAIL_DELETE_TAG])
  }
  logger.Warn("Unknown method")
  return true
}
