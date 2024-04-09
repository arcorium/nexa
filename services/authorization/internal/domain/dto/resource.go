package dto

import (
  "nexa/services/authorization/internal/domain/entity"
  "nexa/shared/types"
  "nexa/shared/wrapper"
)

type ResourceCreateDTO struct {
  Name        string `validate:"required"`
  Description wrapper.NullableString
}

func (r *ResourceCreateDTO) ToDomain() entity.Resource {
  resource := entity.Resource{
    Id:   types.NewId(),
    Name: r.Name,
  }

  wrapper.SetOnNonNull(&resource.Description, r.Description)
  return resource
}

type ResourceUpdateDTO struct {
  Id          string `validate:"require,uuid4"`
  Name        wrapper.NullableString
  Description wrapper.NullableString
}

func (u *ResourceUpdateDTO) ToDomain() entity.Resource {
  resource := entity.Resource{
    Id: types.IdFromString(u.Id),
  }

  wrapper.SetOnNonNull(&resource.Name, u.Name)
  wrapper.SetOnNonNull(&resource.Description, u.Description)
  return resource
}
