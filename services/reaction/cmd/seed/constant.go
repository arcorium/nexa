package main

import (
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/reaction/constant"
)

var REACTION_DEFAULT_PERMS = []types.Action{
  constant.REACTION_CREATE,
  constant.REACTION_GET,
  constant.REACTION_DELETE,
}

var REACTION_SUPER_PERMS = []types.Action{
  constant.REACTION_DELETE_ARB,
}
