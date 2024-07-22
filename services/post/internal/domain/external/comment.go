package external

import (
  "context"
  "github.com/arcorium/nexa/shared/types"
)

type ICommentClient interface {
  GetPostCounts(ctx context.Context, postIds ...types.Id) ([]uint64, error)
  // DeletePostsComments delete comments on each posts
  DeletePostsComments(ctx context.Context, postIds ...types.Id) error
}
