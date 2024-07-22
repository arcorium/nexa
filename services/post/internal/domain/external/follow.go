package external

import (
  "context"
  "github.com/arcorium/nexa/shared/types"
)

type IFollowClient interface {
  IsFollower(ctx context.Context, followerId types.Id, followedId types.Id) (bool, error)
}
