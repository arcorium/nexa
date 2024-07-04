package external

import (
  "context"
  "nexa/services/user/internal/domain/dto"
  "nexa/shared/types"
)

type IAuthenticationClient interface {
  // DeleteCredentials logout all user credentials
  DeleteCredentials(ctx context.Context, userId types.Id) error
  // GenerateToken create new token for specific usage and user
  GenerateToken(ctx context.Context, dto *dto.TokenGenerationDTO) (dto.TokenResponseDTO, error)
  // VerifyToken Verify the token and return the user id
  VerifyToken(ctx context.Context, verificationDTO *dto.TokenVerificationDTO) (types.Id, error)
}
