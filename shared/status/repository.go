package status

import (
  "database/sql"
  "errors"
  "github.com/arcorium/nexa/shared/optional"
  "github.com/arcorium/nexa/shared/types"
  "github.com/arcorium/nexa/shared/util/repo"
  "github.com/arcorium/nexa/shared/variadic"
)

// FromRepository override only the code when the error is sql.ErrNoRows and use the sql.ErrNoRows as the error
func FromRepository(err error, notFoundCode optional.Object[Code]) Object {
  if errors.Is(err, sql.ErrNoRows) {
    return ErrNotFound()
  }
  return Error(REPOSITORY_ERROR, err)
}

func FromRepository2(err error, notFound optional.Object[Object], exists optional.Object[Object]) Object {
  if errors.Is(err, sql.ErrNoRows) {
    return notFound.ValueOr(ErrNotFound())
  } else if errors.Is(err, repo.ErrAlreadyExists) {
    return exists.ValueOr(ErrAlreadyExist())
  }
  return Error(REPOSITORY_ERROR, err)
}

// FromRepositoryOverride override the code and error when the error is sql.ErrNoRows
func FromRepositoryOverride(err error, notFoundOver ...types.Pair[Code, error]) Object {
  va := variadic.New(notFoundOver...)
  if errors.Is(err, sql.ErrNoRows) && va.HasValue() {
    over, _ := va.First()
    return New(over.First, over.Second)
  }
  return Error(REPOSITORY_ERROR, err)
}

// FromRepositoryOverrideObject override the code and error when the error is sql.ErrNoRows
func FromRepositoryOverrideObject(err error, notFoundOver ...Object) Object {
  va := variadic.New(notFoundOver...)
  if errors.Is(err, sql.ErrNoRows) && va.HasValue() {
    over, _ := va.First()
    return *over
  }
  return Error(REPOSITORY_ERROR, err)
}

// FromRepositoryExist helper function to call FromRepositoryOverride which handle sql.ErrNoRows as
// object already exists. Used for inserting new data
func FromRepositoryExist(err error) Object {
  return FromRepositoryOverrideObject(err, ErrAlreadyExist())
}
