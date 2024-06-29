package errors

import "errors"

var ErrRefreshTokenNotFound = errors.New("your session is expired, you need to re-authenticate")

var ErrDifferentScheme = errors.New("token type is not supported")

var ErrMalformedToken = errors.New("your token is malformed")

var ErrBadRelation = errors.New("refresh token is not associated with this access token")
