package model

import (
  "context"
  "github.com/arcorium/nexa/shared/types"
  "github.com/uptrace/bun"
  "nexa/services/mailer/constant"
  "time"
)

var models = []any{
  types.Nil[Tag](),
  types.Nil[Mail](),
  types.Nil[MailTag](),
}

func RegisterBunModels(db *bun.DB) {
  db.RegisterModel(types.Nil[MailTag]())
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

func SeedDefaultTags(db *bun.DB) error {
  defaultTags := []Tag{
    {
      Id:          types.MustCreateId().String(),
      Name:        constant.EMAIL_VERIFICATION_TAG,
      Description: nil,
      CreatedAt:   time.Now().UTC(),
    },
    {
      Id:          types.MustCreateId().String(),
      Name:        constant.RESET_PASSWORD_TAG,
      Description: nil,
      CreatedAt:   time.Now().UTC(),
    },
  }

  _, err := db.NewInsert().
    Model(&defaultTags).
    Returning("NULL").
    Exec(context.Background())

  return err
}
