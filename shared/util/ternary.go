package util

func Ternary[T any](cond bool, trueVal, falseVal T) T {
  if cond {
    return trueVal
  }
  return falseVal
}

// ReturnOnEqual Simple ternary that will return retVal parameter if only expected is equal to comparator parameter,
// it will return the expected parameter otherwise
func ReturnOnEqual[T comparable](expected, comparator, retVal T) T {
  if expected == comparator {
    return retVal
  }
  return expected
}

func TernaryF[T any](cond bool, trueFunc func() T, falseFunc func() T) T {
  if cond {
    return trueFunc()
  }
  return falseFunc()
}

// NilOr return nil when obj is nil, otherwise it will call the function and use the function return value
func NilOr[R, T any](obj *T, f func(obj *T) *R) *R {
  if obj == nil {
    return nil
  }

  return f(obj)
}
