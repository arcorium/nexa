package util

func GenerateMultiple[T any](len int, f func() T) []T {
  result := make([]T, 0, len)
  for i := 0; i < len; i++ {
    result = append(result, f())
  }
  return result
}
