package model

import (
  "github.com/arcorium/nexa/shared/types"
  "github.com/arcorium/nexa/shared/util/repo"
  "github.com/arcorium/nexa/shared/variadic"
  "github.com/uptrace/bun"
  "nexa/services/authorization/internal/domain/entity"
  "time"
)

type PermissionMapOption = repo.DataAccessModelMapOption[*entity.Permission, *Permission]

func FromPermissionDomain(ent *entity.Permission, opts ...PermissionMapOption) Permission {
  permission := Permission{
    Id:        ent.Id.String(),
    Resource:  ent.Resource,
    Action:    ent.Action,
    CreatedAt: ent.CreatedAt,
  }

  variadic.New(opts...).
    DoAll(repo.MapOptionFunc(ent, &permission))

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
