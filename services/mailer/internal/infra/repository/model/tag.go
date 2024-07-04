package model

import (
  "github.com/uptrace/bun"
  domain "nexa/services/mailer/internal/domain/entity"
  "nexa/shared/types"
  "nexa/shared/util/repo"
  "nexa/shared/variadic"
  "time"
)

type TagMapOption = repo.DataAccessModelMapOption[*domain.Tag, *Tag]

func FromTagDomain(domain *domain.Tag, opts ...TagMapOption) Tag {
  obj := Tag{
    Id:          domain.Id.String(),
    Name:        domain.Name,
    Description: domain.Description,
  }

  variadic.New(opts...).DoAll(repo.MapOptionFunc(domain, &obj))

  return obj
}

type Tag struct {
  bun.BaseModel `bun:"table:tags"`

  Id          string `bun:",type:uuid,pk,nullzero"`
  Name        string `bun:",unique,notnull,nullzero"`
  Description string `bun:",nullzero"`

  CreatedAt time.Time `bun:",notnull,nullzero"`
  UpdatedAt time.Time `bun:",nullzero"`
}

func (t *Tag) ToDomain() (domain.Tag, error) {
  id, err := types.IdFromString(t.Id)
  if err != nil {
    return domain.Tag{}, err
  }

  return domain.Tag{
    Id:          id,
    Name:        t.Name,
    Description: t.Description,
  }, nil
}

var DefaultTags = []Tag{
  {
    Id:          types.NewId2().String(),
    Name:        "Email Validation",
    Description: "Email Validation",
    CreatedAt:   time.Now(),
  },
  {
    Id:          types.NewId2().String(),
    Name:        "Forgot Password",
    Description: "Forgot Password",
    CreatedAt:   time.Now(),
  },
}
