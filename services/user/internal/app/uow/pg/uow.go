package pg

import (
  "context"
  "github.com/uptrace/bun"
  "nexa/services/user/internal/app/uow"
  "nexa/services/user/internal/infra/repository/pg"
  sharedUOW "nexa/shared/uow"
)

func NewUserUOW(db bun.IDB) sharedUOW.IUnitOfWork[uow.UserStorage] {
  storage := uow.NewStorage(pg.NewProfile(db), pg.NewUser(db))
  unit := &UserUOW{
    db:    db,
    cache: &storage,
  }

  return unit
}

type UserUOW struct {
  db    bun.IDB
  cache *uow.UserStorage
}

func (m *UserUOW) DoTx(ctx context.Context, f sharedUOW.UOWBlock[uow.UserStorage]) error {
  tx, err := m.db.BeginTx(ctx, nil)
  if err != nil {
    return err
  }

  err = f(ctx, m.repositories(tx))
  if err != nil {
    if err := tx.Rollback(); err != nil {
      return err
    }
    return err
  }

  return tx.Commit()
}

func (m *UserUOW) Repositories() uow.UserStorage {
  //storage := m.repositories(m.db)
  return *m.cache
}

func (m *UserUOW) repositories(db bun.IDB) uow.UserStorage {
  return uow.NewStorage(pg.NewProfile(db), pg.NewUser(db))
}
