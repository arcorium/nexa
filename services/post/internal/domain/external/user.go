package external

import (
  "context"
  "github.com/arcorium/nexa/shared/types"
)

type IUserClient interface {
  GetUserNames(ctx context.Context, userIds ...types.Id) ([]string, error)
}
