package mapper

import (
  authv1 "nexa/proto/gen/go/authentication/v1"
  "nexa/services/authentication/internal/domain/dto"
  sharedErr "nexa/shared/errors"
  "nexa/shared/types"
  sharedUtil "nexa/shared/util"
)

func ToLoginDTO(req *authv1.LoginRequest) (dto.LoginDTO, error) {
  var fieldErrors []sharedErr.FieldError
  if len(req.DeviceName) == 0 {
    fieldErrors = append(fieldErrors, sharedErr.NewFieldError("device_name", sharedErr.ErrFieldEmpty))
  }

  email, err := types.EmailFromString(req.Email)
  if err != nil {
    fieldErrors = append(fieldErrors, sharedErr.NewFieldError("usage", err))
  }

  if len(fieldErrors) == 0 {
    return dto.LoginDTO{}, sharedErr.GrpcFieldErrors2(fieldErrors...)
  }

  password := types.PasswordFromString(req.Password)

  return dto.LoginDTO{
    Email:      email,
    Password:   password,
    DeviceName: req.DeviceName,
  }, nil
}

func ToRegisterDTO(req *authv1.RegisterRequest) (dto.RegisterDTO, error) {
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

func ToRefreshTokenDTO(req *authv1.RefreshTokenRequest) (dto.RefreshTokenDTO, error) {
  dtos := dto.RefreshTokenDTO{
    TokenType:   req.Type,
    AccessToken: req.AccessToken,
  }

  err := sharedUtil.ValidateStruct(&dtos)
  return dtos, err
}

func ToLogoutDTO(input *authv1.LogoutRequest) (dto.LogoutDTO, error) {
  var fieldErrors []sharedErr.FieldError

  // Determines if user id is empty
  var userId *types.Id = nil
  if input.UserId != nil {
    id, err := types.IdFromString(*input.UserId)
    userId = &id
    if err != nil {
      fieldErrors = append(fieldErrors, sharedErr.NewFieldError("user_id", err))
    }
  }

  credIds, ierr := sharedUtil.CastSliceErrs(input.CredIds, types.IdFromString)
  if !ierr.IsNil() {
    fieldErrors = append(fieldErrors, ierr.ToFieldError("cred_ids"))
  }

  if len(fieldErrors) == 0 {
    return dto.LogoutDTO{}, sharedErr.GrpcFieldErrors2(fieldErrors...)
  }

  return dto.LogoutDTO{
    UserId:        types.NewNullable(userId),
    CredentialIds: credIds,
  }, nil
}

func ToProtoRefreshTokenResponse(dtos *dto.RefreshTokenResponseDTO) *authv1.RefreshTokenResponse {
  return &authv1.RefreshTokenResponse{
    Type:        dtos.TokenType,
    AccessToken: dtos.AccessToken,
  }
}

func ToProtoLoginResponse(dtos *dto.LoginResponseDTO) *authv1.LoginResponse {
  return &authv1.LoginResponse{
    TokenType:   dtos.TokenType,
    AccessToken: dtos.Token,
  }
}

func ToProtoCredential(responseDTO *dto.CredentialResponseDTO) *authv1.Credential {
  return &authv1.Credential{
    Id:     responseDTO.Id.String(),
    Device: responseDTO.Device,
  }
}
