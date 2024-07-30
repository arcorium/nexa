package external

import (
  "context"
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/feed/internal/domain/dto"
)

type IPostClient interface {
  GetUsers(ctx context.Context, limit uint64, userIds ...types.Id) (dto.GetUsersPostResponseDTO, error)
}
