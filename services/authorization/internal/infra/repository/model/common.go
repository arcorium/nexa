package model

import (
  "github.com/arcorium/nexa/shared/database"
  "github.com/arcorium/nexa/shared/types"
  "github.com/uptrace/bun"
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
