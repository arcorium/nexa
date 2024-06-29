package entity

import (
  sharedJwt "nexa/shared/jwt"
  "nexa/shared/types"
  sharedUtil "nexa/shared/util"
)

type Role struct {
  Id          types.Id
  Name        string
  Description string

  Permissions []Permission
}

func (r *Role) ToJWT() sharedJwt.Role {
  return sharedJwt.Role{
    Id:   r.Id.String(),
    Role: r.Name,
    Permissions: sharedUtil.CastSliceP(r.Permissions, func(perm *Permission) string {
      return perm.Code
    }),
  }
}
