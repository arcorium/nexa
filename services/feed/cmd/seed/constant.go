package main

import (
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/feed/constant"
)

var FEED_DEFAULT_PERMS = []types.Action{
  constant.FEED_GET,
}

var FEED_SUPER_PERMS = []types.Action{
  constant.FEED_GET_ARB,
}
