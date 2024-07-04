package dto

import (
  "nexa/shared/types"
)

type UserResponseDTO struct {
  UserId   types.Id
  Username string
}
