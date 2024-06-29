package external

import (
  "context"
  "nexa/services/authentication/internal/domain/dto"
  "nexa/shared/types"
)

type IRoleClient interface {
  GetUserRoles(ctx context.Context, userId types.Id) ([]dto.RoleResponseDTO, error)
}
