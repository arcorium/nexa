package entity

import (
  "database/sql"
  "github.com/arcorium/nexa/shared/types"
  "time"
)

func NewUser(username string, email types.Email, password types.Password) (User, error) {
  userId, err := types.NewId()
  if err != nil {
    return User{}, err
  }

  hashedPass, err := password.Hash()
  if err != nil {
    return User{}, err
  }

  return User{
    Id:         userId,
    Username:   username,
    Email:      email,
    Password:   hashedPass,
    IsVerified: false,
  }, nil
}

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

func (u *User) ValidatePassword(password types.Password) error {
  return u.Password.Eq(password)
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
