package mapper

import (
	"nexa/services/user/shared/domain/dto"
	"nexa/services/user/shared/domain/entity"
	"nexa/shared/types"
	"time"
)

func MapUserCreateInput(input *dto.UserCreateInput) (entity.User, entity.Profile) {
	user := entity.User{
		Id:         types.NewId(),
		Username:   input.Username,
		Email:      types.EmailFromString(input.Email),
		Password:   types.PasswordFromString(input.Password),
		IsVerified: false,
		IsDeleted:  false,
	}
	profile := entity.Profile{
		Id:        user.Id,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Bio:       input.Bio,
	}
	return user, profile
}

func MapUserUpdateInput(input *dto.UserUpdateInput) entity.User {
	return entity.User{
		Id:         types.IdFromString(input.Id),
		Username:   input.Username,
		Email:      types.EmailFromString(input.Email),
		IsVerified: false,
		IsDeleted:  false,
	}
}

func MapUserUpdatePasswordInput(input *dto.UserUpdatePasswordInput) entity.User {
	return entity.User{
		Id:       types.IdFromString(input.Id),
		Password: types.PasswordFromString(input.NewPassword),
	}
}

func MapUserResetPasswordInput(input *dto.UserResetPasswordInput) entity.User {
	return entity.User{
		Id:       types.IdFromString(input.Id),
		Password: types.PasswordFromString(input.NewPassword),
	}
}

func MapUserBannedInput(input *dto.UserBannedInput) entity.User {
	return entity.User{
		Id:          types.IdFromString(input.Id),
		BannedUntil: time.Now().Add(input.Duration),
	}
}

func ToUserResponse(user *entity.User) dto.UserResponse {
	return dto.UserResponse{}
}
