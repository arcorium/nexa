package uow

import "nexa/services/file_storage/internal/domain/repository"

func NewStorage(metadata repository.IFileMetadata) FileMetadataStorage {
  return FileMetadataStorage{
    metadata: metadata,
  }
}

type FileMetadataStorage struct {
  metadata repository.IFileMetadata
}

func (m *FileMetadataStorage) Metadata() repository.IFileMetadata {
  return m.metadata
}
