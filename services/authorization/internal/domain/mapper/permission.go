package mapper

import (
	"nexa/services/authorization/internal/domain/dto"
	"nexa/services/authorization/shared/domain/entity"
)

func ToPermissionResponseDTO(permission *entity.Permission) dto.PermissionResponseDTO {
	return dto.PermissionResponseDTO{
		Resource: ToResourceResponseDTO(&permission.Resource),
		Action:   ToActionResponseDTO(&permission.Action),
		Code:     permission.String(),
	}
}
