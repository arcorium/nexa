package main

import (
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/mailer/constant"
)

var MAILER_DEFAULT_PERMS = []types.Action{
  constant.MAIL_GET_TAG,
}

var MAILER_SUPER_PERMS = []types.Action{
  constant.MAIL_CREATE_TAG,
  constant.MAIL_DELETE_TAG,
  constant.MAIL_UPDATE_TAG,
  constant.MAIL_GET,
  constant.MAIL_DELETE,
  constant.MAIL_UPDATE,
}
