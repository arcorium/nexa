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
  return i.Errs == nil
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
