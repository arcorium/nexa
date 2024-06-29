package dto

import (
  sharedJwt "nexa/shared/jwt"
  "nexa/shared/types"
  sharedUtil "nexa/shared/util"
)

type Permission struct {
  Id   types.Id
  Code string
}

type RoleResponseDTO struct {
  Id          types.Id
  Role        string
  Permissions []Permission
}

func (r *RoleResponseDTO) ToJWT() sharedJwt.Role {
  return sharedJwt.Role{
    Id:   r.Id.Underlying().String(),
    Role: r.Role,
    Permissions: sharedUtil.CastSliceP(r.Permissions, func(from *Permission) string {
      return from.Code
    }),
  }
}
