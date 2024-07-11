package mapper

import (
  mailerv1 "github.com/arcorium/nexa/proto/gen/go/mailer/v1"
  sharedErr "github.com/arcorium/nexa/shared/errors"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  "nexa/services/mailer/internal/domain/dto"
)

func ToCreateTagDTO(request *mailerv1.CreateTagRequest) (dto.CreateTagDTO, error) {
  createDto := dto.CreateTagDTO{
    Name:        request.Name,
    Description: types.NewNullable(request.Description),
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
    Name:        types.NewNullable(request.Name),
    Description: types.NewNullable(request.Description),
  }, nil
}

func ToProtoTag(dto *dto.TagResponseDTO) *mailerv1.Tag {
  return &mailerv1.Tag{
    Id:          dto.Id.String(),
    Name:        dto.Name,
    Description: dto.Description,
  }
}
