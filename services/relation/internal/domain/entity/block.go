package entity

import (
  "github.com/arcorium/nexa/shared/types"
  "time"
)

func NewBlock(userId, targetId types.Id) Block {
  return Block{
    BlockerId: userId,
    BlockedId: targetId,
    CreatedAt: time.Now(),
  }
}

type Block struct {
  BlockerId types.Id
  BlockedId types.Id
  CreatedAt time.Time
}

type BlockCount struct {
  UserId       types.Id
  TotalBlocked uint64
}
