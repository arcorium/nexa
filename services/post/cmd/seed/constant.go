package main

import (
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/post/constant"
)

var POST_DEFAULT_PERMS = []types.Action{
  constant.POST_CREATE,
  constant.POST_UPDATE,
  constant.POST_GET,
  constant.POST_DELETE,
}

var POST_SUPER_PERMS = []types.Action{
  constant.POST_GET_ARB,
  constant.POST_DELETE_ARB,
}
