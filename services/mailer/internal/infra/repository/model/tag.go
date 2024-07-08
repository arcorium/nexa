package model

import (
  "github.com/uptrace/bun"
  "nexa/services/mailer/internal/domain/entity"
  "nexa/shared/types"
  "nexa/shared/util/repo"
  "nexa/shared/variadic"
  "time"
)

type TagMapOption = repo.DataAccessModelMapOption[*entity.Tag, *Tag]

type PatchedTagMapOption = repo.DataAccessModelMapOption[*entity.PatchedTag, *Tag]

func FromPatchedTagDomain(domain *entity.PatchedTag, opts ...PatchedTagMapOption) Tag {
  obj := Tag{
    Id:          domain.Id.String(),
    Name:        domain.Name,
    Description: domain.Description.Value(),
  }

  variadic.New(opts...).DoAll(repo.MapOptionFunc(domain, &obj))

  return obj
}

func FromTagDomain(domain *entity.Tag, opts ...TagMapOption) Tag {
  obj := Tag{
    Id:          domain.Id.String(),
    Name:        domain.Name,
    Description: &domain.Description,
  }

  variadic.New(opts...).DoAll(repo.MapOptionFunc(domain, &obj))

  return obj
}

type Tag struct {
  bun.BaseModel `bun:"table:tags"`

  Id          string  `bun:",type:uuid,pk,nullzero"`
  Name        string  `bun:",unique,notnull,nullzero"`
  Description *string `bun:","`

  CreatedAt time.Time `bun:",notnull,nullzero"`
  UpdatedAt time.Time `bun:",nullzero"`
}

func (t *Tag) ToDomain() (entity.Tag, error) {
  id, err := types.IdFromString(t.Id)
  if err != nil {
    return entity.Tag{}, err
  }

  return entity.Tag{
    Id:          id,
    Name:        t.Name,
    Description: types.OnNil(t.Description, ""),
  }, nil
}
