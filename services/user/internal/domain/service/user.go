package service

import (
  "context"
  "nexa/services/user/internal/domain/dto"
  sharedDto "nexa/shared/dto"
  "nexa/shared/status"
  "nexa/shared/types"
)

type IUser interface {
  Create(ctx context.Context, input *dto.UserCreateDTO) (string, status.Object)
  Update(ctx context.Context, input *dto.UserUpdateDTO) status.Object
  UpdatePassword(ctx context.Context, input *dto.UserUpdatePasswordDTO) status.Object
  //UpdateVerified(ctx context.Context, userId string) status.Object
  //FindByEmails(ctx context.Context, emails []types.Email) ([]dto.UserResponseDTO, status.Object)
  BannedUser(ctx context.Context, input *dto.UserBannedDTO) status.Object
  FindAll(ctx context.Context, pagedDto sharedDto.PagedElementDTO) (sharedDto.PagedElementResult[dto.UserResponseDTO], status.Object)
  FindByIds(ctx context.Context, ids ...types.Id) ([]dto.UserResponseDTO, status.Object)
  DeleteById(ctx context.Context, id types.Id) status.Object
  Validate(ctx context.Context, email types.Email, password types.Password) (dto.UserResponseDTO, status.Object)

  VerifyEmail(ctx context.Context, token string) status.Object
  EmailVerificationRequest(ctx context.Context) (dto.TokenResponseDTO, status.Object)
  ForgotPassword(ctx context.Context, email types.Email) (dto.TokenResponseDTO, status.Object)
  ResetPassword(ctx context.Context, input *dto.UserResetPasswordDTO) status.Object
}
