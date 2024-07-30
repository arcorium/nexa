package entity

import (
  sharedErr "github.com/arcorium/nexa/shared/errors"
  "github.com/arcorium/nexa/shared/types"
  "time"
)

func NewVisibility(val uint8) (Visibility, error) {
  vis := Visibility(val)
  if !vis.Valid() {
    return vis, sharedErr.ErrEnumOutOfBounds
  }
  return vis, nil
}

type Visibility uint8

func (v Visibility) Underlying() uint8 {
  return uint8(v)
}

func (v Visibility) Valid() bool {
  return v.Underlying() < VisibilityUnknown.Underlying()
}

// Use hierarchy, OnlyMe visibility will also able to get other visibility posts
const (
  VisibilityPublic Visibility = iota
  VisibilityFollower
  VisibilityOnlyMe
  VisibilityUnknown
)

type Post struct {
  Id         types.Id
  Parent     *Post
  CreatorId  types.Id
  Content    string
  Visibility Visibility

  TotalLikes    int64 // negative numbers means the client should take it separately
  TotalDislikes int64
  TotalComments int64

  TaggedUserId []types.Id
  MediaUrls    []string
  LastEditedAt time.Time
  CreatedAt    time.Time
}

type Feed struct {
  UserId types.Id
  Posts  []Post
}
