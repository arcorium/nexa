package external

import (
  "context"
  "github.com/arcorium/nexa/shared/types"
)

type ICommentClient interface {
  GetPostCommentCounts(ctx context.Context, postIds ...types.Id) ([]uint64, error)
}
