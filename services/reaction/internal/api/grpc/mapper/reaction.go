package mapper

import (
  "github.com/arcorium/nexa/proto/gen/go/common"
  reactionv1 "github.com/arcorium/nexa/proto/gen/go/reaction/v1"
  sharedDto "github.com/arcorium/nexa/shared/dto"
  sharedErr "github.com/arcorium/nexa/shared/errors"
  "google.golang.org/protobuf/types/known/timestamppb"
  "nexa/services/reaction/internal/domain/dto"
  "nexa/services/reaction/internal/domain/entity"
)

func ToEntityItemType(itemType reactionv1.Type) (entity.ItemType, error) {
  switch itemType {
  case reactionv1.Type_POST:
    return entity.ItemPost, nil
  case reactionv1.Type_COMMENT:
    return entity.ItemComment, nil
  }
  return entity.ItemUnknown, sharedErr.ErrEnumOutOfBounds
}

func toProtoReactionType(reactType entity.ReactionType) reactionv1.ReactionType {
  switch reactType {
  case entity.ReactionLike:
    return reactionv1.ReactionType_LIKE
  case entity.ReactionDislike:
    return reactionv1.ReactionType_DISLIKE
  }
  return reactionv1.ReactionType(reactType.Underlying())
}

func ToCommonPagedOutput[T any](result *sharedDto.PagedElementResult[T]) *common.PagedElementOutput {
  return &common.PagedElementOutput{
    Element:       result.Element,
    Page:          result.Page,
    TotalElements: result.TotalElements,
    TotalPages:    result.TotalPages,
  }
}

func ToProtoReaction(ent *dto.ReactionResponseDTO) *reactionv1.Reaction {
  return &reactionv1.Reaction{
    // TODO: Add username
    UserId:       ent.UserId.String(),
    ReactionType: toProtoReactionType(ent.ReactionType),
    CreatedAt:    timestamppb.New(ent.CreatedAt),
  }
}

func ToProtoCount(responseDTO *dto.CountResponseDTO) *reactionv1.Count {
  return &reactionv1.Count{
    TotalLikes:    responseDTO.Like,
    TotalDislikes: responseDTO.Dislike,
  }
}
