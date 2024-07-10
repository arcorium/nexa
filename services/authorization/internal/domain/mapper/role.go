package mapper

import (
  "github.com/arcorium/nexa/shared/util"
  "nexa/services/authorization/internal/domain/dto"
  "nexa/services/authorization/internal/domain/entity"
)

func ToRoleResponseDTO(role *entity.Role) dto.RoleResponseDTO {
  return dto.RoleResponseDTO{
    Id:          role.Id,
    Name:        role.Name,
    Description: role.Description,
    Permissions: util.CastSliceP(role.Permissions, ToPermissionResponseDTO),
  }
}
