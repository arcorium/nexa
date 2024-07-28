package external

import (
  "context"
  "github.com/arcorium/nexa/shared/types"
)

type IUserClient interface {
  Validate(ctx context.Context, userId types.Id) (bool, error)
}
