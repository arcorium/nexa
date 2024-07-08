package model

import (
  "github.com/uptrace/bun"
  "time"
)

type UserRole struct {
  bun.BaseModel `bun:"table:user_roles"`

  UserId string `bun:",nullzero,type:uuid,pk"`
  RoleId string `bun:",nullzero,type:uuid,pk"`

  CreatedAt time.Time `bun:",nullzero,notnull"`

  Role *Role `bun:"rel:belongs-to,join:role_id=id,on_delete:CASCADE"`
}
