package mapper

import (
  authv1 "github.com/arcorium/nexa/proto/gen/go/authentication/v1"
  sharedErr "github.com/arcorium/nexa/shared/errors"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  "google.golang.org/protobuf/types/known/durationpb"
  "nexa/services/authentication/internal/domain/dto"
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

func ToRefreshTokenDTO(req *authv1.RefreshTokenRequest) (dto.RefreshTokenDTO, error) {
  dtos := dto.RefreshTokenDTO{
    TokenType:   req.TokenType,
    AccessToken: req.AccessToken,
  }

  err := sharedUtil.ValidateStruct(&dtos)
  return dtos, err
}

func ToLogoutDTO(req *authv1.LogoutRequest) (dto.LogoutDTO, error) {
  var fieldErrors []sharedErr.FieldError

  userId, err := types.IdFromString(req.UserId)
  if err != nil {
    fieldErrors = append(fieldErrors, sharedErr.NewFieldError("user_id", err))
  }

  credIds, ierr := sharedUtil.CastSliceErrs(req.CredIds, types.IdFromString)
  if !ierr.IsNil() {
    fieldErrors = append(fieldErrors, ierr.ToFieldError("cred_ids"))
  }

  if len(fieldErrors) == 0 {
    return dto.LogoutDTO{}, sharedErr.GrpcFieldErrors2(fieldErrors...)
  }

  return dto.LogoutDTO{
    UserId:        userId,
    CredentialIds: credIds,
  }, nil
}

func ToProtoRefreshTokenResponse(respDTO *dto.RefreshTokenResponseDTO) *authv1.RefreshTokenResponse {
  return &authv1.RefreshTokenResponse{
    TokenType:   respDTO.TokenType,
    AccessToken: respDTO.AccessToken,
    ExpiryTime:  durationpb.New(respDTO.ExpiryTime),
  }
}

func ToProtoLoginResponse(respDTO *dto.LoginResponseDTO) *authv1.LoginResponse {
  return &authv1.LoginResponse{
    TokenType:   respDTO.TokenType,
    AccessToken: respDTO.Token,
    ExpiryTime:  durationpb.New(respDTO.ExpiryTime),
  }
}

func ToProtoCredential(responseDTO *dto.CredentialResponseDTO) *authv1.Credential {
  return &authv1.Credential{
    Id:     responseDTO.Id.String(),
    Device: responseDTO.Device,
  }
}
