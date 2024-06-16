package mapper

import (
  "google.golang.org/protobuf/types/known/timestamppb"
  protoV1 "nexa/proto/generated/golang/file_storage/v1"
  "nexa/services/file_storage/internal/domain/dto"
  "nexa/shared/wrapper"
)

func ToUpdateMetadataDTO(request *protoV1.UpdateFileMetadataRequest) dto.UpdateFileMetadataDTO {
  return dto.UpdateFileMetadataDTO{
    Id:       request.FileId,
    Name:     wrapper.NewNullable(request.Filename),
    IsPublic: wrapper.NewNullable(request.IsPublic),
  }
}

func ToProtoFile(dto *dto.FileMetadataResponseDTO) *protoV1.File {
  return &protoV1.File{
    Id:           dto.Id,
    Name:         dto.Name,
    Size:         dto.Size,
    Path:         dto.Path,
    CreatedAt:    timestamppb.New(dto.CreatedAt),
    LastModified: timestamppb.New(dto.LastModified),
  }
}
