package model

import (
  "github.com/uptrace/bun"
  "nexa/shared/database"
  "nexa/shared/types"
)

var models = []any{
  types.Nil[Permission](),
  types.Nil[Role](),
  types.Nil[RolePermission](),
  types.Nil[UserRole](),
}

func RegisterBunModels(db *bun.DB) {
  database.RegisterBunModels(db, types.Nil[RolePermission]())
  database.RegisterBunModels(db, types.Nil[UserRole]())
  database.RegisterBunModels(db, models...)
}

func CreateTables(db bun.IDB) error {
  return database.CreateTables(db, models...)
}
