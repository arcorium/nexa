package pg

import (
  "context"
  sharedUOW "github.com/arcorium/nexa/shared/uow"
  "github.com/uptrace/bun"
  "nexa/services/mailer/internal/app/uow"
  "nexa/services/mailer/internal/infra/repository/pg"
)

func NewMailUOW(db bun.IDB) sharedUOW.IUnitOfWork[uow.MailStorage] {
  storage := uow.NewStorage(pg.NewMail(db), pg.NewTag(db))
  unit := &MailUOW{
    db:    db,
    cache: &storage,
  }

  return unit
}

type MailUOW struct {
  db    bun.IDB
  cache *uow.MailStorage
}

func (m *MailUOW) DoTx(ctx context.Context, f sharedUOW.UOWBlock[uow.MailStorage]) error {
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

func (m *MailUOW) Repositories() uow.MailStorage {
  //storage := m.repositories(m.db)
  return *m.cache
}

func (m *MailUOW) repositories(db bun.IDB) uow.MailStorage {
  return uow.NewStorage(pg.NewMail(db), pg.NewTag(db))
}
