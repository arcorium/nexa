package model

import (
  "database/sql"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  "github.com/arcorium/nexa/shared/util/repo"
  "github.com/arcorium/nexa/shared/variadic"
  "github.com/uptrace/bun"
  "nexa/services/post/internal/domain/entity"
  "nexa/services/post/util/errs"
  "time"
)

type PostMapOption = repo.DataAccessModelMapOption[*entity.Post, *Post]

func FromPostDomain(ent *entity.Post, opts ...PostMapOption) Post {
  var parentId string
  if ent.ParentPost != nil {
    parentId = ent.ParentPost.Id.String()
  }

  post := Post{
    Id:        ent.Id.String(),
    ParentId:  parentId,
    CreatorId: ent.CreatorId.String(),
    Visibility: sql.NullInt64{
      Int64: int64(ent.Visibility.Underlying()),
      Valid: true,
    },
    CreatedAt: ent.CreatedAt,
  }

  version := PostVersion{
    Id:        types.MustCreateId().String(), // TODO: Move this
    PostId:    post.Id,
    Content:   ent.Content,
    CreatedAt: ent.CreatedAt,
  }

  userTags := sharedUtil.CastSlice(ent.Tags, func(tag entity.TaggedUser) UserTag {
    return UserTag{
      VersionId: version.Id,
      UserId:    tag.Id.String(),
    }
  })

  medias := sharedUtil.CastSlice(ent.Medias, func(media entity.Media) Media {
    return Media{
      VersionId: version.Id,
      FileId:    media.Id.String(),
    }
  })

  version.Medias = medias
  version.UserTags = userTags

  post.Versions = append(post.Versions, version)

  variadic.New(opts...).
    DoAll(repo.MapOptionFunc(ent, &post))

  return post
}

type Post struct {
  bun.BaseModel `bun:"table:posts"`

  Id         string        `bun:",type:uuid,pk"`
  ParentId   string        `bun:",type:uuid,nullzero"`
  CreatorId  string        `bun:",type:uuid,nullzero"`
  Visibility sql.NullInt64 `bun:",type:smallint"`

  CreatedAt time.Time `bun:",notnull,nullzero"`

  Versions []PostVersion `bun:"rel:has-many,join:id=post_id"`
  Parent   *Post         `bun:"rel:belongs-to,join:parent_id=id,on_delete:CASCADE"`
  //TotalShare uint64        `bun:",scanonly"`
}

func (p *Post) ToDomain() (entity.Post, error) {
  if len(p.Versions) == 0 {
    return entity.Post{}, errs.ErrPostWithNoVersion
  }
  vers := &p.Versions[0]

  postId, err := types.IdFromString(p.Id)
  if err != nil {
    return entity.Post{}, err
  }

  creatorId, err := types.IdFromString(p.CreatorId)
  if err != nil {
    return entity.Post{}, err
  }

  var children []entity.ChildPost
  if len(p.Versions) > 1 {
    childs, ierr := sharedUtil.CastSliceErrsP(p.Versions[1:], repo.ToDomainErr[*PostVersion, entity.ChildPost])
    if !ierr.IsNil() {
      return entity.Post{}, ierr
    }
    children = childs
  }

  child, err := vers.ToDomain()
  if err != nil {
    return entity.Post{}, err
  }

  visibility, err := entity.NewVisibility(uint8(p.Visibility.Int64))
  if err != nil {
    return entity.Post{}, err
  }

  // Parent
  var parent *entity.Post
  if p.Parent != nil {
    p, err := p.Parent.ToDomain()
    if err != nil {
      return entity.Post{}, err
    }
    parent = &p
  }

  return entity.Post{
    Id:         postId,
    ParentPost: parent,
    CreatorId:  creatorId,
    Content:    vers.Content,
    Visibility: visibility,
    //Shares:     p.TotalShare,
    LastEdited: vers.CreatedAt,
    CreatedAt:  p.CreatedAt,
    Tags:       child.Tags,
    Medias:     child.Medias,
    EditedPost: children,
  }, nil
}
