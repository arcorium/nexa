package model

import (
  "github.com/arcorium/nexa/shared/types"
  "github.com/arcorium/nexa/shared/util/repo"
  "github.com/arcorium/nexa/shared/variadic"
  "github.com/uptrace/bun"
  "nexa/services/comment/internal/domain/entity"
  "time"
)

type CommentMapOption = repo.DataAccessModelMapOption[*entity.Comment, *Comment]

func FromCommentDomain(ent *entity.Comment, opts ...CommentMapOption) Comment {
  parentId := types.NullId()
  if ent.IsReply() {
    parentId = ent.Parent.Id
  }

  comment := Comment{
    Id:        ent.Id.String(),
    PostId:    ent.PostId.String(),
    ParentId:  parentId.String(),
    UserId:    ent.UserId.String(),
    Content:   ent.Content,
    CreatedAt: time.Now(),
  }

  variadic.New(opts...).
    DoAll(repo.MapOptionFunc(ent, &comment))

  return comment
}

type Comment struct {
  bun.BaseModel `bun:"table:comments"`

  Id       string `bun:",type:uuid,nullzero,pk"`
  PostId   string `bun:",type:uuid,nullzero,notnull"`
  ParentId string `bun:",type:uuid,nullzero"`
  UserId   string `bun:"type:uuid,nullzero,notnull"`
  Content  string `bun:",type:text,nullzero"`

  UpdatedAt time.Time `bun:",nullzero"`
  CreatedAt time.Time `bun:",nullzero,notnull"`

  Parent *Comment `bun:"rel:belongs-to,join:parent_id=id,on_delete:CASCADE"`

  TotalReply uint64 `bun:",scanonly"`
}

func (c *Comment) ToDomain() (entity.Comment, error) {
  id, err := types.IdFromString(c.Id)
  if err != nil {
    return entity.Comment{}, err
  }

  postId, err := types.IdFromString(c.PostId)
  if err != nil {
    return entity.Comment{}, err
  }

  userId, err := types.IdFromString(c.UserId)
  if err != nil {
    return entity.Comment{}, err
  }

  var parent *entity.Comment
  if c.Parent != nil {
    temp, err := c.Parent.ToDomain()
    if err != nil {
      return entity.Comment{}, err
    }
    parent = &temp
  } else {
    // Set parent with only the id, that will be used by domain model to restructure
    if len(c.ParentId) != 0 {
      parentId, err := types.IdFromString(c.ParentId)
      if err != nil {
        return entity.Comment{}, err
      }

      parent = &entity.Comment{
        Id: parentId,
      }
    }
  }

  return entity.Comment{
    Id:        id,
    PostId:    postId,
    UserId:    userId,
    Content:   c.Content,
    UpdatedAt: c.UpdatedAt,
    CreatedAt: c.CreatedAt,
    Parent:    parent,
  }, nil
}
