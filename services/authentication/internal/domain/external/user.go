package external

import (
	"context"
	"nexa/services/authentication/internal/domain/dto"
	userDto "nexa/services/user/shared/domain/dto"
	"nexa/shared/types"
)

type IUserClient interface {
	ValidateUser(ctx context.Context, email types.Email, password string) (userDto.UserResponseDTO, error)
	RegisterUser(ctx context.Context, request *dto.RegisterDTO) error
}
