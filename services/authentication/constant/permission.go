package constant

import authUtil "nexa/shared/util/auth"

const (
  AUTHN_GET_OTHER_CREDENTIALS = "read:cred:other"
  AUTHN_LOGOUT_OTHER          = "delete:cred:other"
)

var AUTHN_PERMISSIONS = authUtil.FullEncode(SERVICE_RESOURCE,
  AUTHN_GET_OTHER_CREDENTIALS,
  AUTHN_LOGOUT_OTHER,
)
