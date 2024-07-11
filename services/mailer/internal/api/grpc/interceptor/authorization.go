package interceptor

import (
  "context"
  mailerv1 "github.com/arcorium/nexa/proto/gen/go/mailer/v1"
  sharedJwt "github.com/arcorium/nexa/shared/jwt"
  "github.com/arcorium/nexa/shared/logger"
  authUtil "github.com/arcorium/nexa/shared/util/auth"
  "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
  "nexa/services/mailer/constant"
)

func AuthSkipSelector(_ context.Context, callMeta interceptors.CallMeta) bool {
  return callMeta.Service == mailerv1.MailerService_ServiceDesc.ServiceName ||
      callMeta.FullMethod() == mailerv1.MailerService_Send_FullMethodName ||
      callMeta.FullMethod() == mailerv1.TagService_Find_FullMethodName ||
      callMeta.FullMethod() == mailerv1.TagService_FindByIds_FullMethodName ||
      callMeta.FullMethod() == mailerv1.TagService_FindByName_FullMethodName
}

func Auth(claims *sharedJwt.UserClaims, meta interceptors.CallMeta) bool {
  switch meta.FullMethod() {
  case mailerv1.MailerService_Find_FullMethodName:
    fallthrough
  case mailerv1.MailerService_FindByIds_FullMethodName:
    fallthrough
  case mailerv1.MailerService_FindByTag_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.MAILER_PERMISSIONS[constant.MAIL_READ])
  case mailerv1.MailerService_Update_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.MAILER_PERMISSIONS[constant.MAIL_UPDATE])
  case mailerv1.MailerService_Remove_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.MAILER_PERMISSIONS[constant.MAIL_DELETE])
  case mailerv1.TagService_Create_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.MAILER_PERMISSIONS[constant.MAIL_CREATE_TAG])
  case mailerv1.TagService_Update_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.MAILER_PERMISSIONS[constant.MAIL_UPDATE_TAG])
  case mailerv1.TagService_Remove_FullMethodName:
    return authUtil.ContainsPermission(claims.Roles, constant.MAILER_PERMISSIONS[constant.MAIL_DELETE_TAG])
  default:
    logger.Warn("Unknown method")
  }
  return true
}
