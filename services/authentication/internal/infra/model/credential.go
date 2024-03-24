package model

import (
	"github.com/uptrace/bun"
	"nexa/services/authentication/shared/domain/valueobject"
	"time"
)

type Credential struct {
	bun.BaseModel `bun:"table:credentials"`

	UserId        string             `bun:",nullzero,type:uuid,pk"`
	AccessTokenId string             `bun:"access_id,nullzero,notnull,type:uuid,pk"`
	Device        valueobject.Device `bun:",nullzero,embed:device_"`
	Token         string             `bun:",nullzero,notnull"`

	UpdatedAt time.Time `bun:",nullzero"`
	CreatedAt time.Time `bun:",nullzero,notnull"`
}
