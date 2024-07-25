package external

import (
  "context"
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/authentication/internal/domain/dto"
)

type IRoleClient interface {
  SetUserAsDefault(ctx context.Context, userId types.Id) error
  SetUserRoles(ctx context.Context, userId types.Id, roleIds ...types.Id) error
  GetUserRoles(ctx context.Context, userId types.Id) ([]dto.RoleResponseDTO, error)
  RemoveUserRoles(ctx context.Context, userId types.Id) error
}
