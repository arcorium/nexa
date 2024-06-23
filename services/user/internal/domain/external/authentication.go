package external

import (
  "context"
  "nexa/services/user/internal/domain/dto"
)

type IAuthenticationClient interface {
  GenerateToken(ctx context.Context, dto *dto.TokenGenerationDTO) (dto.TokenResponseDTO, error)
}
