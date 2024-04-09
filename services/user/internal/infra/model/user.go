package model

import (
  "github.com/uptrace/bun"
  domain "nexa/services/user/internal/domain/entity"
  "nexa/shared/types"
  "nexa/shared/util"
  "nexa/shared/util/repo"
  "nexa/shared/variadic"
  "nexa/shared/wrapper"
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
    IsVerified: false,
  }

  variadic.New(opts...).DoAll(repo.MapOptionFunc(domain, &usr))

  return usr
}

type User struct {
  bun.BaseModel `bun:"table:users"`

  Id         string `bun:",type:uuid,pk"`
  Username   string `bun:",nullzero,notnull,unique"`
  Email      string `bun:",nullzero,notnull,unique"`
  Password   string `bun:",nullzero,notnull"`
  IsVerified bool   `bun:",default:false"`

  BannedUntil time.Time `bun:",nullzero"`
  DeletedAt   time.Time `bun:",nullzero"`
  CreatedAt   time.Time `bun:",nullzero,notnull"`
  UpdatedAt   time.Time `bun:",nullzero"`

  Profile *Profile
}

func (u *User) ToDomain() domain.User {
  return domain.User{
    Id:          wrapper.DropError(types.IdFromString(u.Id)),
    Username:    u.Username,
    Email:       wrapper.DropError(types.EmailFromString(u.Email)),
    Password:    wrapper.DropError(types.PasswordFromString(u.Password)),
    IsVerified:  u.IsVerified,
    BannedUntil: u.BannedUntil,
    IsDeleted:   !u.DeletedAt.IsZero(),
    Profile: util.NilOr[domain.Profile](u.Profile, func(obj *Profile) *domain.Profile {
      temp := u.Profile.ToDomain()
      return &temp
    }),
  }
}
