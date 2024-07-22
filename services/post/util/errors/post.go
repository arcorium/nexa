package errors

import "errors"

var ErrPostWithNoVersion = errors.New("post malformed, it has no versions")
