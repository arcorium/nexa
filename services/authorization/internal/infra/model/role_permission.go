package model

import (
	"github.com/uptrace/bun"
	"time"
)

type RolePermission struct {
	bun.BaseModel `bun:"table:role_permissions"`

	RoleId       string `bun:",nullzero,type:uuid,pk"`
	PermissionId string `bun:",nullzero,type:uuid,pk"`

	CreatedAt time.Time `bun:",nullzero,notnull"`

	Role       *Role       `bun:"rel:belongs-to,join:role_id=id,on_delete:CASCADE"`
	Permission *Permission `bun:"rel:belongs-to,join:permission_id=id,on_delete:CASCADE"`
}
