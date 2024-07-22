package entity

import (
  "github.com/arcorium/nexa/shared/types"
  "time"
)

type Comment struct {
  Id        types.Id
  PostId    types.Id
  UserId    types.Id
  Content   string
  UpdatedAt time.Time
  CreatedAt time.Time

  Parent *Comment
}

func (c *Comment) IsReply() bool {
  return c.Parent != nil
}
