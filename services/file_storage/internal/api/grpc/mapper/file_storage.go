package mapper

import (
  storagev1 "github.com/arcorium/nexa/proto/gen/go/file_storage/v1"
  sharedErr "github.com/arcorium/nexa/shared/errors"
  "github.com/arcorium/nexa/shared/types"
  "google.golang.org/protobuf/types/known/timestamppb"
  "nexa/services/file_storage/internal/domain/dto"
  "time"
)

func ToUpdateMetadataDTO(request *storagev1.UpdateFileRequest) (dto.UpdateFileMetadataDTO, error) {
  fileId, err := types.IdFromString(request.FileId)
  if err != nil {
    return dto.UpdateFileMetadataDTO{}, sharedErr.NewFieldError("file_id", err).ToGrpcError()
  }

  return dto.UpdateFileMetadataDTO{
    Id:       fileId,
    IsPublic: request.IsPublic,
  }, nil
}

func ToProtoFile(dto *dto.FileMetadataResponseDTO) *storagev1.File {
  var lastModified *timestamppb.Timestamp
  if !dto.LastModified.Round(time.Hour).IsZero() {
    lastModified = timestamppb.New(dto.LastModified)
  }
  return &storagev1.File{
    Id:           dto.Id.String(),
    Name:         dto.Name,
    Size:         dto.Size,
    Path:         dto.Path.Path(),
    LastModified: lastModified,
    CreatedAt:    timestamppb.New(dto.CreatedAt),
  }
}

func ToMappedProtoFile(dtos ...dto.FileMetadataResponseDTO) map[string]*storagev1.File {
  result := make(map[string]*storagev1.File)
  for _, dto := range dtos {
    result[dto.Id.String()] = ToProtoFile(&dto)
  }
  return result
}
