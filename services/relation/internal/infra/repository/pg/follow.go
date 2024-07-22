package pg

import (
  "context"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  "github.com/arcorium/nexa/shared/util/repo"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "github.com/uptrace/bun"
  "go.opentelemetry.io/otel/trace"
  "nexa/services/relation/internal/domain/entity"
  "nexa/services/relation/internal/domain/repository"
  "nexa/services/relation/internal/infra/repository/model"
  "nexa/services/relation/util"
  "nexa/services/relation/util/errors"
)

func NewFollow(db bun.IDB) repository.IFollow {
  return &followRepository{
    db:     db,
    tracer: util.GetTracer(),
  }
}

type followRepository struct {
  db     bun.IDB
  tracer trace.Tracer
}

func (f *followRepository) Create(ctx context.Context, follow *entity.Follow) error {
  ctx, span := f.tracer.Start(ctx, "FollowRepository.Create")
  defer span.End()

  dbModel := model.FromFollowDomain(follow)
  res, err := f.db.NewInsert().
    Model(&dbModel).
    Returning("NULL").
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}

func (f *followRepository) Delete(ctx context.Context, follow *entity.Follow) error {
  ctx, span := f.tracer.Start(ctx, "FollowRepository.Delete")
  defer span.End()

  dbModel := model.FromFollowDomain(follow)
  res, err := f.db.NewDelete().
    Model(&dbModel).
    WherePK().
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}

func (f *followRepository) Delsert(ctx context.Context, follow *entity.Follow) error {
  ctx, span := f.tracer.Start(ctx, "FollowRepository.Delsert")
  defer span.End()

  dbModel := model.FromFollowDomain(follow)
  exists, err := f.db.NewSelect().
    Model(&dbModel).
    WherePK().
    Exists(ctx)

  if err != nil {
    spanUtil.RecordError(err, span)
    return err
  }

  if exists {
    return f.Delete(ctx, follow)
  }
  return f.Create(ctx, follow)
}

func (f *followRepository) GetFollowers(ctx context.Context, userId types.Id, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Follow], error) {
  ctx, span := f.tracer.Start(ctx, "FollowRepository.GetFollowers")
  defer span.End()

  var dbModels []model.Follow
  count, err := f.db.NewSelect().
    Model(&dbModels).
    Where("followee_id = ?", userId.String()).
    Limit(int(parameter.Limit)).
    Offset(int(parameter.Offset)).
    OrderExpr("created_at DESC").
    ScanAndCount(ctx)

  result := repo.CheckPaginationResult(dbModels, count, err)
  if result.IsError() {
    spanUtil.RecordError(result.Err, span)
    return repo.NewPaginatedResult[entity.Follow](nil, uint64(count)), result.Err
  }

  resp, ierr := sharedUtil.CastSliceErrsP(dbModels, repo.ToDomainErr[*model.Follow, entity.Follow])
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return repo.NewPaginatedResult[entity.Follow](nil, uint64(count)), ierr
  }

  return repo.NewPaginatedResult(resp, uint64(count)), nil
}

func (f *followRepository) GetFollowings(ctx context.Context, userId types.Id, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Follow], error) {
  ctx, span := f.tracer.Start(ctx, "FollowRepository.GetFollowings")
  defer span.End()

  var dbModels []model.Follow
  count, err := f.db.NewSelect().
    Model(&dbModels).
    Where("follower_id = ?", userId.String()).
    Limit(int(parameter.Limit)).
    Offset(int(parameter.Offset)).
    OrderExpr("created_at DESC").
    ScanAndCount(ctx)

  result := repo.CheckPaginationResult(dbModels, count, err)
  if result.IsError() {
    spanUtil.RecordError(result.Err, span)
    return repo.NewPaginatedResult[entity.Follow](nil, uint64(count)), result.Err
  }

  resp, ierr := sharedUtil.CastSliceErrsP(dbModels, repo.ToDomainErr[*model.Follow, entity.Follow])
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return repo.NewPaginatedResult[entity.Follow](nil, uint64(count)), ierr
  }

  return repo.NewPaginatedResult(resp, uint64(count)), nil
}

func (f *followRepository) IsFollowing(ctx context.Context, userId types.Id, followeeIds ...types.Id) ([]bool, error) {
  ctx, span := f.tracer.Start(ctx, "BlockRepository.IsFollowing")
  defer span.End()

  type Result struct {
    Id          string `bun:"ids,scanonly"`
    IsFollowing bool   `bun:",scanonly"`
  }

  var dbModels []Result
  input := sharedUtil.CastSlice(followeeIds, func(id types.Id) idInput {
    return idInput{Id: id.String()}
  })
  err := f.db.NewSelect().
    With("result", f.db.NewValues(&input).WithOrder()).
    Model(types.Nil[model.Follow]()).
    ColumnExpr("result.id as ids").
    ColumnExpr("COUNT(followee_id) > 0 as is_following").
    GroupExpr("result.id, result._order").
    OrderExpr("result._order").
    //Where("follower_id = ?", userId.String()).
    Join("RIGHT JOIN result ON follow.followee_id = result.id AND follow.follower_id = ?", userId.String()).
    Scan(ctx, &dbModels)

  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  if len(dbModels) != len(followeeIds) {
    // TODO: It should be internal error
    err = errors.ErrResultWithDifferentLength
    spanUtil.RecordError(err, span)
    return nil, err
  }

  counts := sharedUtil.CastSliceP(dbModels, func(from *Result) bool {
    return from.IsFollowing
  })
  return counts, nil
}

func (f *followRepository) GetCounts(ctx context.Context, userIds ...types.Id) ([]entity.FollowCount, error) {
  ctx, span := f.tracer.Start(ctx, "BlockRepository.GetCounts")
  defer span.End()

  var dbModels []model.FollowCount
  input := sharedUtil.CastSlice(userIds, func(id types.Id) idInput {
    return idInput{Id: id.String()}
  })
  err := f.db.NewSelect().
    With("result", f.db.NewValues(&input).WithOrder()).
    Model(types.Nil[model.Follow]()).
    ColumnExpr("result.id as user_id").
    ColumnExpr("SUM(CASE WHEN followee_id = result.id THEN 1 ELSE 0 END) as total_followers").
    ColumnExpr("SUM(CASE WHEN follower_id = result.id THEN 1 ELSE 0 END) as total_followings").
    GroupExpr("result.id, result._order").
    OrderExpr("result._order").
    Join("RIGHT JOIN result ON result.id IN (follow.followee_id, follow.follower_id)").
    Scan(ctx, &dbModels)

  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  if len(dbModels) != len(userIds) {
    // TODO: It should be internal error
    err = errors.ErrResultWithDifferentLength
    spanUtil.RecordError(err, span)
    return nil, err
  }

  counts, ierr := sharedUtil.CastSliceErrsP(dbModels, repo.ToDomainErr[*model.FollowCount, entity.FollowCount])
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return nil, ierr
  }
  return counts, nil
}

// DeleteByUserId delete both the user follower and followed one. Also for this user followers
func (f *followRepository) DeleteByUserId(ctx context.Context, deleteFollower bool, userId types.Id) error {
  ctx, span := f.tracer.Start(ctx, "FollowRepository.DeleteByUserId")
  defer span.End()

  query := f.db.NewSelect().
    Model(types.Nil[model.Block]())
  if deleteFollower {
    // Delete user that follows this user
    query = query.
      Where("follower_id = ? OR followee_id = ?", userId.String(), userId.String())
  } else {
    // Only delete user following list
    query = query.
      Where("follower_id = ? ", userId.String())
  }

  res, err := query.Exec(ctx)
  return repo.CheckResultWithSpan(res, err, span)
}
