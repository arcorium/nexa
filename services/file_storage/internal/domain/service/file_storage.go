package service

import (
  "context"
  "github.com/arcorium/nexa/shared/status"
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/file_storage/internal/domain/dto"
)

type IFileStorage interface {
  // Store store file
  Store(ctx context.Context, file *dto.FileStoreDTO) (types.Id, status.Object)
  // Find read file based on the filename
  Find(ctx context.Context, id types.Id) (dto.FileResponseDTO, status.Object)
  FindMetadata(ctx context.Context, id types.Id) (*dto.FileMetadataResponseDTO, status.Object)
  // Delete remove or place file on bin based on the filename
  Delete(ctx context.Context, id types.Id) status.Object
  // UpdateMetadata patch some field on metadata
  Move(ctx context.Context, input *dto.UpdateFileMetadataDTO) status.Object
}
