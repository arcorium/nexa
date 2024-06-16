package mapper

import (
  "nexa/services/file_storage/internal/domain/dto"
  domain "nexa/services/file_storage/internal/domain/entity"
  "nexa/services/file_storage/internal/domain/external"
  "nexa/services/file_storage/util"
  "nexa/shared/types"
)

func MapFileStoreDTO(store *dto.FileStoreDTO) domain.File {
  return domain.File{
    Name:  store.Name,
    Bytes: store.Data,
    Size:  uint64(len(store.Data)),
  }
}

func MapFileMetadata(storeDTO *dto.FileStoreDTO, storage external.IStorage) (domain.FileMetadata, error) {
  id, err := types.NewId()
  if err != nil {
    return domain.FileMetadata{}, err
  }

  return domain.FileMetadata{
    Id:           id,
    Name:         storeDTO.Name,
    Type:         util.GetMimeType(storeDTO.Name),
    Size:         uint64(len(storeDTO.Data)),
    IsPublic:     storeDTO.IsPublic,
    Provider:     storage.GetProvider(),
    ProviderPath: "", // Needs to be defined outside
  }, nil
}

func ToFileResponse(file *domain.File) dto.FileResponseDTO {
  return dto.FileResponseDTO{
    Name: file.Name,
    Size: file.Size,
    Data: file.Bytes,
  }
}
