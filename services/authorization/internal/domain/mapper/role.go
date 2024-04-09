package mapper

import (
  "nexa/services/authorization/internal/domain/entity"
  sharedDto "nexa/services/authorization/shared/domain/dto"
  "nexa/shared/util"
)

func ToRoleResponseDTO(role *entity.Role) sharedDto.RoleResponseDTO {
  return sharedDto.RoleResponseDTO{
    Id:          role.Id.Underlying().String(),
    Name:        role.Name,
    Description: role.Description,
    Permissions: util.CastSlice(role.Permissions, ToPermissionResponseDTO),
  }
}
