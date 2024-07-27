package model

import (
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  "github.com/arcorium/nexa/shared/util/repo"
  "github.com/uptrace/bun"
  "nexa/services/post/internal/domain/entity"
  "time"
)

type PostVersion struct {
  bun.BaseModel `bun:"table:post_versions"`

  Id      string `bun:",type:uuid,pk"`
  PostId  string `bun:",type:uuid"`
  Content string `bun:",type:text"`
  //IsVisible sql.NullBool `bun:",default:true"`

  CreatedAt time.Time `bun:",nullzero,notnull"`

  Post     *Post     `bun:"rel:belongs-to,join:post_id=id,on_delete:CASCADE"`
  Medias   []Media   `bun:"rel:has-many,join:id=version_id"`
  UserTags []UserTag `bun:"rel:has-many,join:id=version_id"`
}

func (p *PostVersion) ToDomain() (entity.ChildPost, error) {
  verId, err := types.IdFromString(p.Id)
  if err != nil {
    return entity.ChildPost{}, err
  }

  userTags, ierr := sharedUtil.CastSliceErrsP(p.UserTags, repo.ToDomainErr[*UserTag, entity.TaggedUser])
  if !ierr.IsNil() && !ierr.IsEmptySlice() {
    return entity.ChildPost{}, ierr
  }

  medias, ierr := sharedUtil.CastSliceErrsP(p.Medias, repo.ToDomainErr[*Media, entity.Media])
  if !ierr.IsNil() && !ierr.IsEmptySlice() {
    return entity.ChildPost{}, ierr
  }

  return entity.ChildPost{
    Id:        verId,
    Content:   p.Content,
    CreatedAt: p.CreatedAt,
    Tags:      userTags,
    Medias:    medias,
  }, nil
}
