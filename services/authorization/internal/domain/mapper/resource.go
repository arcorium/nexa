package mapper

import (
  "nexa/services/authorization/internal/domain/dto"
  "nexa/services/authorization/internal/domain/entity"
)

func ToResourceResponseDTO(resource *entity.Resource) dto.ResourceResponseDTO {
  return dto.ResourceResponseDTO{
    Id:          resource.Id.Underlying().String(),
    Name:        resource.Name,
    Description: resource.Description,
  }
}
