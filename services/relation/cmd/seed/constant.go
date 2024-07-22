package main

import (
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/relation/constant"
)

var RELATION_DEFAULT_PERMS = []types.Action{
  constant.RELATION_GET_FOLLOW,
  constant.RELATION_DELETE_FOLLOW,
  constant.RELATION_CREATE_BLOCK,
  constant.RELATION_GET_BLOCK,
  constant.RELATION_CREATE_FOLLOW,
  constant.RELATION_DELETE_BLOCK,
}

var RELATION_SUPER_PERMS = []types.Action{
  constant.RELATION_DELETE_BLOCK_ARB,
  constant.RELATION_DELETE_FOLLOW_ARB,
}
