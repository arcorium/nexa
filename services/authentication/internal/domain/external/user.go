package external

import (
  "context"
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/authentication/internal/domain/dto"
)

type IUserClient interface {
  Validate(ctx context.Context, email types.Email, password types.Password) (dto.UserResponseDTO, error)
  Create(ctx context.Context, request *dto.RegisterDTO) (types.Id, error)
}
