package model

import (
  "github.com/uptrace/bun"
  "nexa/services/authorization/internal/domain/entity"
  "nexa/shared/types"
  "nexa/shared/util"
  "nexa/shared/util/repo"
  "nexa/shared/variadic"
  "time"
)

type RoleMapOption = repo.DataAccessModelMapOption[*entity.Role, *Role]

func FromRoleDomain(domain *entity.Role, opts ...RoleMapOption) Role {
  role := Role{
    Id:          domain.Id.Underlying().String(),
    Name:        domain.Name,
    Description: domain.Description,
  }

  variadic.New(opts...).
    DoAll(repo.MapOptionFunc(domain, &role))

  return role
}

type Role struct {
  bun.BaseModel `bun:"table:roles"`

  Id          string `bun:",nullzero,type:uuid,pk"`
  Name        string `bun:",nullzero,unique"`
  Description string `bun:",nullzero"`

  UpdatedAt time.Time `bun:",nullzero"`
  CreatedAt time.Time `bun:",nullzero,notnull"`

  Permissions []Permission `bun:"m2m:role_permissions,join:Role=Permission"`
}

func (r *Role) ToDomain() (entity.Role, error) {
  roleId, err := types.IdFromString(r.Id)
  if err != nil {
    return entity.Role{}, err
  }

  perms, err := util.CastSliceErrsP(r.Permissions, func(from *Permission) (entity.Permission, error) {
    return from.ToDomain()
  })
  if err != nil {
    return entity.Role{}, err
  }

  return entity.Role{
    Id:          roleId,
    Name:        r.Name,
    Description: r.Description,
    Permissions: perms,
  }, nil
}
