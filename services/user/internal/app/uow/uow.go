package uow

import (
	"context"
	"github.com/uptrace/bun"
	"nexa/services/user/internal/domain/repository"
	"nexa/services/user/internal/infra/repository/pg"
	"nexa/shared/uow"
)

func newUserUOWStorage(user repository.IUser, profile repository.IProfile) UserStorage {
	return UserStorage{user: user, profile: profile}
}

// UserStorage repository storage, made members as private, so it couldn't be initialized outside this package
type UserStorage struct {
	user    repository.IUser
	profile repository.IProfile
}

func (u UserStorage) User() repository.IUser {
	return u.user
}

func (u UserStorage) Profile() repository.IProfile {
	return u.profile
}

func NewUserUOW(db bun.IDB) uow.IUnitOfWork[UserStorage] {
	return &UserUOW{
		db: db,
	}
}

type UserUOW struct {
	db bun.IDB

	cache *UserStorage // Cache for real connection repository
}

func (u *UserUOW) DoTx(ctx context.Context, f uow.UOWBlock[UserStorage]) error {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = f(ctx, u.repositories(tx))
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}

func (u *UserUOW) Repositories() UserStorage {
	if u.cache == nil {
		return *u.cache
	}
	user := pg.NewUser(u.db)
	profile := pg.NewProfile(u.db)

	storage := newUserUOWStorage(user, profile)
	u.cache = &storage
	return storage
}

func (u *UserUOW) repositories(db bun.IDB) UserStorage {
	user := pg.NewUser(db)
	profile := pg.NewProfile(db)
	return newUserUOWStorage(user, profile)
}
