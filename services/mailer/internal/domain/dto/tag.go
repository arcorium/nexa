package dto

import "nexa/shared/wrapper"

type CreateTagDTO struct {
  Name        string                   `json:"name" validate:"required"`
  Description wrapper.Nullable[string] `json:"description"`
}

type UpdateTagDTO struct {
  Id          string                   `json:"id" validate:"required,uuid4"`
  Name        wrapper.Nullable[string] `json:"name"`
  Description wrapper.Nullable[string] `json:"description"`
}

type TagResponseDTO struct {
  Id          string `json:"id"`
  Name        string `json:"name"`
  Description string `json:"description"` // TODO: Handle on mapper and service
}
