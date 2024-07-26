package external

import (
  "context"
  "github.com/arcorium/nexa/shared/types"
)

type IPostClient interface {
  Validate(ctx context.Context, postId types.Id) (bool, error)
}
