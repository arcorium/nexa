package entity

import (
	"errors"
	"nexa/services/access_control/shared/domain/entity"
	"nexa/shared/types"
	"time"
)

type User struct {
	Id       types.Id
	Username string
	Email    types.Email
	Password types.Password

	IsVerified  bool
	IsDeleted   bool
	BannedUntil time.Time

	Roles []entity.Role

	toggleDelete bool
}

func (u *User) ValidatePassword(password string) error {
	if !u.Password.Equal(password) {
		return ErrPasswordDifferent
	}
	return nil
}

func (u *User) ValidateEmail() error {
	return u.Email.Validate()
}

// Delete mark user should be deleted later
func (u *User) Delete() {
	if !u.IsDeleted {
		u.toggleDelete = true
	}
}

// ShouldDeleted used for repository either should delete the user or not
func (u *User) ShouldDeleted() bool {
	return u.toggleDelete
}

var ErrPasswordDifferent = errors.New("password different")
var ErrEmailMalformed = errors.New("bad email format")
