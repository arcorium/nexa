package service

import (
  "context"
  "nexa/services/user/internal/domain/dto"
  "nexa/shared/status"
  "nexa/shared/types"
)

type IUser interface {
  Create(ctx context.Context, input *dto.UserCreateDTO) status.Object
  Update(ctx context.Context, input *dto.UserUpdateDTO) status.Object
  UpdatePassword(ctx context.Context, input *dto.UserUpdatePasswordDTO) status.Object
  UpdateVerified(ctx context.Context, id types.Id) status.Object
  ResetPassword(ctx context.Context, input *dto.UserResetPasswordDTO) status.Object
  BannedUser(ctx context.Context, input *dto.UserBannedDTO) status.Object
  FindByEmails(ctx context.Context, emails []types.Email) ([]dto.UserResponseDTO, status.Object)
  FindByIds(ctx context.Context, ids []types.Id) ([]dto.UserResponseDTO, status.Object)
  DeleteById(ctx context.Context, id types.Id) status.Object
}
