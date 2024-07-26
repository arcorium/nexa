package external

import (
  "context"
  "github.com/arcorium/nexa/shared/types"
)

type IUserClient interface {
  Validate(ctx context.Context, userIds ...types.Id) (bool, error)
}
