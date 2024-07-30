package constant

import (
  "github.com/arcorium/nexa/shared/types"
  authUtil "github.com/arcorium/nexa/shared/util/auth"
)

// User
const (
  FEED_GET types.Action = "get"
)

// Super
const (
  FEED_GET_ARB types.Action = "get:arb"
)

var FEED_PERMISSIONS = authUtil.FullEncode(SERVICE_RESOURCE,
  FEED_GET,
  FEED_GET_ARB,
)
