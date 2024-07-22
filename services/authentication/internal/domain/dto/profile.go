package dto

import (
  "github.com/arcorium/nexa/shared/types"
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
