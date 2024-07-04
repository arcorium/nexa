package dto

import (
  "nexa/services/authorization/internal/domain/entity"
  "nexa/shared/types"
  "nexa/shared/wrapper"
)

type RoleCreateDTO struct {
  Name        string `validate:"required"`
  Description wrapper.NullableString
}

func (c *RoleCreateDTO) ToDomain() (entity.Role, error) {
  id, err := types.NewId()
  if err != nil {
    return entity.Role{}, nil
  }

  role := entity.Role{
    Id:   id,
    Name: c.Name,
  }

  wrapper.SetOnNonNull(&role.Description, c.Description)

  return role, nil
}

type RoleUpdateDTO struct {
  RoleId      types.Id
  Name        wrapper.NullableString
  Description wrapper.NullableString
}

func (u *RoleUpdateDTO) ToDomain() entity.Role {
  role := entity.Role{
    Id: u.RoleId,
  }

  wrapper.SetOnNonNull(&role.Name, u.Name)
  wrapper.SetOnNonNull(&role.Description, u.Description)

  return role
}

type ModifyRolesPermissionsDTO struct {
  RoleId        types.Id
  PermissionIds []types.Id `validate:"required"`
}

type ModifyUserRolesDTO struct {
  UserId  types.Id
  RoleIds []types.Id `validate:"required"`
}

type RoleResponseDTO struct {
  Id          types.Id
  Name        string
  Description string

  Permissions []PermissionResponseDTO
}
