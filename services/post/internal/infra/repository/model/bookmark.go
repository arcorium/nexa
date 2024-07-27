package model

import "github.com/uptrace/bun"

type BookmarkPost struct {
  bun.BaseModel `bun:"table:bookmark_posts,alias:bookmark_post"`

  UserId string `bun:",type:uuid,nullzero,pk"`
  PostId string `bun:",type:uuid,nullzero,pk"`

  Post *Post `bun:"rel:belongs-to,join:post_id=id,on_delete:CASCADE"`
}
