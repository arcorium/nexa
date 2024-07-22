package model

import (
  "github.com/arcorium/nexa/shared/types"
  "github.com/arcorium/nexa/shared/util/repo"
  "github.com/arcorium/nexa/shared/variadic"
  "github.com/uptrace/bun"
  "nexa/services/relation/internal/domain/entity"
  "time"
)

type BlockMapOption = repo.DataAccessModelMapOption[*entity.Block, *Block]

func FromBlockDomain(ent *entity.Block, opts ...BlockMapOption) Block {
  block := Block{
    BlockerId: ent.BlockerId.String(),
    BlockedId: ent.BlockedId.String(),
    CreatedAt: ent.CreatedAt,
  }

  variadic.New(opts...).
    DoAll(repo.MapOptionFunc(ent, &block))

  return block
}

type Block struct {
  bun.BaseModel `bun:"table:blocks"`

  BlockerId string `bun:",type:uuid,nullzero,pk"`
  BlockedId string `bun:",type:uuid,nullzero,pk"`

  CreatedAt time.Time `bun:",nullzero,notnull"`
}

func (b *Block) ToDomain() (entity.Block, error) {
  blockerId, err := types.IdFromString(b.BlockerId)
  if err != nil {
    return entity.Block{}, err
  }

  blockedId, err := types.IdFromString(b.BlockedId)
  if err != nil {
    return entity.Block{}, err
  }

  return entity.Block{
    BlockerId: blockerId,
    BlockedId: blockedId,
    CreatedAt: b.CreatedAt,
  }, nil
}
