package service

import (
  "context"
  "github.com/arcorium/nexa/shared/status"
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/user/internal/domain/dto"
)

type IProfile interface {
  Find(ctx context.Context, userIds ...types.Id) ([]dto.ProfileResponseDTO, status.Object)
  Update(ctx context.Context, input *dto.ProfileUpdateDTO) status.Object
  UpdateAvatar(ctx context.Context, input *dto.ProfileAvatarUpdateDTO) status.Object
}
