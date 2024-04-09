package entity

import (
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

  Profile *Profile
}

func (u *User) ValidatePassword(password string) error {
  return u.Password.Equal(password)
}
