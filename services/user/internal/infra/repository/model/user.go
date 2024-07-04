package model

import (
  "database/sql"
  "github.com/uptrace/bun"
  domain "nexa/services/user/internal/domain/entity"
  "nexa/shared/types"
  "nexa/shared/util/repo"
  "nexa/shared/variadic"
  "time"
)

type UserMapOption = repo.DataAccessModelMapOption[*domain.User, *User]

func FromUserDomain(domain *domain.User, opts ...UserMapOption) User {
  // Time based field should be handled individually
  usr := User{
    Id:         domain.Id.Underlying().String(),
    Username:   domain.Username,
    Email:      domain.Email.Underlying(),
    Password:   domain.Password.Underlying(),
    IsVerified: domain.SqlIsVerified(),
  }

  variadic.New(opts...).
    DoAll(repo.MapOptionFunc(domain, &usr))

  return usr
}

type User struct {
  bun.BaseModel `bun:"table:users"`

  Id         string       `bun:",type:uuid,pk"`
  Username   string       `bun:",nullzero,notnull,unique"`
  Email      string       `bun:",nullzero,notnull,unique"`
  Password   string       `bun:",nullzero,notnull"`
  IsVerified sql.NullBool `bun:",notnull,default:false"`

  BannedUntil time.Time `bun:",nullzero"`
  DeletedAt   time.Time `bun:",nullzero"`
  CreatedAt   time.Time `bun:",nullzero,notnull"`
  UpdatedAt   time.Time `bun:",nullzero"`

  Profile *Profile `bun:"rel:has-one,join=id:user_id"`
}

func (u *User) ToDomain() (domain.User, error) {
  id, err := types.IdFromString(u.Id)
  if err != nil {
    return domain.User{}, err
  }

  email, err := types.EmailFromString(u.Email)
  if err != nil {
    return domain.User{}, err
  }

  var profile *domain.Profile = nil
  if u.Profile != nil {
    obj, err := u.Profile.ToDomain()
    if err != nil {
      return domain.User{}, err
    }
    profile = &obj // dangling and allocated on heap
  }

  return domain.User{
    Id:          id,
    Username:    u.Username,
    Email:       email,
    Password:    types.HashedPassword(u.Password),
    IsVerified:  &u.IsVerified.Bool,
    BannedUntil: u.BannedUntil,
    IsDeleted:   !u.DeletedAt.IsZero(),
    Profile:     profile,
  }, nil
}
