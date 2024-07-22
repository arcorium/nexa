package model

import (
  "github.com/arcorium/nexa/shared/database"
  "github.com/arcorium/nexa/shared/types"
  "github.com/uptrace/bun"
)

var models = []any{
  types.Nil[Reaction](),
}

func RegisterBunModels(db *bun.DB) {
  database.RegisterBunModels(db, models...)
}

func CreateTables(db bun.IDB) error {
  return database.CreateTables(db, models...)
}
