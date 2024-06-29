package model

import (
  "github.com/uptrace/bun"
  "nexa/services/authorization/internal/domain/entity"
  "nexa/shared/types"
  "nexa/shared/util/repo"
  "nexa/shared/variadic"
  "time"
)

type PermissionMapOption = repo.DataAccessModelMapOption[*entity.Permission, *Permission]

func FromPermissionDomain(domain *entity.Permission, opts ...PermissionMapOption) Permission {
  permission := Permission{
    Id:   domain.Id.Underlying().String(),
    Code: domain.Code,
  }

  variadic.New(opts...).
    DoAll(repo.MapOptionFunc(domain, &permission))

  return permission
}

type Permission struct {
  bun.BaseModel `bun:"permissions"`

  Id   string `bun:",nullzero,type:uuid,pk"`
  Code string `bun:",nullzero,notnull,unique"`

  CreatedAt time.Time `bun:",nullzero,notnull"`
}

func (p *Permission) ToDomain() (entity.Permission, error) {
  permId, err := types.IdFromString(p.Id)
  if err != nil {
    return entity.Permission{}, err
  }

  return entity.Permission{
    Id:        permId,
    Code:      p.Code,
    CreatedAt: p.CreatedAt,
  }, nil
}
