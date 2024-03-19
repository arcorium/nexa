package entity

import (
	"github.com/uptrace/bun"
	"nexa/shared/types"
	"time"
)

type Role struct {
	bun.BaseModel `bun:"table:roles"`

	Id          types.Id `bun:",nullzero,pk"`
	Name        string   `bun:",nullzero,notnull,unique"`
	Description string   `bun:",nullzero"`

	CreatedAt time.Time `bun:",nullzero,notnull"`
	UpdatedAt time.Time `bun:",nullzero"`
}
