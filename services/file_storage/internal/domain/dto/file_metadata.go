package dto

import (
  domain "nexa/services/file_storage/internal/domain/entity"
  "nexa/shared/types"
  "time"
)

type UpdateFileMetadataDTO struct {
  Id       types.Id
  IsPublic bool
}

func (u *UpdateFileMetadataDTO) ToDomain() domain.FileMetadata {
  md := domain.FileMetadata{
    Id:       u.Id,
    IsPublic: u.IsPublic,
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
