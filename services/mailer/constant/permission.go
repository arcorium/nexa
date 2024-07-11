package constant

import "github.com/arcorium/nexa/shared/util/auth"

const (
  MAIL_CREATE_TAG = "create:tag"
  MAIL_DELETE_TAG = "delete:tag"
  MAIL_UPDATE_TAG = "update:tag"
  MAIL_READ       = "read"
  MAIL_UPDATE     = "update"
  MAIL_DELETE     = "delete"
)

var MAILER_PERMISSIONS = auth.FullEncode(SERVICE_RESOURCE,
  MAIL_READ,
  MAIL_UPDATE,
  MAIL_DELETE,
  MAIL_CREATE_TAG,
  MAIL_DELETE_TAG,
  MAIL_UPDATE_TAG,
)
