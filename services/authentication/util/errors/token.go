package errors

import "errors"

var ErrTokenDifferentUsage = errors.New("token provided has different usage")

var ErrTokenExpired = errors.New("token provided already reach the expiration time")

var ErrTokenNotFound = errors.New("token doesn't exist")
