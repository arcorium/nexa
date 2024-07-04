package entity

import (
  "fmt"
  "nexa/shared/types"
  "time"
)

type Permission struct {
  Id        types.Id
  Resource  string
  Action    string
  CreatedAt time.Time
}

func (p *Permission) Encode() string {
  return fmt.Sprintf("%s:%s", p.Resource, p.Action)
}
