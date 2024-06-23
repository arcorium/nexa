package types

type Enum[T any] interface {
  Underlying() T
}
