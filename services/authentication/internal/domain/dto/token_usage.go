package dto

import (
  "nexa/services/authentication/internal/domain/entity"
  "nexa/shared/types"
  "nexa/shared/wrapper"
)

type TokenAddUsageDTO struct {
  Name        string `validate:"required"`
  Description wrapper.NullableString
}

func (a *TokenAddUsageDTO) ToEntity() entity.TokenUsage {
  usage := entity.TokenUsage{
    Id:   types.NewId(),
    Name: a.Name,
  }

  wrapper.SetOnNonNull(&usage.Description, a.Description)
  return usage
}

type TokenUpdateUsageDTO struct {
  Id          string `validate:"required,uuid4"`
  Name        wrapper.NullableString
  Description wrapper.NullableString
}

func (u *TokenUpdateUsageDTO) ToEntity() entity.TokenUsage {
  usage := entity.TokenUsage{
    Id: types.IdFromString(u.Id),
  }

  wrapper.SetOnNonNull(&usage.Name, u.Name)
  wrapper.SetOnNonNull(&usage.Description, u.Description)
  return usage
}

type TokenUsageResponseDTO struct {
  Id          string
  Name        string
  Description string
}
