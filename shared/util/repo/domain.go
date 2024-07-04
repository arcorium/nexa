package repo

type IDataAccessModel[T any] interface {
  ToDomain() T
}

type IDataAccessModelWithError[T any] interface {
  ToDomain() (T, error)
}

func ToDomainErr[T IDataAccessModelWithError[U], U any](data T) (U, error) {
  return data.ToDomain()
}

func ToDomain[T IDataAccessModel[U], U any](data T) U {
  return data.ToDomain()
}
