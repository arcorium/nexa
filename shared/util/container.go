package util

import "nexa/shared/errors"

func CastSliceP[From, To any](slice []From, f func(*From) To) []To {
  if slice == nil || len(slice) == 0 {
    return nil
  }

  result := make([]To, 0, len(slice))
  for _, val := range slice {
    result = append(result, f(&val))
  }
  return result
}

func CastSlice[From, To any](slice []From, f func(From) To) []To {
  if slice == nil || len(slice) == 0 {
    return nil
  }

  result := make([]To, 0, len(slice))
  for _, val := range slice {
    result = append(result, f(val))
  }
  return result
}

func CastSliceErrP[From, To any](slice []From, f func(*From) (To, error)) ([]To, error) {
  if slice == nil || len(slice) == 0 {
    return nil, errors.ErrEmptySlice
  }

  result := make([]To, 0, len(slice))
  for _, val := range slice {
    res, err := f(&val)
    if err != nil {
      return nil, err
    }
    result = append(result, res)
  }
  return result, nil
}

// CastSliceErrsP cast slice with skipping error. when the function return error it will not stop like what CastSliceErrP do.
// instead, it will skip the element and continue to process next element. element that error when processed will be returned as the
// id along with the error object
func CastSliceErrsP[From, To any](slice []From, f func(*From) (To, error)) ([]To, errors.IndicesError) {
  if slice == nil || len(slice) == 0 {
    return nil, errors.IndicesError{}
  }

  result := make([]To, 0, len(slice))
  var errs []errors.IndexedError
  for i, val := range slice {
    res, err := f(&val)
    if err != nil {
      errs = append(errs, errors.IndexedError{Index: i, Err: err})
      continue
    }
    result = append(result, res)
  }

  return result, errors.IndicesError{Errs: errs}
}

// CastSliceErrs works like CastSliceErrsP, but it takes non-pointer
func CastSliceErrs[From, To any](slice []From, f func(From) (To, error)) ([]To, errors.IndicesError) {
  if slice == nil || len(slice) == 0 {
    return nil, errors.IndicesError{}
  }

  result := make([]To, 0, len(slice))
  var errs []errors.IndexedError
  for i, val := range slice {
    res, err := f(val)
    if err != nil {
      errs = append(errs, errors.IndexedError{Index: i, Err: err})
      continue
    }
    result = append(result, res)
  }

  return result, errors.IndicesError{Errs: errs}
}
