package constant

import (
  "github.com/arcorium/nexa/shared/types"
  authUtil "github.com/arcorium/nexa/shared/util/auth"
)

// Default
const (
  AUTHN_GET_USER       types.Action = "get:user"
  AUTHN_GET_CREDENTIAL types.Action = "get:cred"
  AUTHN_UPDATE_USER    types.Action = "update:user"
  AUTHN_DELETE_USER    types.Action = "delete:user"
  AUTHN_LOGOUT_USER    types.Action = "logout:user"
  AUTHN_VERIF_REQUEST  types.Action = "req:user:verif"
)

// Super
const (
  AUTHN_CREATE_USER        types.Action = "create:user"
  AUTHN_GET_CREDENTIAL_ARB types.Action = "get:cred:arb"
  AUTHN_UPDATE_USER_ARB    types.Action = "update:user:arb"
  AUTHN_BANNED             types.Action = "ban:user"
  AUTHN_LOGOUT_USER_ARB    types.Action = "logout:user:arb"
  AUTHN_DELETE_USER_ARB    types.Action = "delete:user:arb"
)

var AUTHN_PERMISSIONS = authUtil.FullEncode(SERVICE_RESOURCE,
  AUTHN_GET_USER,
  AUTHN_GET_CREDENTIAL,
  AUTHN_UPDATE_USER,
  AUTHN_DELETE_USER,
  AUTHN_LOGOUT_USER,
  AUTHN_VERIF_REQUEST,
  AUTHN_GET_CREDENTIAL_ARB,
  AUTHN_CREATE_USER,
  AUTHN_UPDATE_USER_ARB,
  AUTHN_BANNED,
  AUTHN_LOGOUT_USER_ARB,
  AUTHN_DELETE_USER_ARB,
)
