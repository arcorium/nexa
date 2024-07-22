package model

import (
  "github.com/arcorium/nexa/shared/types"
  "github.com/arcorium/nexa/shared/util/repo"
  "github.com/arcorium/nexa/shared/variadic"
  "github.com/uptrace/bun"
  "nexa/services/reaction/internal/domain/entity"
  "time"
)

type ReactionMapOption = repo.DataAccessModelMapOption[*entity.Reaction, *Reaction]

func FromReactionDomain(ent *entity.Reaction, opts ...ReactionMapOption) Reaction {
  post := Reaction{
    UserId:    ent.UserId.String(),
    Reaction:  ent.ReactionType.Underlying(),
    ItemType:  ent.ItemType.Underlying(),
    ItemId:    ent.ItemId.String(),
    CreatedAt: ent.CreatedAt,
  }

  variadic.New(opts...).
    DoAll(repo.MapOptionFunc(ent, &post))

  return post
}

type Reaction struct {
  bun.BaseModel `bun:"table:reactions"`

  UserId   string `bun:",type:uuid,nullzero,pk"`
  Reaction uint8  `bun:",notnull"`
  ItemType uint8  `bun:",notnull,pk"`
  ItemId   string `bun:",type:uuid,nullzero,pk"`

  CreatedAt time.Time `bun:",notnull,nullzero"`
}

func (r *Reaction) ToDomain() (entity.Reaction, error) {
  userId, err := types.IdFromString(r.UserId)
  if err != nil {
    return entity.Reaction{}, err
  }

  itemId, err := types.IdFromString(r.ItemId)
  if err != nil {
    return entity.Reaction{}, err
  }

  reaction, err := entity.NewReactionType(r.Reaction)
  if err != nil {
    return entity.Reaction{}, err
  }

  itemType, err := entity.NewItemType(r.ItemType)
  if err != nil {
    return entity.Reaction{}, err
  }

  return entity.Reaction{
    UserId:       userId,
    ReactionType: reaction,
    ItemType:     itemType,
    ItemId:       itemId,
    CreatedAt:    r.CreatedAt,
  }, nil
}
