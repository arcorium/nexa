package pg

import (
  "context"
  "database/sql"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  "github.com/arcorium/nexa/shared/util/repo"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "github.com/uptrace/bun"
  "go.opentelemetry.io/otel/trace"
  "nexa/services/authentication/internal/domain/entity"
  "nexa/services/authentication/internal/domain/repository"
  "nexa/services/authentication/internal/infra/repository/model"
  "nexa/services/authentication/util"
  "time"
)

func NewUser(db bun.IDB) repository.IUser {
  return &userRepository{
    db:     db,
    tracer: util.GetTracer(),
  }
}

type userRepository struct {
  db bun.IDB

  tracer trace.Tracer
}

func (u userRepository) Create(ctx context.Context, user *entity.User) error {
  ctx, span := u.tracer.Start(ctx, "UserRepository.Create")
  defer span.End()

  dbModel := model.FromUserDomain(user, func(domain *entity.User, user *model.User) {
    user.CreatedAt = time.Now()
  })

  res, err := u.db.NewInsert().
    Model(&dbModel).
    Returning("NULL").
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}

func (u userRepository) Update(ctx context.Context, user *entity.User) error {
  ctx, span := u.tracer.Start(ctx, "UserRepository.Update")
  defer span.End()

  dbModel := model.FromUserDomain(user, func(domain *entity.User, user *model.User) {
    user.UpdatedAt = time.Now()
  })

  res, err := u.db.NewUpdate().
    Model(&dbModel).
    WherePK().
    ExcludeColumn("id", "created_at").
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}

func (u userRepository) Patch(ctx context.Context, user *entity.PatchedUser) error {
  ctx, span := u.tracer.Start(ctx, "UserRepository.Patch")
  defer span.End()

  dbModel := model.FromPatchedUserDomain(user, func(domain *entity.PatchedUser, user *model.User) {
    user.UpdatedAt = time.Now()
  })

  res, err := u.db.NewUpdate().
    Model(&dbModel).
    OmitZero().
    WherePK().
    ExcludeColumn("id", "created_at").
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}

func (u userRepository) FindByIds(ctx context.Context, userIds ...types.Id) ([]entity.User, error) {
  ctx, span := u.tracer.Start(ctx, "UserRepository.FindByIds")
  defer span.End()

  ids := sharedUtil.CastSlice(userIds, sharedUtil.ToString[types.Id])

  var dbModel []model.User
  err := u.db.NewSelect().
    Model(&dbModel).
    Relation("Profile").
    Where("u.id IN (?)", bun.In(ids)).
    Distinct().
    OrderExpr("created_at DESC").
    Scan(ctx)

  result := repo.CheckSliceResultWithSpan(dbModel, err, span)
  if result.IsError() {
    return nil, result.Err
  }

  users, ierr := sharedUtil.CastSliceErrsP(dbModel, repo.ToDomainErr[*model.User, entity.User])
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return nil, ierr
  }

  return users, nil
}

func (u userRepository) FindByEmails(ctx context.Context, emails ...types.Email) ([]entity.User, error) {
  ctx, span := u.tracer.Start(ctx, "UserRepository.FindByEmails")
  defer span.End()

  var dbModel []model.User
  err := u.db.NewSelect().
    Model(&dbModel).
    Relation("Profile").
    Where("email IN (?)", bun.In(emails)).
    Distinct().
    OrderExpr("created_at DESC").
    Scan(ctx)

  result := repo.CheckSliceResultWithSpan(dbModel, err, span)
  if result.IsError() {
    return nil, result.Err
  }

  users, ierr := sharedUtil.CastSliceErrsP(dbModel, repo.ToDomainErr[*model.User, entity.User])
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return nil, ierr
  }

  return users, nil
}

func (u userRepository) Get(ctx context.Context, query repo.QueryParameter) (repo.PaginatedResult[entity.User], error) {
  ctx, span := u.tracer.Start(ctx, "UserRepository.Get")
  defer span.End()

  var dbModel []model.User
  count, err := u.db.NewSelect().
    Model(&dbModel).
    Offset(int(query.Offset)).
    Limit(int(query.Limit)).
    OrderExpr("created_at DESC").
    ScanAndCount(ctx)

  result := repo.CheckPaginationResultWithSpan(dbModel, count, err, span)
  if result.IsError() {
    return repo.NewPaginatedResult[entity.User](nil, uint64(count)), result.Err
  }

  users, ierr := sharedUtil.CastSliceErrsP(dbModel, repo.ToDomainErr[*model.User, entity.User])
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return repo.NewPaginatedResult[entity.User](nil, uint64(count)), ierr
  }

  return repo.NewPaginatedResult(users, uint64(count)), nil
}

func (u userRepository) Delete(ctx context.Context, ids ...types.Id) error {
  ctx, span := u.tracer.Start(ctx, "UserRepository.Delete")
  defer span.End()

  users := sharedUtil.CastSlice(ids, func(id types.Id) model.User {
    return model.User{
      Id: id.String(),
      DeletedAt: sql.NullTime{
        Time:  time.Now(),
        Valid: true,
      },
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

  return repo.CheckResultWithSpan(res, err, span)
}
