package status

import (
	"nexa/shared/optional"
)

type Code uint

const (
	// Success Code
	SUCCESS Code = iota
	CREATED
	UPDATED
	DELETED

	// Error Code
	// General
	INTERNAL_SERVER_ERROR
	REPOSITORY_ERROR
	// Other
	OBJECT_NOT_FOUND
)

var NullCode = optional.Null[Code]()
