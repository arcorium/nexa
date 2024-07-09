package status

import (
  "database/sql"
  "errors"
  "nexa/shared/optional"
  "nexa/shared/types"
  "nexa/shared/variadic"
)

// FromRepository override only the code when the error is sql.ErrNoRows and use the sql.ErrNoRows as the error
func FromRepository(err error, notFoundCode optional.Object[Code]) Object {
  if errors.Is(err, sql.ErrNoRows) {
    return New(notFoundCode.ValueOr(OBJECT_NOT_FOUND), err)
  }
  return Error(REPOSITORY_ERROR, err)
}

// FromRepositoryOverride override the code and error when the error is sql.ErrNoRows
func FromRepositoryOverride(err error, notFoundOver ...types.Pair[Code, error]) Object {
  va := variadic.New(notFoundOver...)
  if errors.Is(err, sql.ErrNoRows) && va.HasValue() {
    if len(va.Values()) > 1 {
      panic("Too many values, expected only single object")
    }

    over, _ := va.First()
    return New(over.First, over.Second)
  }
  return Error(REPOSITORY_ERROR, err)
}

// FromRepositoryExist helper function to call FromRepositoryOverride which handle sql.ErrNoRows as
// object already exists. Used for inserting new data
func FromRepositoryExist(err error) Object {
  return FromRepositoryOverride(err, types.NewPair(OBJECT_ALREADY_EXIST, errors.New("object already exist")))
}
