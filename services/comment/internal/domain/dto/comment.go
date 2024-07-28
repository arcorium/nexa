package dto

import (
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/comment/internal/domain/entity"
  "time"
)

type CommentResponseDTO struct {
  Id         types.Id
  PostId     types.Id
  UserId     types.Id
  Content    string
  LastEdited types.NullableTime
  CreatedAt  time.Time

  //TotalLikes    uint64
  //TotalDislikes uint64

  Replies []CommentResponseDTO
}

type CreateCommentDTO struct {
  ParentId types.NullableId
  PostId   types.Id
  Content  string
}

func (c *CreateCommentDTO) ToDomain(userId types.Id) (entity.Comment, error) {
  id, err := types.NewId()
  if err != nil {
    return entity.Comment{}, err
  }

  var parent *entity.Comment
  if c.ParentId.HasValue() {
    parent = &entity.Comment{Id: c.ParentId.RawValue()}
  }

  return entity.Comment{
    Id:        id,
    PostId:    c.PostId,
    UserId:    userId,
    Content:   c.Content,
    CreatedAt: time.Now(),
    Parent:    parent,
  }, nil
}

type EditCommentDTO struct {
  Id      types.Id
  Content string
}

type DeleteCommentDTO struct {
  PostId types.Id
  UserId types.NullableId // Allow authorized user to remove other users comments
}

type GetPostsCommentsDTO struct {
  PostId    types.Id
  ShowReply bool
}

type FindCommentByIdDTO struct {
  CommentId types.Id
  ShowReply bool
}

type GetCommentsRepliesDTO struct {
  CommentId types.Id
  ShowReply bool
}

type GetCommentsCountDTO struct {
  ItemType entity.ItemType
  ItemIds  []types.Id
}
