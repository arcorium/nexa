package service

import (
  "context"
  userUow "nexa/services/user/internal/app/uow"
  "nexa/services/user/internal/domain/dto"
  "nexa/services/user/internal/domain/mapper"
  "nexa/services/user/internal/domain/service"
  "nexa/shared/status"
  "nexa/shared/types"
  "nexa/shared/uow"
  "nexa/shared/util"
)

func NewProfile(uow uow.IUnitOfWork[userUow.UserStorage]) service.IProfile {
  return &profileService{unit: uow}
}

type profileService struct {
  unit uow.IUnitOfWork[userUow.UserStorage]
}

func (p profileService) Find(ctx context.Context, userIds []types.Id) ([]dto.ProfileResponseDTO, status.Object) {
  repo := p.unit.Repositories()
  profiles, err := repo.Profile().FindByIds(ctx, userIds...)

  if err != nil {
    return nil, status.FromRepository(err, status.NullCode)
  }
  return util.CastSlice(profiles, mapper.ToProfileResponse), status.Success()
}

func (p profileService) Update(ctx context.Context, input *dto.ProfileUpdateDTO) status.Object {
  profile := mapper.MapProfileUpdateDTO(input)
  repo := p.unit.Repositories()
  err := repo.Profile().Patch(ctx, &profile)
  if err != nil {
    return status.FromRepository(err, status.NullCode)
  }
  return status.Updated()
}

func (p profileService) UpdateAvatar(ctx context.Context, input *dto.ProfilePictureUpdateDTO) status.Object {
  // TODO: Communicate with FileStorage Service
  url := types.FilePath("")
  profile := mapper.MapProfilePictureUpdateDTO(input)
  profile.PhotoURL = url

  repo := p.unit.Repositories()
  err := repo.Profile().Patch(ctx, &profile)
  if err != nil {
    return status.FromRepository(err, status.NullCode)
  }
  return status.Updated()
}
