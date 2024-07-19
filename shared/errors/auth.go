package errors

import (
  "errors"
)

var ErrUnauthenticated = errors.New("you need to provide your token by login")

var ErrUnauthorizedPermission = errors.New("you dont have permission to do this action")

var ErrUnauthorized = errors.New("you are not supposed to access this resource")
