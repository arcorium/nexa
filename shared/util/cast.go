package util

import (
  "fmt"
  "nexa/shared/types"
)

func ToUnderlyingEnum[T types.Enum[U], U any](enum T) U {
  return enum.Underlying()
}

func ToString[T fmt.Stringer](str T) string {
  return str.String()
}
