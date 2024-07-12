package constant

import authUtil "github.com/arcorium/nexa/shared/util/auth"

const (
  AUTHZ_CREATE_ROLE = "create:role"
  AUTHZ_DELETE_ROLE = "delete:role"
  AUTHZ_UPDATE_ROLE = "update:role"
  AUTHZ_READ_ROLE   = "read:role"

  AUTHZ_CREATE_PERMISSION = "create:perms"
  AUTHZ_DELETE_PERMISSION = "delete:perms"
  AUTHZ_READ_PERMISSION   = "read:perms"

  AUTHZ_MODIFY_USER_ROLE        = "modify:role"
  AUTHZ_MODIFY_ROLE_PERMISSIONS = "modify:perms"
)

var AUTHZ_PERMISSIONS = authUtil.FullEncode(SERVICE_RESOURCE,
  AUTHZ_CREATE_ROLE,
  AUTHZ_DELETE_ROLE,
  AUTHZ_UPDATE_ROLE,
  AUTHZ_READ_ROLE,
  AUTHZ_CREATE_PERMISSION,
  AUTHZ_DELETE_PERMISSION,
  AUTHZ_READ_PERMISSION,
  AUTHZ_MODIFY_USER_ROLE,
  AUTHZ_MODIFY_ROLE_PERMISSIONS,
)
