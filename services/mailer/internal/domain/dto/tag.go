package dto

import (
  domain "nexa/services/mailer/internal/domain/entity"
  "nexa/shared/types"
)

type CreateTagDTO struct {
  Name        string `validate:"required"`
  Description types.NullableString
}

func (c *CreateTagDTO) ToDomain() (domain.Tag, error) {
  id, err := types.NewId()
  if err != nil {
    return domain.Tag{}, err
  }

  tag := domain.Tag{
    Id:   id,
    Name: c.Name,
  }
  types.SetOnNonNull(&tag.Description, c.Description)

  return tag, nil
}

type UpdateTagDTO struct {
  Id          types.Id
  Name        types.NullableString
  Description types.NullableString
}

func (u *UpdateTagDTO) ToDomain() domain.Tag {
  tag := domain.Tag{
    Id: u.Id,
  }

  types.SetOnNonNull(&tag.Name, u.Name)
  types.SetOnNonNull(&tag.Description, u.Description)

  return tag
}

type TagResponseDTO struct {
  Id          types.Id
  Name        string
  Description string
}
