package mapper

import (
  mailerv1 "nexa/proto/gen/go/mailer/v1"
  "nexa/services/mailer/internal/domain/dto"
  sharedErr "nexa/shared/errors"
  "nexa/shared/types"
  sharedUtil "nexa/shared/util"
  "nexa/shared/wrapper"
)

func ToCreateTagDTO(request *mailerv1.CreateTagRequest) (dto.CreateTagDTO, error) {
  createDto := dto.CreateTagDTO{
    Name:        request.Name,
    Description: wrapper.NewNullable(request.Description),
  }

  err := sharedUtil.ValidateStruct(&createDto)
  return dto.CreateTagDTO{}, err
}

func ToUpdateTagDTO(request *mailerv1.UpdateTagRequest) (dto.UpdateTagDTO, error) {
  tagId, err := types.IdFromString(request.TagId)
  if err != nil {
    return dto.UpdateTagDTO{}, sharedErr.NewFieldError("tag_id", err).ToGrpcError()
  }

  return dto.UpdateTagDTO{
    Id:          tagId,
    Name:        wrapper.NewNullable(request.Name),
    Description: wrapper.NewNullable(request.Description),
  }, nil
}

func ToProtoTag(dto *dto.TagResponseDTO) *mailerv1.Tag {
  return &mailerv1.Tag{
    Id:          dto.Id.String(),
    Name:        dto.Name,
    Description: dto.Description,
  }
}
