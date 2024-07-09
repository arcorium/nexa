package service

import (
  "context"
  "go.opentelemetry.io/otel/trace"
  "nexa/services/authorization/internal/domain/dto"
  "nexa/services/authorization/internal/domain/entity"
  "nexa/services/authorization/internal/domain/repository"
  "nexa/services/authorization/internal/domain/service"
  "nexa/services/authorization/util"
  sharedErr "nexa/shared/errors"
  sharedJwt "nexa/shared/jwt"
  "nexa/shared/status"
  "nexa/shared/types"
  sharedUtil "nexa/shared/util"
  "nexa/shared/util/auth"
  spanUtil "nexa/shared/util/span"
)

func NewAuthorization(roleRepo repository.IRole) service.IAuthorization {
  return &authorizationService{
    roleRepo: roleRepo,
    tracer:   util.GetTracer(),
  }
}

type authorizationService struct {
  roleRepo repository.IRole
  tracer   trace.Tracer
}

func (a *authorizationService) IsAuthorized(ctx context.Context, authDto *dto.IsAuthorizationDTO) status.Object {
  ctx, span := a.tracer.Start(ctx, "AuthorizationService.IsAuthorized")
  defer span.End()

  roles, err := a.roleRepo.FindByUserId(ctx, authDto.UserId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepositoryOverride(err, types.NewPair(status.NOT_AUTHORIZED_ERROR, sharedErr.ErrUnauthorized))
  }

  jwtRoles := sharedUtil.CastSliceP(roles, func(role *entity.Role) sharedJwt.Role {
    return role.ToJWT()
  })

  if !auth.ContainsPermission(jwtRoles, authDto.ExpectedPermission) {
    spanUtil.RecordError(sharedErr.ErrUnauthorizedPermission, span)
    return status.ErrUnAuthorized(sharedErr.ErrUnauthorizedPermission)
  }
  return status.Success()
}
