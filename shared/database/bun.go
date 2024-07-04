package database

import (
  "context"
  "github.com/uptrace/bun"
)

func RegisterBunModels(db *bun.DB, models ...any) {
  db.RegisterModel(models...)
}

func CreateTables(db bun.IDB, models ...any) error {
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
