package constant

import (
  "github.com/arcorium/nexa/shared/types"
  authUtil "github.com/arcorium/nexa/shared/util/auth"
)

// User
const (
  RELATION_CREATE_FOLLOW types.Action = "create:follow"
  RELATION_GET_FOLLOW    types.Action = "get:follow"
  RELATION_DELETE_FOLLOW types.Action = "delete:follow"
  RELATION_CREATE_BLOCK  types.Action = "create:block"
  RELATION_GET_BLOCK     types.Action = "get:block"
  RELATION_DELETE_BLOCK  types.Action = "delete:block"
)

// Super
const (
  RELATION_DELETE_FOLLOW_ARB types.Action = "delete:follow:arb"
  RELATION_DELETE_BLOCK_ARB  types.Action = "delete:block:arb"
)

var RELATION_PERMISSIONS = authUtil.FullEncode(SERVICE_RESOURCE,
  RELATION_CREATE_FOLLOW,
  RELATION_GET_FOLLOW,
  RELATION_DELETE_FOLLOW,
  RELATION_DELETE_FOLLOW_ARB,
  RELATION_CREATE_BLOCK,
  RELATION_GET_BLOCK,
  RELATION_DELETE_BLOCK,
  RELATION_DELETE_BLOCK_ARB,
)
