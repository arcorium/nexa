package main

import (
  authUtil "github.com/arcorium/nexa/shared/util/auth"
  "nexa/services/authorization/constant"
)

var AUTHZ_DEFAULT_PERMS = authUtil.FullEncode(constant.SERVICE_RESOURCE,
  constant.AUTHZ_READ_ROLE,
  constant.AUTHZ_READ_PERMISSION,
  constant.AUTHZ_DELETE_USER_ROLE,
)

var AUTHZ_SUPER_PERMS = authUtil.FullEncode(constant.SERVICE_RESOURCE,
  constant.AUTHZ_CREATE_ROLE,
  constant.AUTHZ_DELETE_ROLE,
  constant.AUTHZ_UPDATE_ROLE,
  constant.AUTHZ_DELETE_USER_ROLE_OTHER,
  constant.AUTHZ_CREATE_PERMISSION,
  constant.AUTHZ_DELETE_PERMISSION,
  constant.AUTHZ_MODIFY_USER_ROLE,
  constant.AUTHZ_MODIFY_ROLE_PERMISSIONS,
)
