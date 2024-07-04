package types

func Null[T any]() T {
  var t T
  return t
}

func Nil[T any]() *T {
  return (*T)(nil)
}

// NilOrElse return nil when obj is nil, otherwise it will call the function and use the function return value
func NilOrElse[R, T any](obj *T, f func(obj *T) *R) *R {
  if obj == nil {
    return nil
  }

  return f(obj)
}

func NilOrElseErr[R, T any](obj *T, f func(obj *T) (*R, error)) (*R, error) {
  if obj == nil {
    return nil, nil
  }

  return f(obj)
}
