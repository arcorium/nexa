package model

import (
  "context"
  "database/sql"
  "github.com/uptrace/bun"
  domain "nexa/services/file_storage/internal/domain/entity"
  sharedUtil "nexa/shared/util"
  "nexa/shared/wrapper"
  "time"
)

var models = []any{
  sharedUtil.Nil[FileMetadata](),
}

func RegisterBunModels(db *bun.DB) {
  db.RegisterModel(models...)
}

func CreateTables(db bun.IDB) error {
  ctx := context.Background()
  for _, model := range models {
    _, err := db.NewCreateTable().
      Model(model).
      IfNotExists().
      WithForeignKeys().
      Exec(ctx)

    if err != nil {
      return err
    }
  }
  return nil
}

func SeedDatabase(db bun.IDB, domains ...domain.FileMetadata) error {
  result := sharedUtil.CastSliceP(domains, func(from *domain.FileMetadata) FileMetadata {
    return FromFileDomain(from, func(domain *domain.FileMetadata, model *FileMetadata) {
      model.CreatedAt = time.Now()
      model.UpdatedAt = time.Now()
    })
  })

  res, err := db.NewInsert().
    Model(&result).
    Returning("NULL").
    Exec(context.Background())

  if err != nil {
    return err
  }

  if wrapper.Must(res.RowsAffected()) == 0 {
    return sql.ErrNoRows
  }

  return nil
}
