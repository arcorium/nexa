package mapper

import (
	"nexa/services/authentication/internal/domain/dto"
	"nexa/services/authentication/shared/proto"
	"nexa/shared/wrapper"
)

func ToLoginDTO(input *proto.LoginInput) dto.LoginDTO {
	return dto.LoginDTO{
		Email:    input.Email,
		Password: input.Password,
	}
}

func ToRegisterDTO(input *proto.RegisterInput) dto.RegisterDTO {
	return dto.RegisterDTO{
		Username:  input.Username,
		Email:     input.Email,
		Password:  input.Password,
		FirstName: input.FirstName,
		LastName:  wrapper.NewNullable(input.LastName),
		Bio:       wrapper.NewNullable(input.Bio),
	}
}

func ToRefreshTokenDTO(input *proto.RefreshTokenInput) dto.RefreshTokenDTO {
	return dto.RefreshTokenDTO{
		AccessToken: input.AccessToken,
	}
}

func ToCredentialResponse(responseDTO *dto.CredentialResponseDTO) *proto.CredentialResponse {
	return &proto.CredentialResponse{
		Id:     responseDTO.Id,
		Device: responseDTO.Device,
	}
}
