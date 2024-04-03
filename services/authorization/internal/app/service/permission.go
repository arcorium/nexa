package service

import (
	"context"
	"nexa/services/authorization/internal/domain/dto"
	"nexa/services/authorization/internal/domain/mapper"
	"nexa/services/authorization/internal/domain/repository"
	"nexa/services/authorization/internal/domain/service"
	authUtil "nexa/services/authorization/shared/util"
	sharedDto "nexa/shared/dto"
	"nexa/shared/status"
	"nexa/shared/types"
	"nexa/shared/util"
)

func NewPermission(permission repository.IPermission) service.IPermission {
	return &permissionService{permRepo: permission}
}

type permissionService struct {
	permRepo repository.IPermission
}

func (p *permissionService) CheckUserPermission(ctx context.Context, permissionDTO *dto.CheckUserPermissionDTO) status.Object {
	// Validate
	userId := types.IdFromString(permissionDTO.UserId)
	if err := userId.Validate(); err != nil {
		return status.Error(status.BAD_REQUEST_ERROR, err)
	}

	// Get user permissions (bypass role)
	permissions, err := p.permRepo.FindByUserId(ctx, userId)
	if err != nil {
		return status.FromRepository(err, status.NullCode)
	}

	// Check permissions
	return util.Ternary(
		authUtil.HasPermission(permissions, permissionDTO.ToDomain()),
		status.Success(),
		status.ErrUnAuthorized(nil))
}

func (p *permissionService) Find(ctx context.Context, id types.Id) (dto.PermissionResponseDTO, status.Object) {
	permission, err := p.permRepo.FindById(ctx, id)
	return mapper.ToPermissionResponseDTO(&permission), status.FromRepository(err, status.NullCode)
}

func (p *permissionService) FindAll(ctx context.Context, input *sharedDto.PagedElementDTO) (sharedDto.PagedElementResult[dto.PermissionResponseDTO], status.Object) {
	result, err := p.permRepo.FindAll(ctx, input.ToQueryParam())
	responseDTOS := util.CastSlice(result.Data, mapper.ToPermissionResponseDTO)
	return sharedDto.NewPagedElementOutput2(responseDTOS, input, result.Total), status.FromRepository(err, status.NullCode)
}

func (p *permissionService) Create(ctx context.Context, createDTO *dto.PermissionCreateDTO) (types.Id, status.Object) {
	perm := createDTO.ToDomain()
	err := p.permRepo.Create(ctx, &perm)
	return perm.Id, status.FromRepository(err, status.NullCode)
}

func (p *permissionService) Delete(ctx context.Context, id types.Id) status.Object {
	err := p.permRepo.Delete(ctx, id)
	return status.FromRepository(err, status.NullCode)
}
