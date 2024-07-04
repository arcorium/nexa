package mapper

import (
  "nexa/services/authorization/internal/domain/dto"
  "nexa/services/authorization/internal/domain/entity"
)

func ToPermissionResponseDTO(permission *entity.Permission) dto.PermissionResponseDTO {
  return dto.PermissionResponseDTO{
    Id:        permission.Id,
    Code:      permission.Encode(),
    CreatedAt: permission.CreatedAt,
  }
}
