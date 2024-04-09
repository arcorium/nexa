package errors

import "errors"

var ErrRefreshTokenExpired = errors.New("your session is expired, you need to re-authenticate")
