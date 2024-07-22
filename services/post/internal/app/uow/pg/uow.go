package pg

import (
  "context"
  sharedUOW "github.com/arcorium/nexa/shared/uow"
  "github.com/uptrace/bun"
  "nexa/services/post/internal/app/uow"
  "nexa/services/post/internal/infra/repository/pg"
)

func NewPostUOW(db bun.IDB) sharedUOW.IUnitOfWork[uow.PostStorage] {
  storage := uow.NewStorage(pg.NewPost(db))
  unit := &PostUOW{
    db:    db,
    cache: &storage,
  }

  return unit
}

type PostUOW struct {
  db    bun.IDB
  cache *uow.PostStorage
}

func (m *PostUOW) DoTx(ctx context.Context, f sharedUOW.UOWBlock[uow.PostStorage]) error {
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

func (m *PostUOW) Repositories() uow.PostStorage {
  //storage := m.repositories(m.db)
  return *m.cache
}

func (m *PostUOW) repositories(db bun.IDB) uow.PostStorage {
  return uow.NewStorage(pg.NewPost(db))
}
