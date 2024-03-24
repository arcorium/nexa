package entity

import (
	"nexa/shared/types"
)

type Profile struct {
	Id        types.Id
	FirstName string
	LastName  string
	Bio       string
	PhotoURL  types.FilePath
}
