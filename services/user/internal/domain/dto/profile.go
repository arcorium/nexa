package dto

import (
  "nexa/services/user/internal/domain/entity"
  "nexa/shared/types"
  "nexa/shared/wrapper"
)

type ProfileResponseDTO struct {
  //UserId    string `json:"user_id"`
  FirstName string
  LastName  string
  Bio       string
  PhotoURL  string
}

//type ProfileCreateInput struct {
//	UserId    string `json:"user_id"`
//	FirstName string `json:"first_name"`
//	LastName  string `json:"last_name"`
//	Bio       string `json:"bio"`
//}

type ProfileUpdateDTO struct {
  UserId    types.Id
  FirstName wrapper.Nullable[string]
  LastName  wrapper.Nullable[string]
  Bio       wrapper.Nullable[string]
}

func (p *ProfileUpdateDTO) ToDomain() entity.Profile {
  profile := entity.Profile{
    Id: p.UserId,
  }

  wrapper.SetOnNonNull(&profile.FirstName, p.FirstName)
  wrapper.SetOnNonNull(&profile.LastName, p.LastName)
  wrapper.SetOnNonNull(&profile.Bio, p.Bio)

  return profile
}

type ProfileAvatarUpdateDTO struct {
  UserId   types.Id
  Filename string `validate:"required"`
  Bytes    []byte `validate:"required"`
}
