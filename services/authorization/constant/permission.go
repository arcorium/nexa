package constant

import (
  "github.com/arcorium/nexa/shared/types"
  authUtil "github.com/arcorium/nexa/shared/util/auth"
)

// Default
const (
  AUTHZ_GET_ROLE         types.Action = "get:role"
  AUTHZ_DELETE_USER_ROLE types.Action = "delete:user:role"
)

// Super
const (
  AUTHZ_CREATE_ROLE             types.Action = "create:role"
  AUTHZ_DELETE_ROLE             types.Action = "delete:role"
  AUTHZ_UPDATE_ROLE             types.Action = "update:role"
  AUTHZ_ADD_USER_ROLE           types.Action = "add:user:role"
  AUTHZ_DELETE_USER_ROLE_ARB    types.Action = "delete:user:role:arb"
  AUTHZ_CREATE_PERMISSION       types.Action = "create:perms"
  AUTHZ_DELETE_PERMISSION       types.Action = "delete:perms"
  AUTHZ_MODIFY_ROLE_PERMISSIONS types.Action = "modify:perms"
)

var AUTHZ_PERMISSIONS = authUtil.FullEncode(SERVICE_RESOURCE,
  AUTHZ_GET_ROLE,
  AUTHZ_DELETE_USER_ROLE,
  AUTHZ_CREATE_ROLE,
  AUTHZ_DELETE_ROLE,
  AUTHZ_UPDATE_ROLE,
  AUTHZ_DELETE_USER_ROLE_ARB,
  AUTHZ_CREATE_PERMISSION,
  AUTHZ_DELETE_PERMISSION,
  AUTHZ_ADD_USER_ROLE,
  AUTHZ_MODIFY_ROLE_PERMISSIONS,
)
