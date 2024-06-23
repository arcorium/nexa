package util

func Nil[T any]() *T {
  return (*T)(nil)
}

func DoNothing(...any) {}

func CopyWith[T any](val T, f func(*T)) T {
  f(&val)
  return val
}

func CopyWithP[T any](val T, f func(*T)) *T {
  f(&val)
  return &val
}
