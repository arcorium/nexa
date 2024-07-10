package entity

import (
  sharedJwt "github.com/arcorium/nexa/shared/jwt"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
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
      return perm.Encode()
    }),
  }
}

// PatchedRole used as patching field role that also handle nullable field
type PatchedRole struct {
  Id          types.Id
  Name        string
  Description types.NullableString
}
