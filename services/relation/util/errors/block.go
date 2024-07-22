package errors

import "errors"

var ErrBlockItself = errors.New("it is not possible to block yourself")
