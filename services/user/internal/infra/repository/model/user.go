package model

import (
  "database/sql"
  "github.com/arcorium/nexa/shared/types"
  "github.com/arcorium/nexa/shared/util/repo"
  "github.com/arcorium/nexa/shared/variadic"
  "github.com/uptrace/bun"
  "nexa/services/user/internal/domain/entity"
  "time"
)

type UserMapOption = repo.DataAccessModelMapOption[*entity.User, *User]
type PatchedUserMapOption = repo.DataAccessModelMapOption[*entity.PatchedUser, *User]

func FromPatchedUserDomain(ent *entity.PatchedUser, opts ...PatchedUserMapOption) User {
  //var deletedAt *time.Time = nil
  var deletedAt = sql.NullTime{}
  // Wanted to change
  if ent.IsDelete.HasValue() {
    if ent.IsDelete.RawValue() {
      // Delete
      deletedAt = sql.NullTime{
        Time:  time.Now(),
        Valid: true,
      }
    } else {
      // Undelete
      deletedAt = sql.NullTime{
        Time:  time.Time{},
        Valid: true,
      }
    }
  }

  //var bannedUntil *time.Time = nil
  var bannedUntil = sql.NullTime{}
  // Wanted to change
  if ent.BannedDuration.HasValue() {
    if int(ent.BannedDuration.RawValue().Seconds()) == 0 {
      // Unban
      bannedUntil = sql.NullTime{Valid: true}
    } else {
      // Banned
      now := time.Now().Add(ent.BannedDuration.RawValue())
      //bannedUntil = &now
      bannedUntil = sql.NullTime{
        Time:  now,
        Valid: true,
      }
    }
  }
  usr := User{
    Id:          ent.Id.String(),
    Username:    ent.Username,
    Email:       ent.Email.String(),
    Password:    ent.Password.String(),
    IsVerified:  ent.SqlIsVerified(),
    BannedUntil: bannedUntil,
    DeletedAt:   deletedAt,
    UpdatedAt:   time.Now(),
  }

  variadic.New(opts...).
    DoAll(repo.MapOptionFunc(ent, &usr))

  return usr
}

func FromUserDomain(ent *entity.User, opts ...UserMapOption) User {
  usr := User{
    Id:       ent.Id.String(),
    Username: ent.Username,
    Email:    ent.Email.Underlying(),
    Password: ent.Password.Underlying(),
    IsVerified: sql.NullBool{
      Bool:  ent.IsVerified,
      Valid: true,
    },
    BannedUntil: sql.NullTime{
      Time:  ent.BannedUntil,
      Valid: true,
    },
    DeletedAt: sql.NullTime{
      Time:  ent.DeletedAt,
      Valid: true,
    },
    CreatedAt: ent.CreatedAt,
  }

  variadic.New(opts...).
    DoAll(repo.MapOptionFunc(ent, &usr))

  return usr
}

type User struct {
  bun.BaseModel `bun:"table:users,alias:u"`

  Id         string       `bun:",type:uuid,pk"`
  Username   string       `bun:",nullzero,notnull,unique"`
  Email      string       `bun:",nullzero,notnull,unique"`
  Password   string       `bun:",nullzero,notnull"`
  IsVerified sql.NullBool `bun:",notnull,default:false"`

  BannedUntil sql.NullTime `bun:","` // Use nullable type to be able to unban
  DeletedAt   sql.NullTime `bun:","` // Use nullable type to be able to undelete
  UpdatedAt   time.Time    `bun:",nullzero"`
  CreatedAt   time.Time    `bun:",nullzero,notnull"`

  Profile *Profile `bun:"rel:has-one,join:id=user_id"`
}

func (u *User) ToDomain() (entity.User, error) {
  id, err := types.IdFromString(u.Id)
  if err != nil {
    return entity.User{}, err
  }

  email, err := types.EmailFromString(u.Email)
  if err != nil {
    return entity.User{}, err
  }

  var profile *entity.Profile = nil
  if u.Profile != nil {
    obj, err := u.Profile.ToDomain()
    if err != nil {
      return entity.User{}, err
    }
    profile = &obj // dangling and allocated on heap
  }

  return entity.User{
    Id:          id,
    Username:    u.Username,
    Email:       email,
    Password:    types.HashedPassword(u.Password),
    IsVerified:  u.IsVerified.Bool,
    DeletedAt:   u.DeletedAt.Time,
    BannedUntil: u.BannedUntil.Time,
    CreatedAt:   u.CreatedAt,
    Profile:     profile,
  }, nil
}
