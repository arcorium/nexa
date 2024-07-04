package model

import (
  "github.com/uptrace/bun"
  "nexa/shared/database"
  "nexa/shared/types"
)

var models = []any{
  types.Nil[User](),
  types.Nil[Profile](),
}

func RegisterBunModels(db *bun.DB) {
  database.RegisterBunModels(db, models...)
}

func CreateTables(db bun.IDB) error {
  return database.CreateTables(db, models...)
}
