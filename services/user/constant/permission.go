package constant

import authUtil "nexa/shared/util/auth"

// Actions
const (
  USER_DELETE       = "delete"
  USER_DELETE_OTHER = "delete:other"
  USER_BANNED       = "ban"
  USER_UPDATE       = "update"
  USER_UPDATE_OTHER = "update:other"
  USER_GET          = "get"

  USER_GET_PROFILE       = "get:profile"
  USER_GET_PROFILE_OTHER = "get:profile:other"
)

var USER_PERMISSIONS = authUtil.FullEncode(SERVICE_RESOURCE,
  USER_DELETE,
  USER_DELETE_OTHER,
  USER_BANNED,
  USER_UPDATE,
  USER_UPDATE_OTHER,
  USER_GET,
  USER_GET_PROFILE,
  USER_GET_PROFILE_OTHER,
)
