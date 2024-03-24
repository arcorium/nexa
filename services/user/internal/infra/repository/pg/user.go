package pg

import (
	"context"
	"github.com/uptrace/bun"
	"nexa/services/user/internal/domain/repository"
	"nexa/services/user/internal/infra/model"
	"nexa/services/user/shared/domain/entity"
	"nexa/shared/types"
	"nexa/shared/util"
	"nexa/shared/util/repo"
	"time"
)

func NewUser(db bun.IDB) repository.IUser {
	return &userRepository{db: db}
}

type userRepository struct {
	db bun.IDB
}

func (u userRepository) Create(ctx context.Context, user *entity.User) error {
	dbModel := model.FromUserDomain(user, func(domain *entity.User, user *model.User) {
		user.CreatedAt = time.Now()
	})

	res, err := u.db.NewInsert().
		Model(&dbModel).
		Returning("NULL").
		Exec(ctx)

	return repo.CheckResult(res, err)
}

func (u userRepository) Update(ctx context.Context, user *entity.User) error {
	dbModel := model.FromUserDomain(user, func(domain *entity.User, user *model.User) {
		user.UpdatedAt = time.Now()
		if domain.ShouldDeleted() {
			user.DeletedAt = time.Now()
		}
	})

	res, err := u.db.NewUpdate().
		Model(&dbModel).
		ExcludeColumn("id").
		Exec(ctx)

	return repo.CheckResult(res, err)
}

func (u userRepository) Patch(ctx context.Context, user *entity.User) error {
	dbModel := model.FromUserDomain(user, func(domain *entity.User, user *model.User) {
		user.UpdatedAt = time.Now()
		if domain.ShouldDeleted() {
			user.DeletedAt = time.Now()
		}
	})

	res, err := u.db.NewUpdate().
		Model(&dbModel).
		OmitZero().
		ExcludeColumn("id").
		Exec(ctx)

	return repo.CheckResult(res, err)
}

func (u userRepository) FindByIds(ctx context.Context, userIds ...types.Id) ([]entity.User, error) {
	var dbModel []model.User
	err := u.db.NewSelect().
		Model(&dbModel).
		Where("user_id IN (?)", bun.In(userIds)).
		Scan(ctx)

	result := repo.CheckSliceResult(dbModel, err)
	users := util.CastSlice(dbModel, func(from *model.User) entity.User {
		return from.ToDomain()
	})
	return users, result.Err
}

func (u userRepository) FindByEmails(ctx context.Context, emails ...types.Email) ([]entity.User, error) {
	var dbModel []model.User
	err := u.db.NewSelect().
		Model(&dbModel).
		Where("email IN (?)", bun.In(emails)).
		Scan(ctx)

	result := repo.CheckSliceResult(dbModel, err)
	users := util.CastSlice(dbModel, func(from *model.User) entity.User {
		return from.ToDomain()
	})
	return users, result.Err
}

func (u userRepository) FindAllUsers(ctx context.Context, query repo.QueryParameter) (repo.PaginatedResult[entity.
	User], error) {
	var dbModel []model.User
	count, err := u.db.NewSelect().
		Model(&dbModel).
		Offset(int(query.Offset)).
		Limit(int(query.Limit)).
		ScanAndCount(ctx)

	result := repo.CheckPaginationResult(dbModel, count, err)
	users := util.CastSlice(dbModel, func(from *model.User) entity.User {
		return from.ToDomain()
	})
	return repo.NewPaginatedResult(users, uint64(count)), result.Err
}

func (u userRepository) Delete(ctx context.Context, id types.Id) error {
	user := model.User{
		Id:        id.Underlying().String(),
		DeletedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	// Use soft delete
	res, err := u.db.NewUpdate().
		Model(&user).
		WherePK().
		Exec(ctx)

	return repo.CheckResult(res, err)
}
