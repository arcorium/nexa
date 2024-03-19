package entity

import (
	"errors"
	"github.com/uptrace/bun"
	"nexa/services/access_control/shared/domain/entity"
	"nexa/shared/types"
	"time"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	Id         types.Id         `bun:",type:uuid,pk"`
	Username   string           `bun:",nullzero,notnull,unique"`
	Email      types.Email      `bun:",nullzero,notnull,unique"`
	Password   types.HashString `bun:",nullzero,notnull"`
	IsVerified bool             `bun:",default:false"`
	Role       uint8            `bun:",nullzero,notnull"`

	BannedUntil time.Time `bun:",nullzero"`

	DeletedAt time.Time `bun:",nullzero"`
	CreatedAt time.Time `bun:",nullzero,notnull"`
	UpdatedAt time.Time `bun:",nullzero,notnull"`

	RoleDetail entity.Role `bun:",scanonly"`
}

func (u *User) ValidatePassword(password string) error {
	if !u.Password.Equal(password) {
		return ErrPasswordDifferent
	}
	return nil
}

func (u *User) ValidateEmail() error {
	if !u.Email.Validate() {
		return ErrEmailMalformed
	}
	return nil
}

var ErrPasswordDifferent = errors.New("password different")
var ErrEmailMalformed = errors.New("bad email format")
