package errs

import "errors"

var ErrResultWithDifferentLength = errors.New("result has different length from expected one")

var ErrPostNotFound = errors.New("post not found")
