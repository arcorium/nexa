package errors

import "errors"

// Validation
var ErrFieldEmpty = errors.New("field should not be empty")
var ErrZeroEmpty = errors.New("field should be non-zero")

// Container
var ErrEmptySlice = errors.New("slice is empty")
