package entity

import "context"

type IToEntity[T any] interface {
	ToEntity() T
}

type IToEntityCtx[T any] interface {
	ToEntity(context.Context) T
}

func MapToEntityFunc[T IToEntity[U], U any]() func(T) U {
	return func(t T) U {
		return t.ToEntity()
	}
}
