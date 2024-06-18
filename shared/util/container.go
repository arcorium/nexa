package util

import "nexa/shared/errors"

func CastSlice[From, To any](slice []From, f func(*From) To) []To {
  if slice == nil || len(slice) == 0 {
    return nil
  }

  result := make([]To, 0, len(slice))
  for _, val := range slice {
    result = append(result, f(&val))
  }
  return result
}

func CastSlice2[From, To any](slice []From, f func(From) To) []To {
  if slice == nil || len(slice) == 0 {
    return nil
  }

  result := make([]To, 0, len(slice))
  for _, val := range slice {
    result = append(result, f(val))
  }
  return result
}

func CastSliceErr[From, To any](slice []From, f func(*From) (To, error)) ([]To, error) {
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

// CastSliceErrs cast slice with skipping error. when the function return error it will not stop like what CastSliceErr do.
// instead, it will skip the element and continue to process next element. element that error when processed will be returned as the
// id along with the error object
func CastSliceErrs[From, To any](slice []From, f func(*From) (To, error)) ([]To, []errors.IndexedError) {
  if slice == nil || len(slice) == 0 {
    return nil, []errors.IndexedError{{Index: -1, Err: errors.ErrEmptySlice}}
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

  return result, errs
}

// CastSliceErrs2 works like CastSliceErrs, but it takes non-pointer
func CastSliceErrs2[From, To any](slice []From, f func(From) (To, error)) ([]To, []errors.IndexedError) {
  if slice == nil || len(slice) == 0 {
    return nil, []errors.IndexedError{{Index: -1, Err: errors.ErrEmptySlice}}
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

  return result, errs
}
