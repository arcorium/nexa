package mapper

import (
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/file_storage/internal/domain/dto"
  domain "nexa/services/file_storage/internal/domain/entity"
)

func ToFileMetadataResponse(metadata *domain.FileMetadata) dto.FileMetadataResponseDTO {
  return dto.FileMetadataResponseDTO{
    Id:           metadata.Id,
    Name:         metadata.Name,
    Size:         metadata.Size,
    Path:         types.FilePathFromString(metadata.FullPath),
    CreatedAt:    metadata.CreatedAt,
    LastModified: metadata.LastModified,
  }
}
