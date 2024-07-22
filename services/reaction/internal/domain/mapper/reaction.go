package mapper

import (
  "nexa/services/reaction/internal/domain/dto"
  "nexa/services/reaction/internal/domain/entity"
)

func ToReactionResponseDTO(reaction *entity.Reaction) dto.ReactionResponseDTO {
  return dto.ReactionResponseDTO{
    UserId:       reaction.UserId,
    ReactionType: reaction.ReactionType,
    ItemType:     reaction.ItemType,
    ItemId:       reaction.ItemId,
    CreatedAt:    reaction.CreatedAt,
  }
}

func ToCountResponseDTO(count *entity.Count) dto.CountResponseDTO {
  return dto.CountResponseDTO{
    Like:    count.Like,
    Dislike: count.Dislike,
  }
}
