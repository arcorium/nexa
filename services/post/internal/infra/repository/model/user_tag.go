package model

import (
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/post/internal/domain/entity"
)

type UserTag struct {
  Id        uint64 `bun:",autoincrement,pk"`
  VersionId string `bun:"type:uuid,notnull,nullzero"`
  UserId    string `bun:"type:uuid,notnull,nullzero"`

  PostVersion *PostVersion `bun:"rel:belongs-to,join:version_id=id,on_delete:CASCADE"`
}

func (t *UserTag) ToDomain() (entity.TaggedUser, error) {
  userId, err := types.IdFromString(t.UserId)
  if err != nil {
    return entity.TaggedUser{}, err
  }

  return entity.TaggedUser{
    Id: userId,
    //Name: "", // Set by service
  }, nil
}
