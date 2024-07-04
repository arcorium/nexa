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
    Id:       domain.Id.String(),
    Resource: domain.Resource,
    Action:   domain.Action,
  }

  variadic.New(opts...).
    DoAll(repo.MapOptionFunc(domain, &permission))

  return permission
}

type Permission struct {
  bun.BaseModel `bun:"permissions"`

  Id       string `bun:",nullzero,type:uuid,pk"`
  Resource string `bun:",nullzero,notnull,unique:res_action_idx"`
  Action   string `bun:",nullzero,notnull,unique:res_action_idx"`

  CreatedAt time.Time `bun:",nullzero,notnull"`
}

func (p *Permission) ToDomain() (entity.Permission, error) {
  permId, err := types.IdFromString(p.Id)
  if err != nil {
    return entity.Permission{}, err
  }

  return entity.Permission{
    Id:        permId,
    Resource:  p.Resource,
    Action:    p.Action,
    CreatedAt: p.CreatedAt,
  }, nil
}
