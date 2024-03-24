package model

import (
	"github.com/uptrace/bun"
	"time"
)

type Verification struct {
	bun.BaseModel `bun:"table:verifications"`

	Token          string    `bun:",nullzero,pk"`
	UserId         string    `bun:",nullzero,notnull,type:uuid,unique:verif_usage_idx"`
	UsageId        string    `bun:",nullzero,notnull,type:uuid,unique:verif_usage_idx"`
	ExpirationTime time.Time `bun:",nullzero,notnull"`

	UpdatedAt time.Time `bun:",nullzero"`
	CreatedAt time.Time `bun:",nullzero,notnull"`

	VerificationUsage *VerificationUsage `bun:"rel:belongs-to,join:usage_id=id"`
}
