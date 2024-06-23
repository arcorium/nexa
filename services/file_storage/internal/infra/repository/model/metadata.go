package model

import (
  "database/sql"
  "github.com/uptrace/bun"
  domain "nexa/services/file_storage/internal/domain/entity"
  "nexa/shared/types"
  "nexa/shared/util/repo"
  "nexa/shared/variadic"
  "nexa/shared/wrapper"
  "time"
)

type FileMapOption = repo.DataAccessModelMapOption[*domain.FileMetadata, *FileMetadata]

func FromFileDomain(domain *domain.FileMetadata, opts ...FileMapOption) FileMetadata {
  obj := FileMetadata{
    Id:              domain.Id.Underlying().String(),
    Filename:        domain.Name,
    FileType:        sql.NullInt64{Int64: int64(domain.Type.Underlying()), Valid: true},
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

  Id       string        `bun:",type:uuid,pk,nullzero"`
  Filename string        `bun:",unique,notnull,nullzero"`
  FileType sql.NullInt64 `bun:",type:smallint,notnull"`
  Size     uint64        `bun:",notnull,nullzero"`
  IsPublic sql.NullBool  `bun:",notnull,default:false"`

  StorageProvider sql.NullInt64 `bun:",type:smallint,notnull"` // NOTE: Bun only able to use sql.NullBool and sql.NullInt64 for integer
  StoragePath     string        `bun:",notnull"`               // Relative

  CreatedAt time.Time `bun:",nullzero,notnull"`
  UpdatedAt time.Time `bun:",nullzero"`
}

func (m *FileMetadata) ToDomain() domain.FileMetadata {
  return domain.FileMetadata{
    Id:           wrapper.DropError(types.IdFromString(m.Id)),
    Name:         m.Filename,
    Type:         domain.FileType(m.FileType.Int64),
    Size:         m.Size,
    IsPublic:     m.IsPublic.Bool,
    Provider:     domain.StorageProvider(m.StorageProvider.Int64),
    ProviderPath: m.StoragePath,
    CreatedAt:    m.CreatedAt,
    LastModified: m.UpdatedAt,
  }
}
