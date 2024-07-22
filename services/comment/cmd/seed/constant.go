package main

import (
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/comment/constant"
)

var COMMENT_DEFAULT_PERMS = []types.Action{
  constant.COMMENT_CREATE,
  constant.COMMENT_UPDATE,
  constant.COMMENT_GET,
  constant.COMMENT_DELETE,
}

var COMMENT_SUPER_PERMS = []types.Action{
  constant.COMMENT_DELETE_ARB,
}
