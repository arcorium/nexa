package entity

import (
  "github.com/arcorium/nexa/shared/types"
  "time"
)

func NewFollow(userId, targetUserId types.Id) Follow {
  return Follow{
    FollowerId: userId,
    FolloweeId: targetUserId,
    CreatedAt:  time.Now(),
  }
}

type Follow struct {
  FollowerId types.Id
  FolloweeId types.Id
  CreatedAt  time.Time
}

type FollowCount struct {
  UserId          types.Id
  TotalFollowers  uint64
  TotalFollowings uint64
}
