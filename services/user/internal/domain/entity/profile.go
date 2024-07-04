package entity

import (
  "nexa/shared/types"
)

type Profile struct {
  Id        types.Id
  FirstName string
  LastName  string
  Bio       string
  PhotoId   types.Id
  PhotoURL  types.FilePath
}

func (p *Profile) HasAvatar() bool {
  return len(p.PhotoURL.Underlying()) == 0 || p.Id.Eq(types.NullId())
}
