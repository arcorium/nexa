package dto

import (
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/authorization/internal/domain/entity"
  "time"
)

type PermissionCreateDTO struct {
  Resource string `validate:"required"`
  Action   string `validate:"required"`
}

func (p *PermissionCreateDTO) ToDomain() (entity.Permission, error) {
  id, err := types.NewId()
  if err != nil {
    return entity.Permission{}, err
  }

  return entity.Permission{
    Id:        id,
    Resource:  p.Resource,
    Action:    p.Action,
    CreatedAt: time.Now(),
  }, nil
}

type PermissionResponseDTO struct {
  Id        types.Id
  Code      string
  CreatedAt time.Time
}
