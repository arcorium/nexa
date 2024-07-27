package mapper

import (
  sharedUtil "github.com/arcorium/nexa/shared/util"
  "nexa/services/post/internal/domain/dto"
  "nexa/services/post/internal/domain/entity"
)

func ToTaggedUserDTO(user *entity.TaggedUser) dto.TaggedUserDTO {
  return dto.TaggedUserDTO{
    UserId: user.Id,
  }
}

func ToMediaDTO(media *entity.Media) string {
  return media.Url
}

func ToPostResponseDTO(post *entity.Post) dto.PostResponseDTO {
  var parent *dto.PostResponseDTO
  if post.IsShare() {
    temp := ToPostResponseDTO(post.ParentPost) // Escaped
    parent = &temp
  }

  return dto.PostResponseDTO{
    Id:         post.Id,
    ParentPost: parent,
    CreatorId:  post.CreatorId,
    Content:    post.Content,
    Visibility: post.Visibility,
    //TotalShares: post.Shares,
    LastEdited: post.LastEdited,
    CreatedAt:  post.CreatedAt,
    Tags:       sharedUtil.CastSliceP(post.Tags, ToTaggedUserDTO),
    MediaUrls:  sharedUtil.CastSliceP(post.Medias, ToMediaDTO),
  }
}

func ToEditedPostResponseDTO(post *entity.Post) dto.EditedPostResponseDTO {
  var children []dto.ChildPostResponseDTO
  children = append(children, dto.ChildPostResponseDTO{
    Content:   post.Content,
    CreatedAt: post.LastEdited,
    Tags:      sharedUtil.CastSliceP(post.Tags, ToTaggedUserDTO),
    MediaUrls: sharedUtil.CastSliceP(post.Medias, ToMediaDTO),
  })

  for _, edited := range post.EditedPost {
    childPost := dto.ChildPostResponseDTO{
      Content:   edited.Content,
      CreatedAt: edited.CreatedAt,
      Tags:      sharedUtil.CastSliceP(edited.Tags, ToTaggedUserDTO),
      MediaUrls: sharedUtil.CastSliceP(edited.Medias, ToMediaDTO),
    }

    children = append(children, childPost)
  }

  return dto.EditedPostResponseDTO{
    PostId:      post.Id,
    CreatorId:   post.CreatorId,
    EditedPosts: children,
  }
}
