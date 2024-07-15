package main

import (
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/authentication/constant"
)

var AUTHN_DEFAULT_PERMS = []types.Action{
  constant.AUTHN_GET_USER,
  constant.AUTHN_GET_PROFILE,
  constant.AUTHN_GET_CREDENTIAL,
  constant.AUTHN_UPDATE_USER,
  constant.AUTHN_DELETE_USER,
  constant.AUTHN_LOGOUT_USER,
  constant.AUTHN_VERIFY_EMAIL_USER,
}

var AUTHN_SUPER_PERMS = []types.Action{
  constant.AUTHN_GET_PROFILE_ARB,
  constant.AUTHN_GET_CREDENTIAL_ARB,
  constant.AUTHN_UPDATE_USER_ARB,
  constant.AUTHN_BANNED,
  constant.AUTHN_LOGOUT_USER_ARB,
  constant.AUTHN_DELETE_USER_ARB,
}
