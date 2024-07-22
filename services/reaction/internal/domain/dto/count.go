package dto

import "github.com/arcorium/nexa/shared/types"

type CountResponseDTO struct {
  Like    uint64
  Dislike uint64
}

type PostCount struct {
  PostId types.Id
  CountResponseDTO
}

type CommentCount struct {
  CommentId types.Id
  CountResponseDTO
}
