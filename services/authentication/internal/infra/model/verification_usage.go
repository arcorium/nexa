package model

import (
	"github.com/uptrace/bun"
	"time"
)

type VerificationUsage struct {
	bun.BaseModel `bun:"table:verification_usages"`

	Id   string `bun:",nullzero,type:uuid,pk"`
	Name string `bun:",nullzero,notnull,unique"`

	UpdatedAt time.Time `bun:",nullzero"`
	CreatedAt time.Time `bun:",nullzero,notnull"`
}
