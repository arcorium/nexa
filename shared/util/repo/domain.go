package repo

type IDataAccessModel[T any] interface {
	ToDomain() T
}
