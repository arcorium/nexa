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
  return v.Underlying() >= VisibilityUnknown.Underlying()
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
  ParentPost *Post
  CreatorId  types.Id
  Content    string
  Visibility Visibility

  Likes    uint64
  Dislikes uint64
  Comments uint64
  Shares   uint64

  LastEdited time.Time
  CreatedAt  time.Time

  Tags       []TaggedUser
  Medias     []Media
  EditedPost []ChildPost
}

func (p *Post) IsShare() bool {
  return p.ParentPost != nil
}

type ChildPost struct {
  Id        types.Id
  Content   string
  CreatedAt time.Time

  Tags   []TaggedUser
  Medias []Media
}
