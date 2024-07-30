package dto

import (
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/feed/internal/domain/entity"
  "time"
)

type PostResponseDTO struct {
  Id         types.Id
  Parent     *PostResponseDTO
  CreatorId  types.Id
  Content    string
  Visibility entity.Visibility

  TotalLikes    int64 // negative numbers means the client should take it separately
  TotalDislikes int64
  TotalComments int64

  TaggedUserIds []types.Id
  MediaUrls     []string
  LastEditedAt  time.Time
  CreatedAt     time.Time
}

type GetUsersPostResponseDTO struct {
  Posts []PostResponseDTO
}

type PostReactionCountDTO struct {
  TotalLikes    uint64
  TotalDislikes uint64
}
