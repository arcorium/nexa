package errors

import "errors"

var ErrTokenDifferentUsage = errors.New("token provided has different usage")

var ErrTokenDifferentUser = errors.New("token provided does not for you")

var ErrTokenExpired = errors.New("token provided already reach the expiration time")

var ErrTokenNotFound = errors.New("token doesn't exist")

var ErrTokenUsageUnknown = errors.New("token usage unknown")

//var ErrTokenUsageAlreadyExist = errors.New("user already has this token")
