package dto

import "github.com/arcorium/nexa/shared/types"

type IsAuthorizationDTO struct {
  UserId             types.Id
  ExpectedPermission string `validate:"required"`
}
