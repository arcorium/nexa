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

func NewAction(action repository.IAction) service.IAction {
  return &actionService{repo: action}
}

type actionService struct {
  repo repository.IAction
}

func (a *actionService) Find(ctx context.Context, id types.Id) (dto.ActionResponseDTO, status.Object) {
  action, err := a.repo.FindById(ctx, id)
  stats := status.FromRepository(err, status.NullCode)
  return mapper.ToActionResponseDTO(&action), stats
}

func (a *actionService) FindAll(ctx context.Context, input *sharedDto.PagedElementDTO) (sharedDto.PagedElementResult[dto.ActionResponseDTO], status.Object) {
  actions, err := a.repo.FindAll(ctx, input.ToQueryParam())
  responseDTOS := util.CastSlice(actions.Data, mapper.ToActionResponseDTO)

  return sharedDto.NewPagedElementOutput2(responseDTOS, input, actions.Total), status.FromRepository(err, status.NullCode)
}

func (a *actionService) Create(ctx context.Context, input *dto.ActionCreateDTO) (types.Id, status.Object) {
  action := input.ToDomain()
  err := a.repo.Create(ctx, &action)
  return action.Id, status.FromRepository(err, status.NullCode)
}

func (a *actionService) Update(ctx context.Context, input *dto.ActionUpdateDTO) status.Object {
  action := input.ToDomain()
  err := a.repo.Patch(ctx, &action)
  return status.FromRepository(err, status.NullCode)
}

func (a *actionService) Delete(ctx context.Context, id types.Id) status.Object {
  err := a.repo.DeleteById(ctx, id)
  return status.FromRepository(err, status.NullCode)
}
