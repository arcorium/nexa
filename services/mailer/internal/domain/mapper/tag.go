package mapper

import (
  "nexa/services/mailer/internal/domain/dto"
  domain "nexa/services/mailer/internal/domain/entity"
  "nexa/shared/types"
  "nexa/shared/wrapper"
)

func MapCreateTagDTO(tagDTO *dto.CreateTagDTO) domain.Tag {
  tag := domain.Tag{
    Id:   types.NewId2(),
    Name: tagDTO.Name,
  }

  wrapper.SetOnNonNull(&tag.Description, tagDTO.Description)
  return tag
}

func MapUpdateTagDTO(tagDTO *dto.UpdateTagDTO) domain.Tag {
  tag := domain.Tag{
    Id: wrapper.DropError(types.IdFromString(tagDTO.Id)),
  }

  wrapper.SetOnNonNull(&tag.Name, tagDTO.Name)
  wrapper.SetOnNonNull(&tag.Description, tagDTO.Description)

  return tag
}

func ToResponseDTO(tag *domain.Tag) dto.TagResponseDTO {
  return dto.TagResponseDTO{
    Id:   tag.Id.Underlying().String(),
    Name: tag.Name,
  }
}
