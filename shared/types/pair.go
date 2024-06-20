package types

func NewPair[T, U any](first T, second U) Pair[T, U] {
  return Pair[T, U]{
    First:  first,
    Second: second,
  }
}

type Pair[T, U any] struct {
  First  T
  Second U
}
