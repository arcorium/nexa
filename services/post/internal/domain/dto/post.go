package dto

import (
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  "nexa/services/post/internal/domain/entity"
  "time"
)

type PostResponseDTO struct {
  Id         types.Id
  ParentPost *PostResponseDTO
  CreatorId  types.Id
  Content    string
  Visibility entity.Visibility

  TotalLikes    uint64
  TotalDislikes uint64
  TotalComments uint64
  TotalShares   uint64

  LastEdited time.Time
  CreatedAt  time.Time

  Tags      []TaggedUserDTO
  MediaUrls []string
}

type ChildPostResponseDTO struct {
  Content   string
  CreatedAt time.Time

  Tags      []TaggedUserDTO
  MediaUrls []string
}

type EditedPostResponseDTO struct {
  PostId    types.Id
  CreatorId types.Id

  EditedPosts []ChildPostResponseDTO
}

type CreatePostDTO struct {
  SharedPostId types.NullableId
  Content      types.NullableString
  Visibility   entity.Visibility
  MediaIds     []types.Id
  UserIds      []types.Id
}

func (c *CreatePostDTO) ToDomain(creatorId types.Id) entity.Post {
  var parent *entity.Post
  if c.SharedPostId.HasValue() {
    parent = &entity.Post{
      Id: c.SharedPostId.RawValue(),
    }
  }
  return entity.Post{
    Id:         types.MustCreateId(),
    ParentPost: parent,
    CreatorId:  creatorId,
    Content:    c.Content.ValueOr(""),
    Visibility: c.Visibility,
    CreatedAt:  time.Now(),
    Tags: sharedUtil.CastSlice(c.UserIds, func(id types.Id) entity.TaggedUser {
      return entity.TaggedUser{
        Id: id,
      }
    }),
    Medias: sharedUtil.CastSlice(c.MediaIds, func(id types.Id) entity.Media {
      return entity.Media{
        Id: id,
      }
    }),
  }
}

type EditPostDTO struct {
  PostId   types.Id
  Content  string
  MediaIds []types.Id
  UserIds  []types.Id
}

func (e *EditPostDTO) ToDomain(creatorId types.Id) entity.Post {
  return entity.Post{
    Id:         e.PostId,
    CreatorId:  creatorId,
    Content:    e.Content,
    LastEdited: time.Now(),
    Tags: sharedUtil.CastSlice(e.UserIds, func(id types.Id) entity.TaggedUser {
      return entity.TaggedUser{
        Id: id,
      }
    }),
    Medias: sharedUtil.CastSlice(e.MediaIds, func(id types.Id) entity.Media {
      return entity.Media{
        Id: id,
      }
    }),
  }
}

type TaggedUserDTO struct {
  UserId   types.Id
  Username string
}
