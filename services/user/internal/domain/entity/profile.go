package entity

import (
  "nexa/shared/types"
)

type Profile struct {
  Id        types.Id
  UserId    types.Id
  FirstName string
  LastName  string
  Bio       string
  PhotoId   types.Id
  PhotoURL  types.FilePath
}

func (p *Profile) HasAvatar() bool {
  return len(p.PhotoURL.Path()) != 0 && !p.PhotoId.Eq(types.NullId())
}

// PatchedProfile used for patchable profile which contains all nullable data
type PatchedProfile struct {
  Id        types.Id
  FirstName string
  LastName  types.NullableString
  Bio       types.NullableString
  PhotoId   types.NullableId
  PhotoURL  types.NullablePath
}
