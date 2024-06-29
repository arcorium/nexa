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
  spanUtil "nexa/shared/span"
  "nexa/shared/status"
  "nexa/shared/types"
  sharedUtil "nexa/shared/util"
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

func (p *permissionService) Create(ctx context.Context, createDTO *dto.PermissionCreateDTO) (string, status.Object) {
  ctx, span := p.tracer.Start(ctx, "PermissionService.Create")
  defer span.End()

  domain, err := createDTO.ToDomain()
  if err != nil {
    spanUtil.RecordError(err, span)
    return "", status.ErrBadRequest(err)
  }

  err = p.permRepo.Create(ctx, &domain)
  if err != nil {
    return "", status.FromRepository(err, status.NullCode)
  }

  return domain.Id.String(), status.Created()
}

func (p *permissionService) Find(ctx context.Context, permIds ...string) ([]dto.PermissionResponseDTO, status.Object) {
  ctx, span := p.tracer.Start(ctx, "PermissionService.Find")
  defer span.End()

  ids, ierr := sharedUtil.CastSliceErrs(permIds, func(permId string) (types.Id, error) {
    return types.IdFromString(permId)
  })

  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return nil, status.ErrBadRequest(ierr)
  }

  permissions, err := p.permRepo.FindByIds(ctx, ids...)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, status.FromRepository(err, status.NullCode)
  }

  responseDTOS := sharedUtil.CastSliceP(permissions, mapper.ToPermissionResponseDTO)
  return responseDTOS, status.Success()
}

func (p *permissionService) FindByRoles(ctx context.Context, roleIds ...string) ([]dto.PermissionResponseDTO, status.Object) {
  ctx, span := p.tracer.Start(ctx, "PermissionService.FindByRoles")
  defer span.End()

  ids, ierr := sharedUtil.CastSliceErrs(roleIds, func(permId string) (types.Id, error) {
    return types.IdFromString(permId)
  })

  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return nil, status.ErrBadRequest(ierr)
  }

  permissions, err := p.permRepo.FindByRoleIds(ctx, ids...)
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
    return sharedDto.PagedElementResult[dto.PermissionResponseDTO]{}, status.ErrBadRequest(err)
  }

  responseDTOS := sharedUtil.CastSliceP(result.Data, mapper.ToPermissionResponseDTO)
  return sharedDto.NewPagedElementResult2(responseDTOS, input, result.Total), status.Success()
}

func (p *permissionService) Delete(ctx context.Context, permId string) status.Object {
  ctx, span := p.tracer.Start(ctx, "PermissionService.Delete")
  defer span.End()

  id, err := types.IdFromString(permId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrBadRequest(err)
  }

  err = p.permRepo.Delete(ctx, id)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }
  return status.Deleted()
}

//
//func (p *permissionService) Find(ctx context.Context, id types.Id) (dto.PermissionResponseDTO, status.Object) {
//  permission, err := p.permRepo.FindById(ctx, id)
//  return mapper.ToPermissionResponseDTO(&permission), status.FromRepository(err, status.NullCode)
//}
//
//func (p *permissionService) FindAll(ctx context.Context, input *sharedDto.PagedElementDTO) (sharedDto.PagedElementResult[dto.PermissionResponseDTO], status.Object) {
//  result, err := p.permRepo.FindAll(ctx, input.ToQueryParam())
//  responseDTOS := util.CastSlice(result.Data, mapper.ToPermissionResponseDTO)
//  return sharedDto.NewPagedElementOutput2(responseDTOS, input, result.Total), status.FromRepository(err, status.NullCode)
//}
//
//func (p *permissionService) Create(ctx context.Context, createDTO *dto.PermissionCreateDTO) (types.Id, status.Object) {
//  perm := createDTO.ToDomain()
//  err := p.permRepo.Create(ctx, &perm)
//  return perm.Id, status.FromRepository(err, status.NullCode)
//}
//
//func (p *permissionService) Delete(ctx context.Context, id types.Id) status.Object {
//  err := p.permRepo.Delete(ctx, id)
//  return status.FromRepository(err, status.NullCode)
//}
