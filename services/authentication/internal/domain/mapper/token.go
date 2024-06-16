package mapper

import (
  "nexa/services/authentication/internal/domain/dto"
  "nexa/services/authentication/shared/domain/entity"
)

func ToTokenUsageResponse(usage *entity.TokenUsage) dto.TokenUsageResponseDTO {
  return dto.TokenUsageResponseDTO{
    Id:          usage.Id.Underlying().String(),
    Name:        usage.Name,
    Description: usage.Description,
  }
}
