package mapper

import (
	"nexa/services/authorization/internal/domain/dto"
	"nexa/services/authorization/shared/domain/entity"
	"nexa/shared/util"
)

func ToRoleResponseDTO(role *entity.Role) dto.RoleResponseDTO {
	return dto.RoleResponseDTO{
		Id:          role.Id.Underlying().String(),
		Name:        role.Name,
		Description: role.Description,
		Permissions: util.CastSlice(role.Permissions, ToPermissionResponseDTO),
	}
}
