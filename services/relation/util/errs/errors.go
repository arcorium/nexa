package errs

import "errors"

var ErrBlockItself = errors.New("it is not possible to block yourself")

var ErrFollowItself = errors.New("it is not possible to follow yourself")

var ErrUserNotFound = errors.New("user not found")
