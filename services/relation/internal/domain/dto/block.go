package dto

import (
  "github.com/arcorium/nexa/shared/types"
  "time"
)

type BlockCountResponseDTO struct {
  UserId types.Id
  Total  uint64
}

type BlockResponseDTO struct {
  UserId    types.Id
  CreatedAt time.Time
}
