package entity

import "nexa/shared/types"

type Tag struct {
  Id          types.Id
  Name        string
  Description string
}

type PatchedTag struct {
  Id          types.Id
  Name        string
  Description types.NullableString
}
