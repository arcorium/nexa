package mapper

import (
  authv1 "nexa/proto/gen/go/authentication/v1"
  "nexa/services/authentication/internal/domain/dto"
  "nexa/shared/wrapper"
)

func ToLoginDTO(input *authv1.LoginRequest) dto.LoginDTO {
  return dto.LoginDTO{
    Email:    input.Email,
    Password: input.Password,
  }
}

func ToRegisterDTO(input *authv1.RegisterRequest) dto.RegisterDTO {
  return dto.RegisterDTO{
    Username:  input.Username,
    Email:     input.Email,
    Password:  input.Password,
    FirstName: input.FirstName,
    LastName:  wrapper.NewNullable(input.LastName),
    Bio:       wrapper.NewNullable(input.Bio),
  }
}

func ToRefreshTokenDTO(input *authv1.RefreshTokenRequest) dto.RefreshTokenDTO {
  return dto.RefreshTokenDTO{
    AccessToken: input.AccessToken,
  }
}

func ToProtoCredential(responseDTO *dto.CredentialResponseDTO) *authv1.Credential {
  return &authv1.Credential{
    Id:     responseDTO.Id,
    Device: responseDTO.Device,
  }
}
