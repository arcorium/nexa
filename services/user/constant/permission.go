package constant

import authUtil "nexa/shared/util/auth"

// Actions
const (
  USER_GET             = "get"
  USER_UPDATE          = "update"
  USER_UPDATE_ARB      = "update:arb" // arb stands for arbitrary
  USER_DELETE          = "delete"
  USER_DELETE_ARB      = "delete:arb"
  USER_VERIFY_EMAIL    = "verif"
  USER_BANNED          = "ban"
  USER_GET_PROFILE_ARB = "get:profile:arb"
)

var USER_PERMISSIONS = authUtil.FullEncode(USER_SERVICE_RESOURCE,
  USER_DELETE_ARB,
  USER_BANNED,
  USER_UPDATE_ARB,
  USER_GET_PROFILE_ARB,
)
