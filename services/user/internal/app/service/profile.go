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

func (p profileService) Find(ctx context.Context, userIds []types.Id) ([]dto.ProfileResponse, status.Object) {
	profiles, err := p.unit.Repositories().Profile().FindByIds(ctx, userIds...)
	if err != nil {
		return nil, status.FromRepository(err, status.NullCode)
	}
	return util.CastSlice(profiles, mapper.ToProfileResponse), status.Success()
}

func (p profileService) Update(ctx context.Context, input *dto.ProfileUpdateInput) status.Object {
	profile := mapper.MapProfileUpdateInput(input)
	err := p.unit.Repositories().Profile().Patch(ctx, &profile)
	if err != nil {
		return status.FromRepository(err, status.NullCode)
	}
	return status.Updated()
}

func (p profileService) UpdateAvatar(ctx context.Context, input *dto.ProfilePictureUpdateInput) status.Object {
	// TODO: Communicate with FileStorage Service
	url := types.FilePath("")
	profile := mapper.MapProfilePictureUpdateInput(input)
	profile.PhotoURL = url

	err := p.unit.Repositories().Profile().Patch(ctx, &profile)
	if err != nil {
		return status.FromRepository(err, status.NullCode)
	}
	return status.Updated()
}
