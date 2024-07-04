package errors

import (
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
