package entity

import (
  "nexa/shared/types"
  "time"
)

type Credential struct {
  Id            types.Id
  UserId        types.Id
  AccessTokenId types.Id
  Device        Device
  RefreshToken  string
  ExpiresAt     time.Time // Set to refresh token expiration
}
