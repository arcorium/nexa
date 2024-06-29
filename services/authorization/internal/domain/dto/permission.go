package dto

import (
  "nexa/services/authorization/internal/domain/entity"
  "nexa/shared/types"
  "nexa/shared/util"
  "time"
)

type PermissionCreateDTO struct {
  Code string `validate:"required"`
}

func (p *PermissionCreateDTO) ToDomain() (entity.Permission, error) {
  err := util.ValidateStruct(p)
  if err != nil {
    return entity.Permission{}, err
  }

  return entity.Permission{
    Id:        types.NewId2(),
    Code:      p.Code,
    CreatedAt: time.Now(),
  }, nil
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
  Id        string
  Code      string
  CreatedAt time.Time
}
