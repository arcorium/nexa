package mapper

import (
  authv1 "nexa/proto/gen/go/authentication/v1"
  "nexa/services/authentication/internal/domain/dto"
  "nexa/shared/wrapper"
)

func ToLoginDTO(req *authv1.LoginRequest) dto.LoginDTO {
  return dto.LoginDTO{
    Email:    req.Email,
    Password: req.Password,
  }
}

func ToRegisterDTO(req *authv1.RegisterRequest) dto.RegisterDTO {
  return dto.RegisterDTO{
    Username:  req.Username,
    Email:     req.Email,
    Password:  req.Password,
    FirstName: req.FirstName,
    LastName:  wrapper.NewNullable(req.LastName),
    Bio:       wrapper.NewNullable(req.Bio),
  }
}

func ToRefreshTokenDTO(input *authv1.RefreshTokenRequest) dto.RefreshTokenDTO {
  return dto.RefreshTokenDTO{
    AccessToken: input.AccessToken,
  }
}

func ToLogoutDTO(input *authv1.LogoutRequest) dto.LogoutDTO {
  return dto.LogoutDTO{
    UserId:        input.UserId,
    CredentialIds: input.CredIds,
  }
}

func ToProtoRefreshTokenResponse(input *dto.RefreshTokenResponseDTO) *authv1.RefreshTokenResponse {
  return &authv1.RefreshTokenResponse{
    Type:        input.TokenType,
    AccessToken: input.AccessToken,
  }
}

func ToProtoLoginResponse(dto *dto.LoginResponseDTO) *authv1.LoginResponse {
  return &authv1.LoginResponse{
    TokenType:   dto.TokenType,
    AccessToken: dto.Token,
  }
}

func ToProtoCredential(responseDTO *dto.CredentialResponseDTO) *authv1.Credential {
  return &authv1.Credential{
    Id:     responseDTO.Id,
    Device: responseDTO.Device,
  }
}
