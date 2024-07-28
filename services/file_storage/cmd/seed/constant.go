package main

import (
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/file_storage/constant"
)

var FILE_DEFAULT_PERMS = []types.Action{
  constant.FILE_GET_METADATA,
  constant.FILE_STORE,
  constant.FILE_UPDATE,
  constant.FILE_DELETE,
}

var FILE_SUPER_PERMS = []types.Action{
  constant.FILE_GET,
  constant.FILE_STORE_PRIVATE,
}
