package external

import (
  "context"
  "github.com/arcorium/nexa/shared/types"
)

type ICommentClient interface {
  Validate(ctx context.Context, commentIds ...types.Id) (bool, error)
}
