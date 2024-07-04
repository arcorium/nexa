package repo

type QueryParameter struct {
  Offset uint64
  Limit  uint64
}

func NewPaginatedResult[T any](result []T, totalCount uint64) PaginatedResult[T] {
  return PaginatedResult[T]{
    Data:    result,
    Total:   totalCount,
    Element: uint64(len(result)),
  }
}

type PaginatedResult[T any] struct {
  Data    []T
  Total   uint64
  Element uint64
}

func (p *PaginatedResult[T]) HasValue() bool {
  return p.Total != 0 && p.Element != 0
}
