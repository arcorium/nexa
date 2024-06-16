package model

import (
  "context"
  "github.com/uptrace/bun"
  "nexa/shared/util"
)

var models = []any{
  util.Nil[FileMetadata](),
}

func RegisterBunModels(db *bun.DB) {
  db.RegisterModel(models...)
}

func CreateTables(db *bun.DB) error {
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
