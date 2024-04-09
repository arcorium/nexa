package service

import (
  "context"
  "nexa/services/user/internal/domain/dto"
  "nexa/shared/status"
  "nexa/shared/types"
)

type IProfile interface {
  Find(ctx context.Context, userIds []types.Id) ([]dto.ProfileResponseDTO, status.Object)
  Update(ctx context.Context, input *dto.ProfileUpdateDTO) status.Object
  UpdateAvatar(ctx context.Context, input *dto.ProfilePictureUpdateDTO) status.Object
}
