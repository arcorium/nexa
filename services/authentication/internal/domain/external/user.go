package external

import (
  "context"
  "nexa/services/authentication/internal/domain/dto"
  "nexa/shared/types"
)

type IUserClient interface {
  Validate(ctx context.Context, email types.Email, password string) (dto.UserValidateResponseDTO, error)
  Create(ctx context.Context, request *dto.RegisterDTO) error
}
