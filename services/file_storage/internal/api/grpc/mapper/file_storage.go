package mapper

import (
  storagev1 "github.com/arcorium/nexa/proto/gen/go/file_storage/v1"
  sharedErr "github.com/arcorium/nexa/shared/errors"
  "github.com/arcorium/nexa/shared/types"
  "google.golang.org/protobuf/types/known/timestamppb"
  "nexa/services/file_storage/internal/domain/dto"
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
  return &storagev1.File{
    Id:           dto.Id.String(),
    Name:         dto.Name,
    Size:         dto.Size,
    Path:         dto.Path.Path(),
    CreatedAt:    timestamppb.New(dto.CreatedAt),
    LastModified: timestamppb.New(dto.LastModified),
  }
}
