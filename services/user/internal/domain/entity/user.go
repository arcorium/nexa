package entity

import (
  "database/sql"
  "nexa/shared/types"
  "time"
)

type User struct {
  Id       types.Id
  Username string
  Email    types.Email
  Password types.HashedPassword

  IsVerified  *bool
  IsDeleted   bool
  BannedUntil time.Time

  Profile *Profile
}

func (u *User) SqlIsVerified() sql.NullBool {
  return sql.NullBool{
    Bool:  *u.IsVerified,
    Valid: u.IsVerified != nil,
  }
}

func (u *User) ValidatePassword(password string) error {
  return u.Password.Equal(password)
}
