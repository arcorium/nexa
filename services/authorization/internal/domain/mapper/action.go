package mapper

import (
	"nexa/services/authorization/internal/domain/dto"
	"nexa/services/authorization/shared/domain/entity"
)

func ToActionResponseDTO(action *entity.Action) dto.ActionResponseDTO {
	if action == nil {
		return dto.ActionResponseDTO{}
	}

	return dto.ActionResponseDTO{
		Id:          action.Id.Underlying().String(),
		Name:        action.Name,
		Description: action.Name,
	}
}
