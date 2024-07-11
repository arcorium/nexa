package dto

import (
  "github.com/arcorium/nexa/shared/types"
)

type UserResponseDTO struct {
  UserId   types.Id
  Username string
}
