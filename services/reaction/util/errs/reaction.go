package errs

import "errors"

var ErrResultWithDifferentLength = errors.New("result with different length")

var ErrItemNotFound = errors.New("expected item with those id is not found")
