package entity

import sharedErr "github.com/arcorium/nexa/shared/errors"

func NewFollowStatus(val uint8) (FollowStatus, error) {
  stat := FollowStatus(val)
  if !stat.Valid() {
    return FollowStatusUnknown, sharedErr.ErrEnumOutOfBounds
  }
  return stat, nil
}

type FollowStatus uint8

const (
  FollowStatusNone FollowStatus = iota
  FollowStatusFollower
  FollowStatusMutual
  FollowStatusUnknown
)

func (s FollowStatus) Underlying() uint8 {
  return uint8(s)
}

func (s FollowStatus) Valid() bool {
  return s.Underlying() < FollowStatusUnknown.Underlying()
}

func NewBlockStatus(val uint8) (BlockStatus, error) {
  stat := BlockStatus(val)
  if !stat.Valid() {
    return BlockStatusUnknown, sharedErr.ErrEnumOutOfBounds
  }
  return stat, nil
}

type BlockStatus uint8

const (
  BlockStatusNone BlockStatus = iota
  BlockStatusBlocked
  BlockStatusUnknown
)

func (s BlockStatus) Underlying() uint8 {
  return uint8(s)
}

func (s BlockStatus) Valid() bool {
  return s.Underlying() < BlockStatusUnknown.Underlying()
}
