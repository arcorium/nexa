package dto

import (
  "github.com/arcorium/nexa/shared/types"
  "time"
)

type FollowResponseDTO struct {
  UserId    types.Id
  CreatedAt time.Time
}

type FollowCountResponseDTO struct {
  UserId          types.Id
  TotalFollowers  uint64
  TotalFollowings uint64
}
