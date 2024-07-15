package constant

import (
  "github.com/arcorium/nexa/shared/types"
  authUtil "github.com/arcorium/nexa/shared/util/auth"
)

// Default
const (
  AUTHZ_READ_ROLE        types.Action = "read:role"
  AUTHZ_READ_PERMISSION  types.Action = "read:perms"
  AUTHZ_DELETE_USER_ROLE types.Action = "delete:user:role"
)

// Super
const (
  AUTHZ_CREATE_ROLE             types.Action = "create:role"
  AUTHZ_DELETE_ROLE             types.Action = "delete:role"
  AUTHZ_UPDATE_ROLE             types.Action = "update:role"
  AUTHZ_DELETE_USER_ROLE_OTHER  types.Action = "delete:user:role:arb"
  AUTHZ_CREATE_PERMISSION       types.Action = "create:perms"
  AUTHZ_DELETE_PERMISSION       types.Action = "delete:perms"
  AUTHZ_MODIFY_USER_ROLE        types.Action = "modify:role"
  AUTHZ_MODIFY_ROLE_PERMISSIONS types.Action = "modify:perms"
)

var AUTHZ_PERMISSIONS = authUtil.FullEncode(SERVICE_RESOURCE,
  AUTHZ_READ_ROLE,
  AUTHZ_READ_PERMISSION,
  AUTHZ_DELETE_USER_ROLE,
  AUTHZ_CREATE_ROLE,
  AUTHZ_DELETE_ROLE,
  AUTHZ_UPDATE_ROLE,
  AUTHZ_DELETE_USER_ROLE_OTHER,
  AUTHZ_CREATE_PERMISSION,
  AUTHZ_DELETE_PERMISSION,
  AUTHZ_MODIFY_USER_ROLE,
  AUTHZ_MODIFY_ROLE_PERMISSIONS,
)
