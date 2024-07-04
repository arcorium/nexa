package types

type Field[T any] struct {
  Name string
  Val  T
}

func NewField[T any](name string, val T) Field[T] {
  return Field[T]{
    Name: name,
    Val:  val,
  }
}
