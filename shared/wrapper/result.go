package wrapper

func Some[T any](data T, err error) Result[T] {
	return Result[T]{
		Data: data,
		Err:  err,
	}
}

func None[T any](err error) Result[T] {
	return Result[T]{
		Err: err,
	}
}

func NoneF[T any](none func() T, err error) Result[T] {
	return Result[T]{
		Data: none(),
		Err:  err,
	}
}

func SomeF[T any](f func() (T, error)) Result[T] {
	return Some(f())
}

func DropError[T any](val T, err error) T {
	return Some(val, err).Data
}

func SomeF1[T, P1 any](f func(P1) (T, error), param *P1) Result[T] {
	return Some(f(*param))
}

type Result[T any] struct {
	Data T
	Err  error
}

func (r Result[T]) Value() (T, error) {
	return r.Data, r.Err
}
