package entity

import (
  "github.com/arcorium/nexa/shared/types"
)

func NewProfile(userId types.Id, firstName string) (Profile, error) {
  profileId, err := types.NewId()
  if err != nil {
    return Profile{}, err
  }

  return Profile{
    Id:        profileId,
    UserId:    userId,
    FirstName: firstName,
  }, nil
}

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
  //Id        types.Id
  UserId    types.Id
  FirstName string
  LastName  types.NullableString
  Bio       types.NullableString
  PhotoId   types.NullableId
  PhotoURL  types.NullablePath
}
