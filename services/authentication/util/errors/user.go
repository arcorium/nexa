package errors

import "errors"

var ErrResetPasswordWithoutUserId = errors.New("expected user id, when there are no token provided")
