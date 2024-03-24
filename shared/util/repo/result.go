package repo

import (
	"database/sql"
	"nexa/shared/wrapper"
)

// CheckResult check the result if the rows affected is 0 it will return sql.ErrNoRows
func CheckResult(result sql.Result, err error) error {
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// CheckSliceResult check the slice (from scanning ORMs) and err will become sql.ErrNoRows if the slice size is 0
func CheckSliceResult[T any](slice []T, err error) wrapper.Result[[]T] {
	if err != nil {
		return wrapper.None[[]T](err)
	}

	if len(slice) <= 0 {
		return wrapper.None[[]T](sql.ErrNoRows)
	}

	return wrapper.Some(slice, nil)
}

func CheckPaginationResult[T any](slice []T, count int, err error) wrapper.Result[[]T] {
	if err != nil {
		return wrapper.None[[]T](err)
	}

	if count == 0 || len(slice) == 0 {
		return wrapper.None[[]T](sql.ErrNoRows)
	}

	return wrapper.Some(slice, nil)
}
