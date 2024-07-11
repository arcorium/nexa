package dto

import (
  "crypto/sha512"
  "fmt"
  "github.com/arcorium/nexa/shared/types"
  entity "nexa/services/file_storage/internal/domain/entity"
  "nexa/services/file_storage/util"
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

func (s *FileStoreDTO) ToDomain(provider entity.StorageProvider) (entity.File, entity.FileMetadata, error) {
  id, err := types.NewId()
  if err != nil {
    return entity.File{}, entity.FileMetadata{}, err
  }

  ext := path.Ext(s.Name)
  filename := fmt.Sprintf("%s.%s", id.Hash(sha512.New()), ext)

  file := entity.File{
    Name:     filename,
    Bytes:    s.Data,
    Size:     uint64(len(s.Data)),
    IsPublic: s.IsPublic,
  }

  metadata := entity.FileMetadata{
    Id:       id,
    Name:     filename,
    MimeType: util.GetMimeType(s.Name),
    Size:     uint64(len(s.Data)),
    IsPublic: s.IsPublic,
    Provider: provider,
  }

  return file, metadata, nil
}
