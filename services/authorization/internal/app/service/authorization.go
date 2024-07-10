package service

import (
  "context"
  sharedErr "github.com/arcorium/nexa/shared/errors"
  sharedJwt "github.com/arcorium/nexa/shared/jwt"
  "github.com/arcorium/nexa/shared/status"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  "github.com/arcorium/nexa/shared/util/auth"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "go.opentelemetry.io/otel/trace"
  "nexa/services/authorization/internal/domain/dto"
  "nexa/services/authorization/internal/domain/entity"
  "nexa/services/authorization/internal/domain/repository"
  "nexa/services/authorization/internal/domain/service"
  "nexa/services/authorization/util"
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
