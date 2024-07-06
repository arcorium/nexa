package optional

func New[T any](val *T) Object[T] {
  if val == nil {
    return Null[T]()
  }
  return Some(*val)
}

func Some[T any](val T) Object[T] {
  return Object[T]{data: val, valid: true}
}

type Object[T any] struct {
  data  T
  valid bool
}

func (o Object[T]) HasValue() bool {
  return o.valid
}

func (o Object[T]) Value() *T {
  return &o.data
}

func (o Object[T]) ValueOr(val T) T {
  if o.HasValue() {
    return o.data
  }
  return val
}

func (o Object[T]) ValueOrElse(f func() T) T {
  if o.HasValue() {
    return o.data
  }
  return f()
}

func Null[T any]() Object[T] {
  return Object[T]{valid: false}
}
