package entity

import sharedErr "github.com/arcorium/nexa/shared/errors"

const (
  ItemPostComment ItemType = iota
  ItemCommentReply
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
