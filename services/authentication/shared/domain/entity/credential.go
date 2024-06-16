package entity

import (
  "nexa/services/authentication/shared/domain/valueobject"
  "nexa/shared/types"
)

type Credential struct {
  Id            types.Id
  UserId        types.Id
  AccessTokenId types.Id
  Device        valueobject.Device
  RefreshToken  string
}
