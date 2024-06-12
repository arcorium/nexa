package pg

import (
  "context"
  "github.com/uptrace/bun"
  "nexa/services/user/internal/domain/entity"
  "nexa/services/user/internal/domain/repository"
  "nexa/services/user/internal/infra/model"
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
  })

  res, err := u.db.NewUpdate().
    Model(&dbModel).
    WherePK().
    ExcludeColumn("id", "created_at", "deleted_at").
    Exec(ctx)

  return repo.CheckResult(res, err)
}

func (u userRepository) Patch(ctx context.Context, user *entity.User) error {
  dbModel := model.FromUserDomain(user, func(domain *entity.User, user *model.User) {
    user.UpdatedAt = time.Now()
  })

  res, err := u.db.NewUpdate().
    Model(&dbModel).
    OmitZero().
    WherePK().
    ExcludeColumn("id").
    Exec(ctx)

  return repo.CheckResult(res, err)
}

func (u userRepository) FindByIds(ctx context.Context, userIds ...types.Id) ([]entity.User, error) {
  uuids := util.CastSlice2(userIds, func(from types.Id) string {
    return from.Underlying().String()
  })

  var dbModel []model.User
  err := u.db.NewSelect().
    Model(&dbModel).
    Where("id IN (?)", bun.In(uuids)).
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

func (u userRepository) FindAllUsers(ctx context.Context, query repo.QueryParameter) (repo.PaginatedResult[entity.User], error) {
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

func (u userRepository) Delete(ctx context.Context, ids ...types.Id) error {
  users := util.CastSlice(ids, func(from *types.Id) model.User {
    return model.User{
      Id:        from.Underlying().String(),
      DeletedAt: time.Now(),
      UpdatedAt: time.Now(),
    }
  })
  // Use soft delete
  res, err := u.db.NewUpdate().
    Model(&users).
    WherePK().
    Column("updated_at", "deleted_at").
    Bulk().
    Exec(ctx)

  return repo.CheckResult(res, err)
}
