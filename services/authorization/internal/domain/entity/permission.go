package entity

import (
  "nexa/shared/types"
  "time"
)

type Permission struct {
  Id        types.Id
  Code      string
  CreatedAt time.Time
}
