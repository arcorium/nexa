package mapper

import (
	"nexa/services/user/shared/domain/dto"
	"nexa/services/user/shared/proto"
	"nexa/shared/wrapper"
)

// Proto
func ToDTOCreateInput(request *proto.CreateUserRequest) dto.UserCreateInput {
	return dto.UserCreateInput{
		Username:  request.Username,
		Email:     request.Email,
		Password:  request.Password,
		FirstName: request.FirstName,
		LastName:  wrapper.NewNullable(request.LastName),
		Bio:       wrapper.NewNullable(request.Bio),
	}
}

func ToDTOUserUpdateInput(request *proto.UpdateUserRequest) dto.UserUpdateInput {
	return dto.UserUpdateInput{
		Id:       request.Id,
		Username: wrapper.NewNullable(request.Username),
		Email:    wrapper.NewNullable(request.Email),
	}
}

func ToDTOUserUpdatePasswordInput(request *proto.UpdateUserPasswordRequest) dto.UserUpdatePasswordInput {
	return dto.UserUpdatePasswordInput{
		Id:           request.Id,
		LastPassword: request.LastPassword,
		NewPassword:  request.NewPassword,
	}
}

func ToDTOUserBannedInput(request *proto.BannedUserRequest) dto.UserBannedInput {
	return dto.UserBannedInput{
		Id:       request.Id,
		Duration: request.Duration.AsDuration(),
	}
}

func ToDTOUserResetPasswordInput(request *proto.ResetUserPasswordRequest) dto.UserResetPasswordInput {
	return dto.UserResetPasswordInput{
		Id:          request.Id,
		NewPassword: request.NewPassword,
	}
}

func ToProtoUserResponse(user *dto.UserResponse) proto.User {
	return proto.User{
		Id:         user.Id.Underlying().String(),
		Username:   user.Username,
		Email:      user.Email,
		IsVerified: user.IsVerified,
	}
}
