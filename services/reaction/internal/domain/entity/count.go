package entity

import "github.com/arcorium/nexa/shared/types"

type Count struct {
  ItemType ItemType
  ItemId   types.Id
  Like     uint64
  Dislike  uint64
}
