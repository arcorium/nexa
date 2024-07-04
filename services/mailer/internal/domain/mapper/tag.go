package mapper

import (
  "nexa/services/mailer/internal/domain/dto"
  domain "nexa/services/mailer/internal/domain/entity"
)

func ToTagResponseDTO(tag *domain.Tag) dto.TagResponseDTO {
  return dto.TagResponseDTO{
    Id:          tag.Id,
    Name:        tag.Name,
    Description: tag.Description,
  }
}
