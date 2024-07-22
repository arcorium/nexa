package dto

import (
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/reaction/internal/domain/entity"
  "time"
)

type ReactionResponseDTO struct {
  UserId       types.Id
  ReactionType entity.ReactionType
  ItemType     entity.ItemType
  ItemId       types.Id
  CreatedAt    time.Time
}
