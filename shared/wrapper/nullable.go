package wrapper

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"reflect"
)

func NewNullable[T any](data *T) Nullable[T] {
	return Nullable[T]{data}
}

type INullable[T any] interface {
	HasValue() bool
	Value() *T
	Value2() T
}

// Nullable *T wrapper, to be used for optional JSON, it works like sql.Null
type Nullable[T any] struct {
	data *T
}

func (n Nullable[T]) MarshalJSON() ([]byte, error) {
	if n.data == nil {
		return nil, nil
	}

	return json.Marshal(*n.data)
}

func (n Nullable[T]) UnmarshalJSON(bytes []byte) error {
	if len(bytes) == 0 {
		n.data = nil
		return nil
	}
	var t T
	err := json.Unmarshal(bytes, &t)
	if err != nil {
		n.data = &t
	}
	return err
}

func (n Nullable[T]) HasValue() bool {
	return n.data != nil
}

func (n Nullable[T]) Value() *T {
	return n.data
}

// Value2 Works like Value, but it will copy except for data type that has pointer as underlying, for example string, slice, map
func (n Nullable[T]) Value2() T {
	return *n.Value()
}

func nullableValidation[T any](sl reflect.Value) any {
	n, ok := sl.Interface().(Nullable[T])
	if !ok {
		return nil
	}
	return n.data
}

type (
	NullableString = Nullable[string]
	NullableInt    = Nullable[int]
	NullableInt8   = Nullable[int8]
	NullableInt16  = Nullable[int16]
	NullableInt32  = Nullable[int32]
	NullableInt64  = Nullable[int64]
	NullableUInt   = Nullable[uint]
	NullableUInt8  = Nullable[uint8]
	NullableUInt16 = Nullable[uint16]
	NullableUInt32 = Nullable[uint32]
	NullableUInt64 = Nullable[uint64]
)

func RegisterValidation[T any](validate *validator.Validate) {
	validate.RegisterCustomTypeFunc(nullableValidation[T], Nullable[T]{})
}

func RegisterDefaultNullableValidations(validate *validator.Validate) {
	RegisterValidation[string](validate)
	RegisterValidation[int](validate)
	RegisterValidation[int8](validate)
	RegisterValidation[int16](validate)
	RegisterValidation[int32](validate)
	RegisterValidation[int64](validate)
	RegisterValidation[uint](validate)
	RegisterValidation[uint8](validate)
	RegisterValidation[uint16](validate)
	RegisterValidation[uint32](validate)
	RegisterValidation[uint64](validate)
}

// SetOnNonNull set value on dest if only the nullable object has value
func SetOnNonNull[T any, U INullable[T]](dest *T, nullable U) {
	if nullable.HasValue() {
		*dest = nullable.Value2()
	}
}

// SetOnNonNullCasted set value on dest if only the nullable object has value. Nullable object can have different type from dest parameter.
// castFunc will be used to do casting from nullable type to destination type
func SetOnNonNullCasted[T, V any, U INullable[V]](dest *T, nullable U, castFunc func(V) T) {
	if nullable.HasValue() {
		*dest = castFunc(nullable.Value2())
	}
}
