package mapper

import (
  authNv1 "github.com/arcorium/nexa/proto/gen/go/authentication/v1"
  sharedErr "github.com/arcorium/nexa/shared/errors"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  "nexa/services/authentication/internal/domain/dto"
)

func ToRegisterDTO(req *authNv1.RegisterRequest) (dto.RegisterDTO, error) {
  email, err := types.EmailFromString(req.Email)
  if err != nil {
    return dto.RegisterDTO{}, sharedErr.NewFieldError("email", err).ToGrpcError()
  }
  password := types.PasswordFromString(req.Password)

  return dto.RegisterDTO{
    Username:  req.Username,
    Email:     email,
    Password:  password,
    FirstName: req.FirstName,
    LastName:  types.NewNullable(req.LastName),
    Bio:       types.NewNullable(req.Bio),
  }, nil
}

func ToUserCreateDTO(request *authNv1.CreateUserRequest) (dto.UserCreateDTO, error) {
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

func ToUserUpdateDTO(request *authNv1.UpdateUserRequest) (dto.UserUpdateDTO, error) {
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
    Id:        id,
    Username:  types.NewNullable(request.Username),
    Email:     types.NewNullable(emails),
    FirstName: types.NewNullable(request.FirstName),
    LastName:  types.NewNullable(request.LastName),
    Bio:       types.NewNullable(request.Bio),
  }, nil
}

func ToUserUpdatePasswordDTO(request *authNv1.UpdateUserPasswordRequest) (dto.UserUpdatePasswordDTO, error) {
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

func ToDTOUserBannedInput(request *authNv1.BannedUserRequest) (dto.UserBannedDTO, error) {
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

func ToResetUserPasswordDTO(request *authNv1.ResetUserPasswordRequest) (dto.ResetUserPasswordDTO, error) {
  // Empty validation
  eerr := sharedUtil.StringEmptyValidates(
    types.NewField("new_password", request.NewPassword))
  if !eerr.IsNil() {
    return dto.ResetUserPasswordDTO{}, eerr.ToGRPCError()
  }

  userId, err := types.IdFromString(request.UserId)
  if err != nil {
    return dto.ResetUserPasswordDTO{}, sharedErr.NewFieldError("user_id", err).ToGrpcError()
  }

  password := types.PasswordFromString(request.NewPassword)

  dtos := dto.ResetUserPasswordDTO{
    UserId:      userId,
    LogoutAll:   request.LogoutAll,
    NewPassword: password,
  }

  err = sharedUtil.ValidateStruct(&dtos)
  return dtos, err
}

func ToResetPasswordByTokenDTO(request *authNv1.ResetPasswordByTokenRequest) (dto.ResetPasswordWithTokenDTO, error) {
  // Empty validation
  eerr := sharedUtil.StringEmptyValidates(
    types.NewField("token", request.Token),
    types.NewField("new_password", request.NewPassword),
  )
  if !eerr.IsNil() {
    return dto.ResetPasswordWithTokenDTO{}, eerr.ToGRPCError()
  }

  password := types.PasswordFromString(request.NewPassword)

  dtos := dto.ResetPasswordWithTokenDTO{
    Token:       request.Token,
    LogoutAll:   request.LogoutAll,
    NewPassword: password,
  }

  err := sharedUtil.ValidateStruct(&dtos)
  return dtos, err
}

func ToProtoUser(responseDTO *dto.UserResponseDTO) *authNv1.User {
  return &authNv1.User{
    Id:         responseDTO.Id.String(),
    Username:   responseDTO.Username,
    Email:      responseDTO.Email.String(),
    IsVerified: responseDTO.IsVerified,
    FirstName:  responseDTO.Profile.FirstName,
    LastName:   responseDTO.Profile.LastName,
    Bio:        responseDTO.Profile.Bio,
    ImagePath:  responseDTO.Profile.PhotoURL.String(),
  }
}
