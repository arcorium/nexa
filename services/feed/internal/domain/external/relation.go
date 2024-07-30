package external

import (
  "context"
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/feed/internal/domain/dto"
)

type IRelationClient interface {
  GetFollowings(ctx context.Context, userId types.Id) (dto.GetFollowingResponseDTO, error)
}
