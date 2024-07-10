package model

import (
  "github.com/arcorium/nexa/shared/types"
  "github.com/arcorium/nexa/shared/util"
  "github.com/arcorium/nexa/shared/util/repo"
  "github.com/arcorium/nexa/shared/variadic"
  "github.com/uptrace/bun"
  "nexa/services/authorization/internal/domain/entity"
  "time"
)

type RoleMapOption = repo.DataAccessModelMapOption[*entity.Role, *Role]

type PatchedRoleMapOption = repo.DataAccessModelMapOption[*entity.PatchedRole, *Role]

func FromPatchedRoleDomain(ent *entity.PatchedRole, opts ...PatchedRoleMapOption) Role {
  role := Role{
    Id:          ent.Id.Underlying().String(),
    Name:        ent.Name,
    Description: ent.Description.Value(),
  }

  variadic.New(opts...).
    DoAll(repo.MapOptionFunc(ent, &role))

  return role
}

func FromRoleDomain(ent *entity.Role, opts ...RoleMapOption) Role {
  role := Role{
    Id:          ent.Id.Underlying().String(),
    Name:        ent.Name,
    Description: &ent.Description,
  }

  variadic.New(opts...).
    DoAll(repo.MapOptionFunc(ent, &role))

  return role
}

type Role struct {
  bun.BaseModel `bun:"table:roles"`

  Id          string  `bun:",nullzero,type:uuid,pk"`
  Name        string  `bun:",nullzero,notnull,unique"`
  Description *string `bun:","`

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

  return entity.Role{
    Id:          roleId,
    Name:        r.Name,
    Description: types.OnNil(r.Description, ""),
    Permissions: perms,
  }, nil
}
