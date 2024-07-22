package external

import (
  "context"
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/post/internal/domain/dto"
)

type IReactionClient interface {
  GetPostCounts(ctx context.Context, postIds ...types.Id) ([]dto.LikeCountResponseDTO, error)
  DeletePostsLikes(ctx context.Context, postIds ...types.Id) error
}
