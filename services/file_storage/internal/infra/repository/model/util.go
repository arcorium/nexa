package model

import (
  "github.com/uptrace/bun"
  domain "nexa/services/file_storage/internal/domain/entity"
  "nexa/shared/database"
  "nexa/shared/types"
  sharedUtil "nexa/shared/util"
  "time"
)

var models = []any{
  types.Nil[FileMetadata](),
}

func RegisterBunModels(db *bun.DB) {
  database.RegisterBunModels(db, models...)
}

func CreateTables(db bun.IDB) error {
  return database.CreateTables(db, models...)
}

func SeedFromDomain(db bun.IDB, domains ...domain.FileMetadata) error {
  result := sharedUtil.CastSliceP(domains, func(from *domain.FileMetadata) FileMetadata {
    return FromFileDomain(from, func(domain *domain.FileMetadata, model *FileMetadata) {
      model.CreatedAt = time.Now()
      model.UpdatedAt = time.Now()
    })
  })
  return database.Seed(db, result...)
}
