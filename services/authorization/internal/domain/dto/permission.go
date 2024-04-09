package dto

import (
  entity2 "nexa/services/authorization/internal/domain/entity"
  "nexa/shared/types"
  "nexa/shared/util"
)

type PermissionCreateDTO struct {
  ResourceId string `validate:"required,uuid4"`
  ActionId   string `validate:"required,uuid4"`
}

func (p *PermissionCreateDTO) ToDomain() entity2.Permission {
  return entity2.Permission{
    Id:       types.NewId(),
    Resource: entity2.Resource{Id: types.IdFromString(p.ResourceId)},
    Action:   entity2.Action{Id: types.IdFromString(p.ActionId)},
  }
}

type InternalCheckUserPermissionDTO struct {
  Resource string `validate:"required"`
  Action   string `validate:"required"`
}

func (i *InternalCheckUserPermissionDTO) ToDomain() entity2.Permission {
  return entity2.Permission{
    Resource: entity2.Resource{Name: i.Resource},
    Action:   entity2.Action{Name: i.Action},
  }
}

type CheckUserPermissionDTO struct {
  UserId      string                           `validate:"required,uuid4"`
  Permissions []InternalCheckUserPermissionDTO `validate:"required"`
}

func (c *CheckUserPermissionDTO) ToDomain() []entity2.Permission {
  return util.CastSlice(c.Permissions, func(from *InternalCheckUserPermissionDTO) entity2.Permission {
    return from.ToDomain()
  })
}
