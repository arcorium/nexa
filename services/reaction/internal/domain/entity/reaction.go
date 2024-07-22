package entity

import (
  "github.com/arcorium/nexa/shared/types"
  "time"
)

func NewLikeReaction(userId types.Id, itemType ItemType, itemId types.Id) Reaction {
  return Reaction{
    UserId:       userId,
    ReactionType: ReactionLike,
    ItemType:     itemType,
    ItemId:       itemId,
    CreatedAt:    time.Now(),
  }
}

func NewDislikeReaction(userId types.Id, itemType ItemType, itemId types.Id) Reaction {
  return Reaction{
    UserId:       userId,
    ReactionType: ReactionDislike,
    ItemType:     itemType,
    ItemId:       itemId,
    CreatedAt:    time.Now(),
  }
}

type Reaction struct {
  UserId       types.Id
  ReactionType ReactionType
  ItemType     ItemType
  ItemId       types.Id
  CreatedAt    time.Time
}

func (r *Reaction) IsPost() bool {
  return r.ItemType == ItemPost
}

func (r *Reaction) IsComment() bool {
  return r.ItemType == ItemComment
}

type Item = types.Pair[ItemType, []types.Id]
