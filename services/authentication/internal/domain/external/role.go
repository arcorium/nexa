package external

import (
  "context"
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/authentication/internal/domain/dto"
)

type IRoleClient interface {
  GetUserRoles(ctx context.Context, userId types.Id) ([]dto.RoleResponseDTO, error)
  RemoveUserRoles(ctx context.Context, userId types.Id) error
}
