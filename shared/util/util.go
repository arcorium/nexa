package util

func Nil[T any]() *T {
	return (*T)(nil)
}
