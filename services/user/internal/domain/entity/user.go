package entity

import (
  "database/sql"
  "github.com/arcorium/nexa/shared/types"
  "time"
)

type User struct {
  Id         types.Id
  Username   string
  Email      types.Email
  Password   types.HashedPassword
  IsVerified bool

  DeletedAt   time.Time
  BannedUntil time.Time
  CreatedAt   time.Time

  Profile *Profile
}

func (u *User) ValidatePassword(password string) error {
  return u.Password.Equal(password)
}

type PatchedUser struct {
  Id         types.Id
  Username   string
  Email      types.Email
  Password   types.HashedPassword
  IsVerified types.NullableBool
  IsDelete   types.NullableBool

  BannedDuration types.Nullable[time.Duration] // Set to 0 to unban
}

func (u *PatchedUser) SqlIsVerified() sql.NullBool {
  if !u.IsVerified.HasValue() {
    return sql.NullBool{
      Bool:  false,
      Valid: false,
    }
  }
  return sql.NullBool{
    Bool:  u.IsVerified.RawValue(),
    Valid: true,
  }
}
