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
  role := entity.Role{
    Id: types.NewId2(),
  }

  wrapper.SetOnNonNull(&role.Description, c.Description)
  return role, nil
}

type RoleUpdateDTO struct {
  Id          string `validate:"required,uuid4"`
  Name        wrapper.NullableString
  Description wrapper.NullableString
}

func (u *RoleUpdateDTO) ToDomain() (entity.Role, error) {
  id, err := types.IdFromString(u.Id)
  if err != nil {
    return entity.Role{}, err
  }

  role := entity.Role{
    Id: id,
  }

  wrapper.SetOnNonNull(&role.Name, u.Name)
  wrapper.SetOnNonNull(&role.Description, u.Description)

  return role, nil
}

type RoleAddPermissionsDTO struct {
  RoleId        string   `validate:"required,uuid4"`
  PermissionIds []string `validate:"required,dive,uuid4"`
}

type RoleRemovePermissionsDTO struct {
  RoleId        string   `validate:"required,uuid4"`
  PermissionIds []string `validate:"required,dive,uuid4"`
}

type RoleAddUsersDTO struct {
  UserId  string   `validate:"required,uuid4"`
  RoleIds []string `validate:"required,dive,uuid4"`
}

type RoleRemoveUsersDTO struct {
  UserId  string   `validate:"required,uuid4"`
  RoleIds []string `validate:"required,dive,uuid4"`
}

type RoleResponseDTO struct {
  Id          string
  Name        string
  Description string

  Permissions []PermissionResponseDTO
}
