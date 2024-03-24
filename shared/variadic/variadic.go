package variadic

import "errors"

func New[T any](values ...T) Object[T] {
	return Object[T]{data: values}
}

type Object[T any] struct {
	data []T
}

// HasOne check if variadic only contains single value
func (o Object[T]) HasOne() bool {
	return o.data != nil && len(o.data) == 1
}

// HasValue check if variadic has values
func (o Object[T]) HasValue() bool {
	return o.data != nil && len(o.data) > 0
}

// HasAtLeast check if variadic size is at minimum the same with the parameter
func (o Object[T]) HasAtLeast(min int) bool {
	return o.data != nil && len(o.data) >= min
}

// First get pointer to first object on variadic if only variadic has value
func (o Object[T]) First() (*T, error) {
	if !o.HasValue() {
		return nil, ErrObjectNotFound
	}
	return &o.data[0], nil
}

// At get pointer to object on variadic if only the variadic has at least size from parameter
func (o Object[T]) At(index uint) (*T, error) {
	if !o.HasAtLeast(int(index)) {
		return nil, ErrObjectNotFound
	}
	return &o.data[index], nil
}

func (o Object[T]) Values() []T {
	return o.data
}

// DoAtFirst call f function if only the variadic has at least single object and f function will not be called when
// the variadic is empty
func (o Object[T]) DoAtFirst(f func(*T)) {
	first, err := o.First()
	if err != nil {
		return
	}
	f(first)
}

var ErrObjectNotFound = errors.New("variadic object not found")
