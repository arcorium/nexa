package dto

import "nexa/shared/types"

type IsAuthorizationDTO struct {
  UserId             types.Id
  ExpectedPermission string `validate:"required"`
}
