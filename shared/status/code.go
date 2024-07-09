package status

import (
  "github.com/arcorium/nexa/shared/optional"
)

type Code uint

const (
  // Success Code
  SUCCESS_INTERNAL Code = iota // Should not be used for response
  SUCCESS
  CREATED
  UPDATED
  DELETED

  // Error Code
  // General
  INTERNAL_SERVER_ERROR
  EXTERNAL_SERVICE_ERROR
  REPOSITORY_ERROR
  BAD_REQUEST_ERROR
  FIELD_VALIDATION_ERROR
  NOT_AUTHORIZED_ERROR
  NOT_AUTHENTICATED_ERROR
  SERVICE_UNAVAILABLE_ERROR
  // Other
  OBJECT_NOT_FOUND
  OBJECT_ALREADY_EXIST
)

var NullCode = optional.Null[Code]()

// Success Code Helper
func Success() Object {
  return New(SUCCESS, nil)
}
func SuccessInternal() Object {
  return New(SUCCESS_INTERNAL, nil)
}
func Created() Object {
  return New(CREATED, nil)
}
func Updated() Object {
  return New(UPDATED, nil)
}
func Deleted() Object {
  return New(DELETED, nil)
}

// Error Code Helper
func ErrInternal(err error) Object {
  return New(INTERNAL_SERVER_ERROR, err)
}
func ErrExternal(err error) Object { return New(EXTERNAL_SERVICE_ERROR, err) }
func ErrUnAuthorized(err error) Object {
  return New(NOT_AUTHORIZED_ERROR, err)
}
func ErrUnAuthenticated(err error) Object {
  return New(NOT_AUTHENTICATED_ERROR, err)
}
func ErrBadRequest(err error) Object {
  return New(BAD_REQUEST_ERROR, err)
}

func ErrFieldValidation(err error) Object {
  return New(FIELD_VALIDATION_ERROR, err)
}
