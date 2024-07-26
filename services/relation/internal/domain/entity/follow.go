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

func NewFollows(userId types.Id, targetUserIds ...types.Id) []Follow {
  var follows []Follow
  for _, targetUserId := range targetUserIds {
    follows = append(follows, NewFollow(userId, targetUserId))
  }
  return follows
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
