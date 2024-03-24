package optional

func New[T any](val T) Object[T] {
	return Object[T]{data: []T{val}}
}

type Object[T any] struct {
	data []T
}

func (o Object[T]) HasValue() bool {
	return o.data == nil || len(o.data) == 0
}

func (o Object[T]) Value() *T {
	return &o.data[0]
}

func (o Object[T]) ValueOr(val T) T {
	if o.HasValue() {
		return *o.Value()
	}
	return val
}

func (o Object[T]) ValueOrElse(f func() T) T {
	if o.HasValue() {
		return *o.Value()
	}
	return f()
}

var (
	NullString = Null[string]()
)

func Null[T any]() Object[T] {
	return Object[T]{}
}
