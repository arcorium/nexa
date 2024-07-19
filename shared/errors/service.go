package errors

import "errors"

var ErrServiceUnavailable = errors.New("service currently unavailable, try again later")

var ErrServiceRecovering = errors.New("service currently recovering, try again later")
