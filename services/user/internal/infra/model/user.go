package model

import (
	"github.com/uptrace/bun"
	domain "nexa/services/user/shared/domain/entity"
	"nexa/shared/types"
	"nexa/shared/variadic"
	"time"
)

type UserMapOption = DataAccessModelMapOption[*domain.User, *User]

func FromUserDomain(user *domain.User, opts ...UserMapOption) User {
	// Time based field should be handled individually
	usr := User{
		Id:         user.Id.Underlying().String(),
		Username:   user.Username,
		Email:      user.Email.Underlying(),
		Password:   user.Password.Underlying(),
		IsVerified: false,
	}

	variadic.New(opts...).DoAtFirst(func(i *UserMapOption) {
		(*i)(user, &usr)
	})

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

	//Roles []entity.Role `bun:"rel:has-many,join:id=user_id"`
}

func (u *User) ToDomain() domain.User {
	return domain.User{
		Id:          types.IdFromString(u.Id),
		Username:    u.Username,
		Email:       types.EmailFromString(u.Email),
		Password:    types.PasswordFromString(u.Password),
		IsVerified:  u.IsVerified,
		BannedUntil: u.BannedUntil,
		IsDeleted:   !u.DeletedAt.IsZero(),
	}
}
