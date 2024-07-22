package model

import (
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/relation/internal/domain/entity"
)

type FollowCount struct {
  UserId          string `bun:"user_id,scanonly"`
  TotalFollowers  uint64 `bun:",scanonly"`
  TotalFollowings uint64 `bun:",scanonly"`
}

func (f *FollowCount) ToDomain() (entity.FollowCount, error) {
  userId, err := types.IdFromString(f.UserId)
  if err != nil {
    return entity.FollowCount{}, err
  }

  return entity.FollowCount{
    UserId:          userId,
    TotalFollowers:  f.TotalFollowers,
    TotalFollowings: f.TotalFollowings,
  }, nil
}

type BlockedCount struct {
  UserId       string `bun:"user_id,scanonly"`
  TotalBlocked uint64 `bun:"total_blocked,scanonly"`
}

func (f *BlockedCount) ToDomain() (entity.BlockCount, error) {
  userId, err := types.IdFromString(f.UserId)
  if err != nil {
    return entity.BlockCount{}, err
  }

  return entity.BlockCount{
    UserId:       userId,
    TotalBlocked: f.TotalBlocked,
  }, nil
}
