package errors

import (
  "errors"
  "fmt"
  "strings"
)

type IndicesError struct {
  Errs []IndexedError
}

func (i IndicesError) Error() string {
  var str strings.Builder
  for _, err := range i.Errs {
    str.WriteString(err.Error() + "\n")
  }
  return str.String()
}

func (i IndicesError) IsNil() bool {
  return len(i.Errs) == 0
}

// IsEmptySlice check if the error is due to empty slice. Always check if it has error, otherwise
// it will panic due to indexing nil slice
func (i IndicesError) IsEmptySlice() bool {
  return i.Errs[0].Index == -1 || errors.Is(i.Errs[0].Err, ErrEmptySlice)
}

func (i IndicesError) ToGRPCError(fieldName string) error {
  return GrpcFieldIndexedErrors(fieldName, i)
}

func (i IndicesError) ToFieldError(fieldName string) FieldError {
  return NewFieldError(fieldName, i)
}

type IndexedError struct {
  Index int
  Err   error
}

func (i IndexedError) Error() string {
  return fmt.Sprintf("%d: %v", i.Index, i.Err)
}

func (i IndexedError) IsNil() bool {
  return i.Err == nil
}
