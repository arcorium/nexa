package mapper

import (
  userv1 "github.com/arcorium/nexa/proto/gen/go/user/v1"
  sharedErr "github.com/arcorium/nexa/shared/errors"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  "nexa/services/user/internal/domain/dto"
)

func ToUserCreateDTO(request *userv1.CreateUserRequest) (dto.UserCreateDTO, error) {
  // Hash Password
  pass := types.PasswordFromString(request.Password)

  email, err := types.EmailFromString(request.Email)
  if err != nil {
    return dto.UserCreateDTO{}, err
  }

  dtos := dto.UserCreateDTO{
    Username:  request.Username,
    Email:     email,
    Password:  pass,
    FirstName: request.FirstName,
    LastName:  types.NewNullable(request.LastName),
    Bio:       types.NewNullable(request.Bio),
  }

  err = sharedUtil.ValidateStruct(&dtos)
  return dtos, err
}

func ToUserUpdateDTO(request *userv1.UpdateUserRequest) (dto.UserUpdateDTO, error) {
  id, err := types.IdFromString(request.Id)
  if err != nil {
    err = sharedErr.GrpcFieldErrors2(sharedErr.NewFieldError("id", err))
    return dto.UserUpdateDTO{}, err
  }

  var emails *types.Email = nil
  if request.Email != nil {
    email, err := types.EmailFromString(*request.Email)
    if err != nil {
      return dto.UserUpdateDTO{}, sharedErr.NewFieldError("email", err).ToGrpcError()
    }
    emails = &email
  }

  return dto.UserUpdateDTO{
    Id:       id,
    Username: types.NewNullable(request.Username),
    Email:    types.NewNullable(emails),
  }, nil
}

func ToUserUpdatePasswordDTO(request *userv1.UpdateUserPasswordRequest) (dto.UserUpdatePasswordDTO, error) {
  id, err := types.IdFromString(request.Id)
  if err != nil {
    err = sharedErr.GrpcFieldErrors2(sharedErr.NewFieldError("id", err))
    return dto.UserUpdatePasswordDTO{}, err
  }

  lastPassword := types.PasswordFromString(request.LastPassword)
  newPassword := types.PasswordFromString(request.NewPassword)

  dtos := dto.UserUpdatePasswordDTO{
    Id:           id,
    LastPassword: lastPassword,
    NewPassword:  newPassword,
  }

  // Validate
  err = sharedUtil.ValidateStruct(&dtos)
  if err != nil {
    return dto.UserUpdatePasswordDTO{}, err
  }

  return dtos, nil
}

func ToDTOUserBannedInput(request *userv1.BannedUserRequest) (dto.UserBannedDTO, error) {
  id, err := types.IdFromString(request.Id)
  if err != nil {
    err = sharedErr.GrpcFieldErrors2(sharedErr.NewFieldError("id", err))
    return dto.UserBannedDTO{}, err
  }

  dtos := dto.UserBannedDTO{
    Id:       id,
    Duration: request.Duration.AsDuration(),
  }

  err = sharedUtil.ValidateStruct(&dtos)
  return dtos, err
}

func ToResetUserPasswordDTO(request *userv1.ResetUserPasswordRequest) (dto.ResetUserPasswordDTO, error) {
  // Empty validation
  eerr := sharedUtil.StringEmptyValidates(
    types.NewField("new_password", request.NewPassword))
  if !eerr.IsNil() {
    return dto.ResetUserPasswordDTO{}, eerr.ToGRPCError()
  }

  var userId *types.Id = nil
  if request.UserId != nil {
    id, err := types.IdFromString(*request.UserId)
    if err != nil {
      return dto.ResetUserPasswordDTO{}, sharedErr.NewFieldError("user_id", err).ToGrpcError()
    }
    userId = &id
  }

  password := types.PasswordFromString(request.NewPassword)

  dtos := dto.ResetUserPasswordDTO{
    Token:       types.NewNullable(request.Token),
    UserId:      types.NewNullable(userId),
    LogoutAll:   request.LogoutAll,
    NewPassword: password,
  }

  err := sharedUtil.ValidateStruct(&dtos)
  return dtos, err
}

func ToProtoUser(responseDTO *dto.UserResponseDTO) *userv1.User {
  return &userv1.User{
    Id:         responseDTO.Id.String(),
    Username:   responseDTO.Username,
    Email:      responseDTO.Email.String(),
    IsVerified: responseDTO.IsVerified,
    Profile:    ToProtoProfile(responseDTO.Profile),
  }
}
