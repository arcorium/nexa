package model

import (
  "context"
  "github.com/uptrace/bun"
  "nexa/shared/types"
)

var models = []any{
  types.Nil[Tag](),
  types.Nil[Mail](),
  types.Nil[MailTag](),
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

// InsertDefaultData insert default data for Tag
func InsertDefaultData(db *bun.DB) error {
  ctx := context.Background()

  RegisterBunModels(db)
  var err error
  for _, model := range DefaultTags {
    _, err = db.NewInsert().
      Model(&model).
      Returning("NULL").
      Exec(ctx)
  }
  return err
}
