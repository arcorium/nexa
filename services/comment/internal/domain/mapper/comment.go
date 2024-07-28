package mapper

import (
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/comment/internal/domain/dto"
  "nexa/services/comment/internal/domain/entity"
)

func ToCommentResponse(comment *entity.Comment) dto.CommentResponseDTO {
  return dto.CommentResponseDTO{
    Id:         comment.Id,
    PostId:     comment.PostId,
    UserId:     comment.UserId,
    Content:    comment.Content,
    LastEdited: types.NewNullableTime(comment.UpdatedAt),
    CreatedAt:  comment.CreatedAt,
    //TotalLikes:    responseDTO.TotalLikes,
    //TotalDislikes: responseDTO.TotalDislikes,
    //Replies:       nil,
  }
}

func ToCommentsResponse(comments []entity.Comment) []dto.CommentResponseDTO {
  type Wrapper struct {
    Index int
  }

  var result []dto.CommentResponseDTO
  var ids = map[types.Id][]Wrapper{} // index on replies

  for i := 0; i < len(comments); i++ {
    val := &comments[i]
    //react := &reactions[i]

    // Base comment
    if !val.IsReply() {
      result = append(result, ToCommentResponse(val))
      ids[val.Id] = []Wrapper{{i}}
    } else {
      // Replies
      parentIndices := ids[val.Parent.Id]
      // It could be nil if the parent doesn't exists (for reply only)
      if parentIndices == nil {
        result = append(result, ToCommentResponse(val))
        ids[val.Id] = []Wrapper{{i}} // Could become parent
        continue
      }
      var parent = &result[parentIndices[0].Index]
      for j := 1; j < len(parentIndices); j++ {
        parent = &parent.Replies[parentIndices[j].Index]
      }

      // Set as child
      parent.Replies = append(parent.Replies, ToCommentResponse(val))
      insertedIndex := len(parent.Replies) - 1

      // Add to indexes
      copied := make([]Wrapper, len(parentIndices))
      copy(copied, parentIndices)
      copied = append(copied, Wrapper{insertedIndex})
      ids[val.Id] = copied
    }
  }
  return result
}
