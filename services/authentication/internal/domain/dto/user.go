package dto

import (
  "nexa/shared/types"
)

type UserValidateResponseDTO struct {
  UserId   types.Id
  Username string
}
