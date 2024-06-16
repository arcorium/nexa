package model

import (
	"github.com/uptrace/bun"
	"time"
)

type UserRoles struct {
	bun.BaseModel `bun:"table:user_roles"`

	UserId string `bun:",type:uuid,pk"`
	RoleId string `bun:",type:uuid,pk"`

	User *User `bun:"rel:belongs-to,join:user_id=id,on_delete:CASCADE"`

	CreatedAt time.Time `bun:",nullzero,notnull"`
}
