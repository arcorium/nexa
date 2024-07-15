package external

import (
  "context"
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/authentication/internal/domain/dto"
)

type ITokenClient interface {
  // Generate new token for specific usage and user
  Generate(ctx context.Context, dto *dto.TokenGenerationDTO) (dto.TokenResponseDTO, error)
  // Verify the token and return the user id related to the token
  Verify(ctx context.Context, verificationDTO *dto.TokenVerificationDTO) (types.Id, error)
}
