package dto

import (
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/authentication/internal/domain/entity"
)

type ProfileResponseDTO struct {
  Id        types.Id
  FirstName string
  LastName  string
  Bio       string
  PhotoURL  types.FilePath
}

//type ProfileCreateInput struct {
//	UserId    string
//	FirstName string
//	LastName  string
//	Bio       string
//}

type ProfileUpdateDTO struct {
  Id        types.Id
  FirstName types.NullableString
  LastName  types.NullableString
  Bio       types.NullableString
}

func (p *ProfileUpdateDTO) ToDomain() entity.PatchedProfile {
  profile := entity.PatchedProfile{
    Id:       p.Id,
    LastName: p.LastName,
    Bio:      p.Bio,
  }

  types.SetOnNonNull(&profile.FirstName, p.FirstName)

  return profile
}

type ProfileAvatarUpdateDTO struct {
  Id       types.Id
  Filename string `validate:"required"`
  Bytes    []byte `validate:"required"`
}
