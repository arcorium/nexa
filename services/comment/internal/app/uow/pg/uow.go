package pg

import (
  "context"
  sharedUOW "github.com/arcorium/nexa/shared/uow"
  "github.com/uptrace/bun"
  "nexa/services/comment/internal/app/uow"
  "nexa/services/comment/internal/infra/repository/pg"
)

func NewCommentUOW(db bun.IDB) sharedUOW.IUnitOfWork[uow.CommentStorage] {
  storage := uow.NewStorage(pg.NewComment(db))
  unit := &PostUOW{
    db:    db,
    cache: &storage,
  }

  return unit
}

type PostUOW struct {
  db    bun.IDB
  cache *uow.CommentStorage
}

func (m *PostUOW) DoTx(ctx context.Context, f sharedUOW.UOWBlock[uow.CommentStorage]) error {
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

func (m *PostUOW) Repositories() uow.CommentStorage {
  //storage := m.repositories(m.db)
  return *m.cache
}

func (m *PostUOW) repositories(db bun.IDB) uow.CommentStorage {
  return uow.NewStorage(pg.NewComment(db))
}
