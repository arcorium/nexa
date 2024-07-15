package constant

import (
  "github.com/arcorium/nexa/shared/types"
  "github.com/arcorium/nexa/shared/util/auth"
)

// Default
const MAIL_GET_TAG types.Action = "get:tag"

// Super
const (
  MAIL_CREATE_TAG types.Action = "create:tag"
  MAIL_DELETE_TAG types.Action = "delete:tag"
  MAIL_UPDATE_TAG types.Action = "update:tag"
  MAIL_GET        types.Action = "get:mail"
  MAIL_DELETE     types.Action = "delete:mail"
  MAIL_UPDATE     types.Action = "update:mail"
)

var MAILER_PERMISSIONS = auth.FullEncode(SERVICE_RESOURCE,
  MAIL_GET,
  MAIL_UPDATE,
  MAIL_DELETE,
  MAIL_GET_TAG,
  MAIL_CREATE_TAG,
  MAIL_DELETE_TAG,
  MAIL_UPDATE_TAG,
)
