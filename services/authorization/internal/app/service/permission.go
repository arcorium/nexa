package service

import (
  "context"
  sharedDto "github.com/arcorium/nexa/shared/dto"
  "github.com/arcorium/nexa/shared/status"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "go.opentelemetry.io/otel/trace"
  "nexa/services/authorization/internal/domain/dto"
  "nexa/services/authorization/internal/domain/entity"
  "nexa/services/authorization/internal/domain/mapper"
  "nexa/services/authorization/internal/domain/repository"
  "nexa/services/authorization/internal/domain/service"
  "nexa/services/authorization/util"
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
    return types.NullId(), status.FromRepositoryExist(err)
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

func (p *permissionService) GetAll(ctx context.Context, input *sharedDto.PagedElementDTO) (sharedDto.PagedElementResult[dto.PermissionResponseDTO], status.Object) {
  ctx, span := p.tracer.Start(ctx, "PermissionService.Get")
  defer span.End()

  result, err := p.permRepo.Get(ctx, input.ToQueryParam())
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

func (p *permissionService) Seed(ctx context.Context, seedDTO []dto.PermissionCreateDTO) ([]types.Id, status.Object) {
  ctx, span := p.tracer.Start(ctx, "PermissionService.Seed")
  defer span.End()

  entities, ierr := sharedUtil.CastSliceErrsP(seedDTO, func(createDto *dto.PermissionCreateDTO) (entity.Permission, error) {
    return createDto.ToDomain()
  })
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return nil, status.ErrInternal(ierr)
  }

  err := p.permRepo.Creates(ctx, entities...)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, status.FromRepositoryExist(err)
  }

  ids := sharedUtil.CastSliceP(entities, func(perm *entity.Permission) types.Id {
    return perm.Id
  })

  return ids, status.Created()
}
