package external

import (
  "context"
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/comment/internal/domain/dto"
)

type IReactionClient interface {
  // DeleteComments remove all reaction in the comments
  DeleteComments(ctx context.Context, commentIds ...types.Id) error
  GetCommentsCounts(ctx context.Context, commentIds ...types.Id) ([]dto.ReactionCountDTO, error)
}
