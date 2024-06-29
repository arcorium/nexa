package types

func Null[T any]() T {
  var t T
  return t
}

func Nil[T any]() *T {
  return (*T)(nil)
}
