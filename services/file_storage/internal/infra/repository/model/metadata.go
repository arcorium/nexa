package model

import (
  "database/sql"
  "github.com/uptrace/bun"
  domain "nexa/services/file_storage/internal/domain/entity"
  "nexa/shared/types"
  "nexa/shared/util/repo"
  "nexa/shared/variadic"
  "time"
)

type FileMapOption = repo.DataAccessModelMapOption[*domain.FileMetadata, *FileMetadata]

func FromFileDomain(domain *domain.FileMetadata, opts ...FileMapOption) FileMetadata {
  obj := FileMetadata{
    Id:              domain.Id.Underlying().String(),
    Filename:        domain.Name,
    MimeType:        domain.MimeType,
    Size:            domain.Size,
    IsPublic:        sql.NullBool{Bool: domain.IsPublic, Valid: true},
    StorageProvider: sql.NullInt64{Int64: int64(domain.Provider.Underlying()), Valid: true},
    StoragePath:     domain.ProviderPath,
  }

  variadic.New(opts...).DoAll(repo.MapOptionFunc(domain, &obj))
  return obj
}

type FileMetadata struct {
  bun.BaseModel `bun:"table:file_metadata"`

  Id       string       `bun:",type:uuid,pk,nullzero"`
  Filename string       `bun:",unique,notnull,nullzero"`
  MimeType string       `bun:",nullzero,notnull"`
  Size     uint64       `bun:",notnull,nullzero"`
  IsPublic sql.NullBool `bun:",notnull,default:false"`

  StorageProvider sql.NullInt64 `bun:",type:smallint,notnull"` // NOTE: Bun only able to use sql.NullBool and sql.NullInt64 for integer
  StoragePath     string        `bun:",notnull"`               // Relative

  CreatedAt time.Time `bun:",nullzero,notnull"`
  UpdatedAt time.Time `bun:",nullzero"`
}

func (m *FileMetadata) ToDomain() (domain.FileMetadata, error) {
  id, err := types.IdFromString(m.Id)
  if err != nil {
    return domain.FileMetadata{}, err
  }

  provider, err := domain.NewStorageProvider(uint8(m.StorageProvider.Int64))
  if err != nil {
    return domain.FileMetadata{}, err
  }

  return domain.FileMetadata{
    Id:           id,
    Name:         m.Filename,
    MimeType:     m.MimeType,
    Size:         m.Size,
    IsPublic:     m.IsPublic.Bool,
    Provider:     provider,
    ProviderPath: m.StoragePath,
    CreatedAt:    m.CreatedAt,
    LastModified: m.UpdatedAt,
  }, nil
}
