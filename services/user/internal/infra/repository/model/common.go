package model

import (
  "context"
  "github.com/arcorium/nexa/shared/database"
  "github.com/arcorium/nexa/shared/types"
  "github.com/uptrace/bun"
)

var modelsPair = []types.Pair[any, bool]{
  {types.Nil[User](), false},
  {types.Nil[Profile](), true},
}

var models = []any{
  types.Nil[User](),
  types.Nil[Profile](),
}

func RegisterBunModels(db *bun.DB) {
  database.RegisterBunModels(db, models...)
}

func CreateTables(db bun.IDB) error {
  // Need custom, because the user should not have foreign keys to profiles
  ctx := context.Background()
  for _, pair := range modelsPair {
    q := db.NewCreateTable().
      Model(pair.First).
      IfNotExists()

    if pair.Second {
      q = q.WithForeignKeys()
    }

    _, err := q.Exec(ctx)

    if err != nil {
      return err
    }
  }
  return nil
}
