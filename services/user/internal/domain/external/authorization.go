package external

import (
  "context"
  "nexa/shared/types"
)

type IAuthorizationClient interface {
  // HasPermission check is the user has permission to access the resource with those actions
  HasPermission(ctx context.Context, roleId types.Id, resource string, actions ...string) error
}
