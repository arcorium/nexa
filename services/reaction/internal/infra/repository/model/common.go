package model

import (
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/reaction/internal/domain/entity"
)

type Count struct {
  ItemId       string `bun:",scanonly"`
  LikeCount    uint64 `bun:",scanonly"`
  DislikeCount uint64 `bun:",scanonly"`
}

func (c *Count) ToDomain(itemType entity.ItemType) (entity.Count, error) {
  itemId, err := types.IdFromString(c.ItemId)
  if err != nil {
    return entity.Count{}, err
  }

  return entity.Count{
    ItemType: itemType,
    ItemId:   itemId,
    Like:     c.LikeCount,
    Dislike:  c.DislikeCount,
  }, nil
}
