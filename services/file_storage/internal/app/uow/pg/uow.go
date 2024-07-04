package pg

import (
  "context"
  "github.com/uptrace/bun"
  "nexa/services/file_storage/internal/app/uow"
  "nexa/services/file_storage/internal/infra/repository/pg"
  sharedUOW "nexa/shared/uow"
)

func NewStorageUOW(db bun.IDB) sharedUOW.IUnitOfWork[uow.FileMetadataStorage] {
  storage := uow.NewStorage(pg.NewFileMetadataRepository(db))
  unit := &FileStorageUOW{
    db:    db,
    cache: &storage,
  }

  return unit
}

type FileStorageUOW struct {
  db    bun.IDB
  cache *uow.FileMetadataStorage
}

func (m *FileStorageUOW) DoTx(ctx context.Context, f sharedUOW.UOWBlock[uow.FileMetadataStorage]) error {
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

func (m *FileStorageUOW) Repositories() uow.FileMetadataStorage {
  //storage := m.repositories(m.db)
  return *m.cache
}

func (m *FileStorageUOW) repositories(db bun.IDB) uow.FileMetadataStorage {
  return uow.NewStorage(pg.NewFileMetadataRepository(db))
}
