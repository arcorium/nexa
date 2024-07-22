package constant

import (
  "github.com/arcorium/nexa/shared/types"
  authUtil "github.com/arcorium/nexa/shared/util/auth"
)

// User
const (
  COMMENT_CREATE     types.Action = "create:comment"
  COMMENT_UPDATE     types.Action = "update:comment"
  COMMENT_GET        types.Action = "get:comment"
  COMMENT_DELETE     types.Action = "delete:comment"
  COMMENT_DELETE_ARB types.Action = "delete:comment:arb"
)

var COMMENT_PERMISSIONS = authUtil.FullEncode(SERVICE_RESOURCE,
  COMMENT_CREATE,
  COMMENT_UPDATE,
  COMMENT_GET,
  COMMENT_DELETE,
  COMMENT_DELETE_ARB,
)
