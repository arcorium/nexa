package dto

import (
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/authorization/internal/domain/entity"
)

type RoleCreateDTO struct {
  Name        string `validate:"required"`
  Description types.NullableString
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

  types.SetOnNonNull(&role.Description, c.Description)

  return role, nil
}

type RoleUpdateDTO struct {
  RoleId      types.Id
  Name        types.NullableString
  Description types.NullableString
}

func (u *RoleUpdateDTO) ToDomain() entity.PatchedRole {
  role := entity.PatchedRole{
    Id:          u.RoleId,
    Description: u.Description,
  }

  types.SetOnNonNull(&role.Name, u.Name)
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
