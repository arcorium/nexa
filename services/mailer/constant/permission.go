package constant

import "nexa/shared/util/auth"

const (
  MAIL_CREATE_TAG = "create:tag"
  MAIL_DELETE_TAG = "delete:tag"
  MAIL_UPDATE_TAG = "update:tag"
  MAIL_READ_TAG   = "read:tag"
  MAIL_READ       = "read"
  MAIL_UPDATE     = "update"
  MAIL_DELETE     = "delete"
)

var MAILER_PERMISSIONS = map[string]string{
  MAIL_READ:       auth.Encode(SERVICE_RESOURCE, MAIL_READ),
  MAIL_UPDATE:     auth.Encode(SERVICE_RESOURCE, MAIL_UPDATE),
  MAIL_DELETE:     auth.Encode(SERVICE_RESOURCE, MAIL_DELETE),
  MAIL_CREATE_TAG: auth.Encode(SERVICE_RESOURCE, MAIL_CREATE_TAG),
  MAIL_DELETE_TAG: auth.Encode(SERVICE_RESOURCE, MAIL_DELETE_TAG),
  MAIL_READ_TAG:   auth.Encode(SERVICE_RESOURCE, MAIL_READ_TAG),
  MAIL_UPDATE_TAG: auth.Encode(SERVICE_RESOURCE, MAIL_UPDATE_TAG),
}
