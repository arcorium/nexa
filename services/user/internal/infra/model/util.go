package model

import (
	"github.com/uptrace/bun"
	"nexa/shared/util"
)

func RegisterBunModels(db *bun.DB) {
	db.RegisterModel(util.Nil[User](), util.Nil[Profile]())
}
