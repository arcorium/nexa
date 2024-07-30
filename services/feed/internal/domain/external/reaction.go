package external

import (
  "context"
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/feed/internal/domain/dto"
)

type IReactionClient interface {
  GetPostReactionCounts(ctx context.Context, postIds ...types.Id) ([]dto.PostReactionCountDTO, error)
}
