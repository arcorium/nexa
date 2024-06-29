package service

import (
  "context"
  "go.opentelemetry.io/otel/trace"
  "nexa/services/authorization/internal/domain/dto"
  "nexa/services/authorization/internal/domain/entity"
  "nexa/services/authorization/internal/domain/repository"
  "nexa/services/authorization/internal/domain/service"
  "nexa/shared/auth"
  sharedErr "nexa/shared/errors"
  sharedJwt "nexa/shared/jwt"
  spanUtil "nexa/shared/span"
  "nexa/shared/status"
  "nexa/shared/types"
  sharedUtil "nexa/shared/util"
)

func NewAuthorization() service.IAuthorization {
  return &authorizationService{}
}

type authorizationService struct {
  roleRepo repository.IRole
  tracer   trace.Tracer
}

func (a *authorizationService) IsAuthorized(ctx context.Context, authDto *dto.IsAuthorizationDTO) status.Object {
  ctx, span := a.tracer.Start(ctx, "AuthorizationService.IsAuthorized")
  defer span.End()

  id, err := types.IdFromString(authDto.UserId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrBadRequest(err)
  }

  roles, err := a.roleRepo.FindByUserId(ctx, id)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }

  jwtRoles := sharedUtil.CastSliceP(roles, func(role *entity.Role) sharedJwt.Role {
    return role.ToJWT()
  })

  result := auth.ContainsPermissions(jwtRoles, authDto.ExpectedPermission...)
  if !result {
    spanUtil.RecordError(sharedErr.ErrUnauthorized, span)
    return status.ErrUnAuthorized(sharedErr.ErrUnauthorized)
  }
  return status.Success()
}
