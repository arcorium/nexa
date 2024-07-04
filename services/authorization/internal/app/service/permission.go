package service

import (
  "context"
  "go.opentelemetry.io/otel/trace"
  "nexa/services/authorization/internal/domain/dto"
  "nexa/services/authorization/internal/domain/mapper"
  "nexa/services/authorization/internal/domain/repository"
  "nexa/services/authorization/internal/domain/service"
  "nexa/services/authorization/util"
  sharedDto "nexa/shared/dto"
  "nexa/shared/status"
  "nexa/shared/types"
  sharedUtil "nexa/shared/util"
  spanUtil "nexa/shared/util/span"
)

func NewPermission(permission repository.IPermission) service.IPermission {
  return &permissionService{
    permRepo: permission,
    tracer:   util.GetTracer(),
  }
}

type permissionService struct {
  permRepo repository.IPermission

  tracer trace.Tracer
}

func (p *permissionService) Create(ctx context.Context, createDTO *dto.PermissionCreateDTO) (types.Id, status.Object) {
  ctx, span := p.tracer.Start(ctx, "PermissionService.Create")
  defer span.End()

  domain, err := createDTO.ToDomain()
  if err != nil {
    spanUtil.RecordError(err, span)
    return types.NullId(), status.ErrInternal(err)
  }

  err = p.permRepo.Create(ctx, &domain)
  if err != nil {
    spanUtil.RecordError(err, span)
    return types.NullId(), status.FromRepository(err, status.NullCode)
  }

  return domain.Id, status.Created()
}

func (p *permissionService) Find(ctx context.Context, permIds ...types.Id) ([]dto.PermissionResponseDTO, status.Object) {
  ctx, span := p.tracer.Start(ctx, "PermissionService.FindByIds")
  defer span.End()

  permissions, err := p.permRepo.FindByIds(ctx, permIds...)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, status.FromRepository(err, status.NullCode)
  }

  responseDTOS := sharedUtil.CastSliceP(permissions, mapper.ToPermissionResponseDTO)
  return responseDTOS, status.Success()
}

func (p *permissionService) FindByRoles(ctx context.Context, roleIds ...types.Id) ([]dto.PermissionResponseDTO, status.Object) {
  ctx, span := p.tracer.Start(ctx, "PermissionService.FindByRoles")
  defer span.End()

  permissions, err := p.permRepo.FindByRoleIds(ctx, roleIds...)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, status.FromRepository(err, status.NullCode)
  }

  responseDTOS := sharedUtil.CastSliceP(permissions, mapper.ToPermissionResponseDTO)
  return responseDTOS, status.Success()
}

func (p *permissionService) FindAll(ctx context.Context, input *sharedDto.PagedElementDTO) (sharedDto.PagedElementResult[dto.PermissionResponseDTO], status.Object) {
  ctx, span := p.tracer.Start(ctx, "PermissionService.FindAll")
  defer span.End()

  result, err := p.permRepo.FindAll(ctx, input.ToQueryParam())
  if err != nil {
    spanUtil.RecordError(err, span)
    return sharedDto.PagedElementResult[dto.PermissionResponseDTO]{}, status.FromRepository(err, status.NullCode)
  }

  responseDTOS := sharedUtil.CastSliceP(result.Data, mapper.ToPermissionResponseDTO)
  return sharedDto.NewPagedElementResult2(responseDTOS, input, result.Total), status.Success()
}

func (p *permissionService) Delete(ctx context.Context, permId types.Id) status.Object {
  ctx, span := p.tracer.Start(ctx, "PermissionService.Delete")
  defer span.End()

  err := p.permRepo.Delete(ctx, permId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }
  return status.Deleted()
}
