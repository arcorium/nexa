package constant

import authUtil "github.com/arcorium/nexa/shared/util/auth"

// Actions
// Format: {resource}:{action}:{extended}
const (
  USER_GET          = "get"
  USER_UPDATE       = "update"
  USER_UPDATE_ARB   = "update:arb"
  USER_DELETE       = "delete"
  USER_DELETE_ARB   = "delete:arb"
  USER_VERIFY_EMAIL = "verif"
  USER_BANNED       = "ban"

  // NOTE: Used when the user could have multiple profiles
  PROFILE_GET     = "get"
  PROFILE_GET_ARB = "get:arb"
)

var USER_PERMISSIONS = authUtil.FullEncode(USER_SERVICE_RESOURCE,
  USER_GET,
  USER_UPDATE,
  USER_UPDATE_ARB,
  USER_DELETE,
  USER_DELETE_ARB,
  USER_VERIFY_EMAIL,
  USER_BANNED,
)

var PROFILE_PERMISSIONS = authUtil.FullEncode(PROFILE_SERVICE_RESOURCE,
  PROFILE_GET,
  PROFILE_GET_ARB,
)
