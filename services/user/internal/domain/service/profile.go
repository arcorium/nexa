package service

import (
	"context"
	"nexa/services/user/internal/domain/dto"
	"nexa/shared/status"
	"nexa/shared/types"
)

type IProfile interface {
	Find(ctx context.Context, userIds []types.Id) ([]dto.ProfileResponse, status.Object)
	Update(ctx context.Context, input *dto.ProfileUpdateInput) status.Object
	UpdateAvatar(ctx context.Context, input *dto.ProfilePictureUpdateInput) status.Object
}
