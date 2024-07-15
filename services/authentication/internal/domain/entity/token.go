package entity

import "github.com/arcorium/nexa/shared/types"

type JWTToken struct {
  Id    types.Id
  Token string
}

type PairTokens struct {
  Access  JWTToken
  Refresh JWTToken
}
