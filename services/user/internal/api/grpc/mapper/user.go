package mapper

import (
  proto "nexa/proto/generated/golang/user/v1"
  "nexa/services/user/internal/domain/dto"
  "nexa/shared/wrapper"
)

func ToDTOCreateInput(request *proto.CreateUserRequest) dto.UserCreateDTO {
  return dto.UserCreateDTO{
    Username:  request.Username,
    Email:     request.Email,
    Password:  request.Password,
    FirstName: request.FirstName,
    LastName:  wrapper.NewNullable(request.LastName),
    Bio:       wrapper.NewNullable(request.Bio),
  }
}

func ToDTOUserUpdateInput(request *proto.UpdateUserRequest) dto.UserUpdateDTO {
  return dto.UserUpdateDTO{
    Id:       request.Id,
    Username: wrapper.NewNullable(request.Username),
    Email:    wrapper.NewNullable(request.Email),
  }
}

func ToDTOUserUpdatePasswordInput(request *proto.UpdateUserPasswordRequest) dto.UserUpdatePasswordDTO {
  return dto.UserUpdatePasswordDTO{
    Id:           request.Id,
    LastPassword: request.LastPassword,
    NewPassword:  request.NewPassword,
  }
}

func ToDTOUserBannedInput(request *proto.BannedUserRequest) dto.UserBannedDTO {
  return dto.UserBannedDTO{
    Id:       request.Id,
    Duration: request.Duration.AsDuration(),
  }
}

func ToDTOUserResetPasswordInput(request *proto.ResetUserPasswordRequest) dto.UserResetPasswordDTO {
  return dto.UserResetPasswordDTO{
    Id:          request.Id,
    NewPassword: request.NewPassword,
  }
}

func ToProtoUser(responseDTO *dto.UserResponseDTO) *proto.User {
  return &proto.User{
    Id:         responseDTO.Id,
    Username:   responseDTO.Username,
    Email:      responseDTO.Email,
    IsVerified: responseDTO.IsVerified,
    Profile:    ToProtoProfile(responseDTO.Profile),
  }
}
