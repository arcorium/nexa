package mapper

import (
  mailerv1 "nexa/proto/gen/go/mailer/v1"
  "nexa/services/mailer/internal/domain/dto"
  "nexa/shared/wrapper"
)

func ToCreateTagDTO(request *mailerv1.CreateTagRequest) dto.CreateTagDTO {
  return dto.CreateTagDTO{
    Name:        request.Name,
    Description: wrapper.NewNullable(request.Description),
  }
}

func ToUpdateTagDTO(request *mailerv1.UpdateTagRequest) dto.UpdateTagDTO {
  return dto.UpdateTagDTO{
    Id:          request.TagId,
    Name:        wrapper.NewNullable(request.Name),
    Description: wrapper.NewNullable(request.Description),
  }
}

func ToProtoTag(dto *dto.TagResponseDTO) *mailerv1.Tag {
  return &mailerv1.Tag{
    Id:          dto.Id,
    Name:        dto.Name,
    Description: dto.Description,
  }
}
