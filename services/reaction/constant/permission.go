package constant

import (
  "github.com/arcorium/nexa/shared/types"
  authUtil "github.com/arcorium/nexa/shared/util/auth"
)

// User
const (
  REACTION_CREATE     types.Action = "create:reaction"
  REACTION_GET        types.Action = "get:reaction"
  REACTION_DELETE     types.Action = "delete:reaction"
  REACTION_DELETE_ARB types.Action = "delete:reaction:arb"
)

var REACTION_PERMISSIONS = authUtil.FullEncode(SERVICE_RESOURCE,
  REACTION_CREATE,
  REACTION_GET,
  REACTION_DELETE,
  REACTION_DELETE_ARB,
)
