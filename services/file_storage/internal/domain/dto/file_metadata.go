package dto

import (
  entity "nexa/services/file_storage/internal/domain/entity"
  "nexa/shared/types"
  "time"
)

type UpdateFileMetadataDTO struct {
  Id       types.Id
  IsPublic bool
}

func (u *UpdateFileMetadataDTO) ToDomain(newRelPath string) entity.PatchedFileMetadata {
  md := entity.PatchedFileMetadata{
    Id:           u.Id,
    IsPublic:     types.SomeNullable(u.IsPublic),
    ProviderPath: newRelPath,
    FullPath:     types.SomeNullable(""), // Set fullpath to empty
  }

  return md
}

type FileMetadataResponseDTO struct {
  Id           types.Id
  Name         string
  Size         uint64
  Path         types.FilePath
  CreatedAt    time.Time
  LastModified time.Time
}
