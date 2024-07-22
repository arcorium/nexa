package entity

import (
  sharedErr "github.com/arcorium/nexa/shared/errors"
)

const (
  ItemPost ItemType = iota
  ItemComment
  ItemUnknown
)

func NewItemType(val uint8) (ItemType, error) {
  itemType := ItemType(val)
  if !itemType.Valid() {
    return ItemUnknown, sharedErr.ErrEnumOutOfBounds
  }
  return itemType, nil
}

type ItemType uint8

func (i ItemType) Underlying() uint8 {
  return uint8(i)
}

func (i ItemType) Valid() bool {
  return i.Underlying() < ItemUnknown.Underlying()
}

const (
  ReactionLike ReactionType = iota
  ReactionDislike
  ReactionUnknown
)

func NewReactionType(val uint8) (ReactionType, error) {
  reactType := ReactionType(val)
  if !reactType.Valid() {
    return ReactionUnknown, sharedErr.ErrEnumOutOfBounds
  }
  return reactType, nil
}

type ReactionType uint8

func (i ReactionType) Underlying() uint8 {
  return uint8(i)
}

func (i ReactionType) Valid() bool {
  return i.Underlying() < ReactionUnknown.Underlying()
}
