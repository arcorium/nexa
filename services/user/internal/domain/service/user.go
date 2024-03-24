package service

import (
	"context"
	"nexa/services/user/shared/domain/dto"
	"nexa/shared/status"
	"nexa/shared/types"
)

type IUser interface {
	Create(ctx context.Context, input *dto.UserCreateInput) status.Object
	Update(ctx context.Context, input *dto.UserUpdateInput) status.Object
	UpdatePassword(ctx context.Context, input *dto.UserUpdatePasswordInput) status.Object
	UpdateVerified(ctx context.Context, id types.Id) status.Object
	ResetPassword(ctx context.Context, input *dto.UserResetPasswordInput) status.Object
	BannedUser(ctx context.Context, input *dto.UserBannedInput) status.Object
	FindByEmails(ctx context.Context, emails []types.Email) ([]dto.UserResponse, status.Object)
	FindByIds(ctx context.Context, ids []types.Id) ([]dto.UserResponse, status.Object)
	DeleteById(ctx context.Context, id types.Id) status.Object
}
