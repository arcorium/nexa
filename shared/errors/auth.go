package errors

import (
  "errors"
)

var ErrUnauthorizedPermission = errors.New("you dont have permission to access do this action")

var ErrUnauthorized = errors.New("you are not supposed to access this resource")
