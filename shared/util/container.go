package util

import (
  sharedErr "nexa/shared/errors"
)

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

func MapToSlice[KFrom comparable, VFrom, To any](maps map[KFrom]VFrom, f func(KFrom, VFrom) To) []To {
  if maps == nil || len(maps) == 0 {
    return nil
  }

  result := make([]To, 0, len(maps))
  for key, val := range maps {
    result = append(result, f(key, val))
  }
  return result
}

func CastSliceErrP[From, To any](slice []From, f func(*From) (To, error)) ([]To, error) {
  if slice == nil || len(slice) == 0 {
    return nil, sharedErr.ErrEmptySlice
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
// instead, it will skip the element and continue to process next element. element with an error will be returned as the
// index along with the error object and not appended on returned slice
func CastSliceErrsP[From, To any](slice []From, f func(*From) (To, error)) ([]To, sharedErr.IndicesError) {
  if slice == nil || len(slice) == 0 {
    return nil, sharedErr.IndicesError{Errs: []sharedErr.IndexedError{{-1, sharedErr.ErrEmptySlice}}}
  }

  result := make([]To, 0, len(slice))
  var errs []sharedErr.IndexedError
  for i, val := range slice {
    res, err := f(&val)
    if err != nil {
      errs = append(errs, sharedErr.IndexedError{Index: i, Err: err})
      continue
    }
    result = append(result, res)
  }

  return result, sharedErr.IndicesError{Errs: errs}
}

// CastSliceErrs works like CastSliceErrsP, but it takes non-pointer
func CastSliceErrs[From, To any](slice []From, f func(From) (To, error)) ([]To, sharedErr.IndicesError) {
  if slice == nil || len(slice) == 0 {
    return nil, sharedErr.IndicesError{Errs: []sharedErr.IndexedError{{-1, sharedErr.ErrEmptySlice}}}
  }

  result := make([]To, 0, len(slice))
  var errs []sharedErr.IndexedError
  for i, val := range slice {
    res, err := f(val)
    if err != nil {
      errs = append(errs, sharedErr.IndexedError{Index: i, Err: err})
      continue
    }
    result = append(result, res)
  }

  return result, sharedErr.IndicesError{Errs: errs}
}
