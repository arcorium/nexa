package entity

import (
	"fmt"
	"nexa/shared/types"
)

type Permission struct {
	Id       types.Id
	Resource Resource
	Action   Action
}

func (p *Permission) String() string {
	return fmt.Sprintf("%s:%s", p.Resource.Name, p.Action.Name)
}
