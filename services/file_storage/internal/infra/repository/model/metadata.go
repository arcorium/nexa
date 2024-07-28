package model

import (
  "database/sql"
  "github.com/arcorium/nexa/shared/types"
  "github.com/arcorium/nexa/shared/util/repo"
  "github.com/arcorium/nexa/shared/variadic"
  "github.com/uptrace/bun"
  entity "nexa/services/file_storage/internal/domain/entity"
  "time"
)

type FileMapOption = repo.DataAccessModelMapOption[*entity.FileMetadata, *FileMetadata]

type PatchedFileMapOption = repo.DataAccessModelMapOption[*entity.PatchedFileMetadata, *FileMetadata]

func FromPatchedDomain(ent *entity.PatchedFileMetadata, opts ...PatchedFileMapOption) FileMetadata {
  obj := FileMetadata{
    Id:              ent.Id.String(),
    IsPublic:        ent.SqlIsPublic(),
    StorageProvider: ent.SqlStorageProvider(),
    StoragePath:     ent.ProviderPath,
    FullPath:        ent.FullPath.ValueOrNil(),
  }

  variadic.New(opts...).DoAll(repo.MapOptionFunc(ent, &obj))
  return obj
}

func FromFileDomain(ent *entity.FileMetadata, opts ...FileMapOption) FileMetadata {
  var fullpath *string
  if len(ent.FullPath) > 0 {
    fullpath = &ent.FullPath
  }
  obj := FileMetadata{
    Id:              ent.Id.String(),
    Filename:        ent.Name,
    MimeType:        ent.MimeType,
    Size:            ent.Size,
    IsPublic:        sql.NullBool{Bool: ent.IsPublic, Valid: true},
    StorageProvider: sql.NullInt64{Int64: int64(ent.Provider.Underlying()), Valid: true},
    StoragePath:     ent.ProviderPath,
    FullPath:        fullpath,
  }

  variadic.New(opts...).DoAll(repo.MapOptionFunc(ent, &obj))
  return obj
}

type FileMetadata struct {
  bun.BaseModel `bun:"table:file_metadata"`

  Id       string       `bun:",type:uuid,pk,nullzero"`
  Filename string       `bun:",notnull,nullzero"`
  MimeType string       `bun:",nullzero,notnull"`
  Size     uint64       `bun:",notnull,nullzero"`
  IsPublic sql.NullBool `bun:",notnull,default:false"`

  StorageProvider sql.NullInt64 `bun:",type:smallint,notnull"` // NOTE: Bun only able to use sql.NullBool and sql.NullInt64 for integer
  StoragePath     string        `bun:",nullzero,notnull"`      // Relative
  FullPath        *string       `bun:","`

  CreatedAt time.Time `bun:",nullzero,notnull"`
  UpdatedAt time.Time `bun:",nullzero"`
}

func (m *FileMetadata) ToDomain() (entity.FileMetadata, error) {
  id, err := types.IdFromString(m.Id)
  if err != nil {
    return entity.FileMetadata{}, err
  }

  provider, err := entity.NewStorageProvider(uint8(m.StorageProvider.Int64))
  if err != nil {
    return entity.FileMetadata{}, err
  }

  return entity.FileMetadata{
    Id:           id,
    Name:         m.Filename,
    MimeType:     m.MimeType,
    Size:         m.Size,
    IsPublic:     m.IsPublic.Bool,
    Provider:     provider,
    ProviderPath: m.StoragePath,
    FullPath:     types.OnNil(m.FullPath, ""),
    CreatedAt:    m.CreatedAt,
    LastModified: m.UpdatedAt,
  }, nil
}
