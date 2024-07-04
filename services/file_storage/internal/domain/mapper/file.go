package mapper

import (
  "nexa/services/file_storage/internal/domain/dto"
  domain "nexa/services/file_storage/internal/domain/entity"
)

func ToFileResponse(file *domain.File) dto.FileResponseDTO {
  return dto.FileResponseDTO{
    Name: file.Name,
    Size: file.Size,
    Data: file.Bytes,
  }
}
