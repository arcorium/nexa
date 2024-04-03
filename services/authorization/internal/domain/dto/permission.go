package dto

import (
	"nexa/services/authorization/shared/domain/entity"
	"nexa/shared/types"
	"nexa/shared/util"
)

type PermissionCreateDTO struct {
	ResourceId string `validate:"required,uuid4"`
	ActionId   string `validate:"required,uuid4"`
}

func (p *PermissionCreateDTO) ToDomain() entity.Permission {
	return entity.Permission{
		Id:       types.NewId(),
		Resource: entity.Resource{Id: types.IdFromString(p.ResourceId)},
		Action:   entity.Action{Id: types.IdFromString(p.ActionId)},
	}
}

type InternalCheckUserPermissionDTO struct {
	Resource string `validate:"required"`
	Action   string `validate:"required"`
}

func (i *InternalCheckUserPermissionDTO) ToDomain() entity.Permission {
	return entity.Permission{
		Resource: entity.Resource{Name: i.Resource},
		Action:   entity.Action{Name: i.Action},
	}
}

type CheckUserPermissionDTO struct {
	UserId      string                           `validate:"required,uuid4"`
	Permissions []InternalCheckUserPermissionDTO `validate:"required"`
}

func (c *CheckUserPermissionDTO) ToDomain() []entity.Permission {
	return util.CastSlice(c.Permissions, func(from *InternalCheckUserPermissionDTO) entity.Permission {
		return from.ToDomain()
	})
}

type PermissionResponseDTO struct {
	Resource ResourceResponseDTO
	Action   ActionResponseDTO
	Code     string
}
