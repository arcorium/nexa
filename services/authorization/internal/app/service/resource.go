package service

import (
	"context"
	"nexa/services/authorization/internal/domain/dto"
	"nexa/services/authorization/internal/domain/mapper"
	"nexa/services/authorization/internal/domain/repository"
	"nexa/services/authorization/internal/domain/service"
	sharedDto "nexa/shared/dto"
	"nexa/shared/status"
	"nexa/shared/types"
	"nexa/shared/util"
)

func NewResource(repo repository.IResource) service.IResource {
	return &resourceService{repo: repo}
}

type resourceService struct {
	repo repository.IResource
}

func (r *resourceService) Find(ctx context.Context, id types.Id) (dto.ResourceResponseDTO, status.Object) {
	resource, err := r.repo.FindById(ctx, id)
	return mapper.ToResourceResponseDTO(&resource), status.FromRepository(err, status.NullCode)
}

func (r *resourceService) FindAll(ctx context.Context, input *sharedDto.PagedElementDTO) (sharedDto.PagedElementResult[dto.ResourceResponseDTO], status.Object) {
	result, err := r.repo.FindAll(ctx, input.ToQueryParam())
	responseDTOS := util.CastSlice(result.Data, mapper.ToResourceResponseDTO)

	return sharedDto.NewPagedElementOutput2(responseDTOS, input, result.Total), status.FromRepository(err, status.NullCode)
}

func (r *resourceService) Create(ctx context.Context, createDTO *dto.ResourceCreateDTO) (types.Id, status.Object) {
	resource := createDTO.ToDomain()
	err := r.repo.Create(ctx, &resource)
	return resource.Id, status.FromRepository(err, status.NullCode)
}

func (r *resourceService) Update(ctx context.Context, updateDTO *dto.ResourceUpdateDTO) status.Object {
	resource := updateDTO.ToDomain()
	err := r.repo.Patch(ctx, &resource)
	return status.FromRepository(err, status.NullCode)
}

func (r *resourceService) Delete(ctx context.Context, id types.Id) status.Object {
	err := r.repo.DeleteById(ctx, id)
	return status.FromRepository(err, status.NullCode)
}
