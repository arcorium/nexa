package uow

import (
  "context"
)

type UOWBlock[T any] func(context.Context, T) error

// IUnitOfWork interface unit of work for all services. T should be storage for repositories
type IUnitOfWork[T any] interface {
  // DoTx create transaction and run function f. Repository storage on f should be in transaction
  // error returned is error forwarded from UOWBlock
  DoTx(ctx context.Context, f UOWBlock[T]) error
  // Repositories return T type of repository storage that is not in transaction
  Repositories() T
}
