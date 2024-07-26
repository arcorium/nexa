package repo

import (
  "database/sql"
  "errors"
  "github.com/arcorium/nexa/shared/types"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "github.com/uptrace/bun/driver/pgdriver"
  "go.opentelemetry.io/otel/trace"
)

func checkPgDriverViolation(err error) error {
  pg, ok := err.(pgdriver.Error)
  if !ok {
    return err
  }
  switch pg.Field('C') {
  case "23000", "23502", "23503", "23505", "23514", "23P01":
    return errors.Join(ErrAlreadyExists, err)
  case "23001": // Null on not null fields
    return err
  default:
    return err
  }
}

// CheckResult check the result if the rows affected is 0 it will return sql.ErrNoRows
func CheckResult(result sql.Result, err error) error {
  if err != nil {
    return checkPgDriverViolation(err)
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

// CheckResultWithSpan works like CheckResult, but it will record error for span
func CheckResultWithSpan(result sql.Result, err error, span trace.Span) error {
  resultErr := CheckResult(result, err)
  if resultErr != nil {
    spanUtil.RecordError(resultErr, span)
  }
  return resultErr
}

// CheckSliceResult check the slice (from scanning ORMs) and err will become sql.ErrNoRows if the slice size is 0
func CheckSliceResult[T any](slice []T, err error) types.Result[[]T] {
  if err != nil {
    return types.None[[]T](checkPgDriverViolation(err))
  }

  if len(slice) <= 0 {
    return types.None[[]T](sql.ErrNoRows)
  }

  return types.Some(slice, nil)
}

// CheckSliceResultWithSpan works like CheckSliceResult, but it will record error for span
func CheckSliceResultWithSpan[T any](slice []T, err error, span trace.Span) types.Result[[]T] {
  result := CheckSliceResult(slice, err)
  if result.IsError() {
    spanUtil.RecordError(result.Err, span)
  }
  return result
}

func CheckPaginationResult[T any](slice []T, count int, err error) types.Result[[]T] {
  if err != nil {
    return types.None[[]T](checkPgDriverViolation(err))
  }

  if count == 0 || len(slice) == 0 {
    return types.None[[]T](sql.ErrNoRows)
  }

  return types.Some(slice, nil)
}

func CheckPaginationResultWithSpan[T any](slice []T, count int, err error, span trace.Span) types.Result[[]T] {
  result := CheckPaginationResult(slice, count, err)
  if result.IsError() {
    spanUtil.RecordError(result.Err, span)
  }
  return result
}
