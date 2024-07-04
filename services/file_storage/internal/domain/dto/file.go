package dto

import (
  "crypto/sha512"
  "fmt"
  domain "nexa/services/file_storage/internal/domain/entity"
  "nexa/services/file_storage/internal/domain/external"
  "nexa/services/file_storage/util"
  "nexa/shared/types"
  "path"
)

type FileResponseDTO struct {
  Name string
  Size uint64
  Data []byte
}

type FileStoreDTO struct {
  Name     string // Filename is ignored and only extension is used
  Data     []byte `validate:"required"`
  IsPublic bool   `validate:"required"`
}

func (s *FileStoreDTO) ToDomain(storage external.IStorage) (domain.File, domain.FileMetadata, error) {
  id, err := types.NewId()
  if err != nil {
    return domain.File{}, domain.FileMetadata{}, err
  }

  ext := path.Ext(s.Name)
  filename := fmt.Sprintf("%s.%s", id.Hash(sha512.New()), ext)

  file := domain.File{
    Name:     filename,
    Bytes:    s.Data,
    Size:     uint64(len(s.Data)),
    IsPublic: s.IsPublic,
  }

  metadata := domain.FileMetadata{
    Id:       id,
    Name:     filename,
    MimeType: util.GetMimeType(s.Name),
    Size:     uint64(len(s.Data)),
    IsPublic: s.IsPublic,
    Provider: storage.GetProvider(),
  }

  return file, metadata, nil
}
