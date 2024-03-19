package entity

import (
	"github.com/uptrace/bun"
	"nexa/shared/file"
	"nexa/shared/types"
	"time"
)

type Profile struct {
	bun.BaseModel `bun:"table:profiles"`

	Id        types.Id        `bun:",nullzero,type:uuid,pk"`
	UserId    types.Id        `bun:",nullzero,notnull,unique,type:uuid"`
	FirstName string          `bun:",notnull"`
	LastName  string          `bun:",nullzero"`
	PhotoURL  file.RemotePath `bun:",nullzero"`
	Bio       string          `bun:",nullzero"`

	UpdatedAt time.Time `bun:",nullzero,notnull"`

	User *User `bun:"rel:belongs-to,join:user_id=id,on_delete:CASCADE"`
}
