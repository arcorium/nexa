package constant

import (
  "github.com/arcorium/nexa/shared/types"
  authUtil "github.com/arcorium/nexa/shared/util/auth"
)

// User
const (
  POST_GET_ARB    types.Action = "get:post:arb"
  POST_GET        types.Action = "get:post"
  POST_CREATE     types.Action = "create:post"
  POST_UPDATE     types.Action = "update:post"
  POST_DELETE     types.Action = "delete:post"
  POST_DELETE_ARB types.Action = "delete:post:arb"
)

var POST_PERMISSIONS = authUtil.FullEncode(SERVICE_RESOURCE,
  POST_GET_ARB,
  POST_GET,
  POST_CREATE,
  POST_UPDATE,
  POST_DELETE,
  POST_DELETE_ARB,
)
