package constant

import (
  "github.com/arcorium/nexa/shared/types"
  authUtil "github.com/arcorium/nexa/shared/util/auth"
)

// Default
const (
  FILE_GET_METADATA types.Action = "get:md"
  FILE_STORE        types.Action = "store:file"
  FILE_UPDATE       types.Action = "update:file"
  FILE_DELETE       types.Action = "delete:file"
)

// Super
const (
  FILE_GET           types.Action = "get:file"
  FILE_STORE_PRIVATE types.Action = "store:file:priv"
)

var FILE_STORAGE_PERMISSIONS = authUtil.FullEncode(SERVICE_RESOURCE,
  FILE_GET,
  FILE_GET_METADATA,
  FILE_STORE,
  FILE_UPDATE,
  FILE_DELETE,
  FILE_STORE_PRIVATE,
)
