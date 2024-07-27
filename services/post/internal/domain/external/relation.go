package external

import (
  "context"
  "github.com/arcorium/nexa/shared/types"
)

type IRelationClient interface {
  IsFollower(ctx context.Context, followerId types.Id, followedId types.Id) (bool, error)
  IsBlocked(ctx context.Context, blockerId types.Id) (bool, error)
}
