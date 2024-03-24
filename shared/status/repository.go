package status

import (
	"database/sql"
	"errors"
	"nexa/shared/optional"
)

func FromRepository(err error, notFoundCode optional.Object[Code]) Object {
	if errors.Is(err, sql.ErrNoRows) {
		return New(notFoundCode.ValueOr(OBJECT_NOT_FOUND), err)
	}
	return Error(REPOSITORY_ERROR, err)
}
