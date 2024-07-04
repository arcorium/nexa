package util

// ReturnOnEqual Simple ternary that will return retVal parameter if only expected is equal to comparator parameter,
// it will return the expected parameter otherwise
func ReturnOnEqual[T comparable](expected, comparator, retVal T) T {
  if expected == comparator {
    return retVal
  }
  return expected
}

func Ternary[T any](condition bool, trueVal, falseVal T) T {
  if condition {
    return trueVal
  }
  return falseVal
}
