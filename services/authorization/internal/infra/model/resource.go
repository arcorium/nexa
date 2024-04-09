package model

import (
  "github.com/uptrace/bun"
  "nexa/services/authorization/internal/domain/entity"
  "nexa/shared/types"
  "nexa/shared/util/repo"
  "nexa/shared/variadic"
  "time"
)

type ResourceMapOption = repo.DataAccessModelMapOption[*entity.Resource, *Resource]

func FromResourceDomain(domain *entity.Resource, opts ...ResourceMapOption) Resource {
  resource := Resource{
    Id:          domain.Id.Underlying().String(),
    Name:        domain.Name,
    Description: domain.Description,
  }

  variadic.New(opts...).
    DoAll(repo.MapOptionFunc(domain, &resource))

  return resource
}

type Resource struct {
  bun.BaseModel `bun:"table:resources"`

  Id          string `bun:",nullzero,type:uuid,pk"`
  Name        string `bun:",nullzero,unique"`
  Description string `bun:",nullzero"`

  UpdatedAt time.Time `bun:",nullzero"`
  CreatedAt time.Time `bun:",nullzero,notnull"`
}

func (r *Resource) ToDomain() entity.Resource {
  return entity.Resource{
    Id:          types.IdFromString(r.Id),
    Name:        r.Name,
    Description: r.Description,
  }
}
