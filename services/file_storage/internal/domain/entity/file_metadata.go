package entity

import (
  "database/sql"
  "nexa/shared/types"
  "time"
)

type FileMetadata struct {
  Id       types.Id
  Name     string
  MimeType string
  Size     uint64
  IsPublic bool

  Provider     StorageProvider
  ProviderPath string // File path on provider (relative)
  FullPath     string

  CreatedAt    time.Time
  LastModified time.Time
}

type PatchedFileMetadata struct {
  Id types.Id

  IsPublic     types.NullableBool
  Provider     types.Nullable[StorageProvider]
  ProviderPath string
  FullPath     types.NullableString
}

func (p *PatchedFileMetadata) SqlIsPublic() sql.NullBool {
  if !p.IsPublic.HasValue() {
    return sql.NullBool{}
  }

  return sql.NullBool{
    Bool:  p.IsPublic.RawValue(),
    Valid: true,
  }
}

func (p *PatchedFileMetadata) SqlStorageProvider() sql.NullInt64 {
  if !p.Provider.HasValue() {
    return sql.NullInt64{}
  }

  return sql.NullInt64{
    Int64: int64(p.Provider.RawValue()),
    Valid: true,
  }
}
