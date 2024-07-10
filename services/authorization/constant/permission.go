package constant

import authUtil "github.com/arcorium/nexa/shared/util/auth"

const (
  AUTHZ_CREATE_ROLE = "create:role"
  AUTHZ_DELETE_ROLE = "delete:role"
  AUTHZ_UPDATE_ROLE = "update:role"
  AUTHZ_READ_ROLE   = "read:role"

  AUTHZ_MODIFY_USER_ROLE        = "modify:role"
  AUTHZ_MODIFY_ROLE_PERMISSIONS = "modify:perms"
)

var AUTHZ_PERMISSIONS = authUtil.FullEncode(ROLE_RESOURCE,
  AUTHZ_CREATE_ROLE,
  AUTHZ_DELETE_ROLE,
  AUTHZ_UPDATE_ROLE,
  AUTHZ_READ_ROLE,
  AUTHZ_MODIFY_ROLE_PERMISSIONS,
  AUTHZ_MODIFY_USER_ROLE,
)

var PERM_PERMISSIONS = authUtil.FullEncode(PERMISSION_RESOURCE,
  AUTHZ_MODIFY_ROLE_PERMISSIONS,
)
