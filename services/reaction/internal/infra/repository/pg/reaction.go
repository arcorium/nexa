package pg

import (
  "context"
  "github.com/arcorium/nexa/shared/optional"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  "github.com/arcorium/nexa/shared/util/repo"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "github.com/uptrace/bun"
  "go.opentelemetry.io/otel/trace"
  "nexa/services/reaction/internal/domain/entity"
  "nexa/services/reaction/internal/domain/repository"
  "nexa/services/reaction/internal/infra/repository/model"
  "nexa/services/reaction/util"
  "nexa/services/reaction/util/errors"
)

func NewReaction(db bun.IDB) repository.IReaction {
  return &reactionRepository{
    db:     db,
    tracer: util.GetTracer(),
  }
}

type reactionRepository struct {
  db     bun.IDB
  tracer trace.Tracer
}

type idInput struct {
  Id string `bun:",type:uuid"`
}

func (r *reactionRepository) delete(ctx context.Context, reaction *model.Reaction) error {
  ctx, span := r.tracer.Start(ctx, "ReactionRepository.delete")
  defer span.End()

  res, err := r.db.NewDelete().
    Model(reaction).
    WherePK().
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}

func (r *reactionRepository) insert(ctx context.Context, reaction *model.Reaction) error {
  ctx, span := r.tracer.Start(ctx, "ReactionRepository.insert")
  defer span.End()

  res, err := r.db.NewInsert().
    Model(reaction).
    Returning("NULL").
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}

func (r *reactionRepository) Delsert(ctx context.Context, reaction *entity.Reaction) error {
  ctx, span := r.tracer.Start(ctx, "ReactionRepository.Delsert")
  defer span.End()

  dbModel := model.FromReactionDomain(reaction)
  exists, err := r.db.NewSelect().
    Model(&dbModel).
    WherePK().
    Exists(ctx)

  if err != nil {
    spanUtil.RecordError(err, span)
  }

  if exists {
    return r.delete(ctx, &dbModel)
  }
  return r.insert(ctx, &dbModel)
}

func (r *reactionRepository) FindByItemId(ctx context.Context, itemType entity.ItemType, itemId types.Id, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Reaction], error) {
  ctx, span := r.tracer.Start(ctx, "ReactionRepository.FindByItemId")
  defer span.End()

  var dbModels []model.Reaction
  count, err := r.db.NewSelect().
    Model(&dbModels).
    Where("item_type = ? AND item_id = ?", itemType.Underlying(), itemId.String()).
    Limit(int(parameter.Limit)).
    Offset(int(parameter.Offset)).
    OrderExpr("created_at DESC").
    ScanAndCount(ctx)

  result := repo.CheckPaginationResult(dbModels, count, err)
  if result.IsError() {
    spanUtil.RecordError(result.Err, span)
    return repo.NewPaginatedResult[entity.Reaction](nil, uint64(count)), result.Err
  }

  reactions, ierr := sharedUtil.CastSliceErrsP(result.Data, repo.ToDomainErr[*model.Reaction, entity.Reaction])
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return repo.NewPaginatedResult[entity.Reaction](nil, uint64(count)), ierr
  }

  return repo.NewPaginatedResult[entity.Reaction](reactions, uint64(count)), nil
}

func (r *reactionRepository) FindByUserId(ctx context.Context, userId types.Id, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Reaction], error) {
  ctx, span := r.tracer.Start(ctx, "ReactionRepository.FindByUserId")
  defer span.End()

  var dbModels []model.Reaction
  count, err := r.db.NewSelect().
    Model(&dbModels).
    Where("user_id = ?", userId.String()).
    Limit(int(parameter.Limit)).
    Offset(int(parameter.Offset)).
    OrderExpr("created_at DESC").
    ScanAndCount(ctx)

  result := repo.CheckPaginationResult(dbModels, count, err)
  if result.IsError() {
    spanUtil.RecordError(result.Err, span)
    return repo.NewPaginatedResult[entity.Reaction](nil, uint64(count)), err
  }

  reactions, ierr := sharedUtil.CastSliceErrsP(result.Data, repo.ToDomainErr[*model.Reaction, entity.Reaction])
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return repo.NewPaginatedResult[entity.Reaction](nil, uint64(count)), ierr
  }

  return repo.NewPaginatedResult[entity.Reaction](reactions, uint64(count)), nil
}

func (r *reactionRepository) GetCounts(ctx context.Context, itemType entity.ItemType, itemIds ...types.Id) ([]entity.Count, error) {
  ctx, span := r.tracer.Start(ctx, "ReactionRepository.GetCounts")
  defer span.End()

  ids := sharedUtil.CastSlice(itemIds, sharedUtil.ToString[types.Id])

  var dbModels []model.Count
  input := sharedUtil.CastSlice(itemIds, func(id types.Id) idInput {
    return idInput{Id: id.String()}
  })
  err := r.db.NewSelect().
    With("result", r.db.NewValues(&input).WithOrder()).
    Model(types.Nil[model.Reaction]()).
    ColumnExpr("result.id as item_id").
    ColumnExpr("COUNT(CASE WHEN reaction = ? THEN 1 END) as like_count", entity.ReactionLike.Underlying()).
    ColumnExpr("COUNT(CASE WHEN reaction = ? THEN 1 END) as dislike_count", entity.ReactionDislike.Underlying()).
    GroupExpr("result.id, result._order").
    OrderExpr("result._order").
    Join("RIGHT JOIN result ON item_id = result.id AND item_type = ?", itemType.Underlying()).
    Scan(ctx, &dbModels)

  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  // Result length should be the same with the ids (because it is expected in same order)
  if len(dbModels) != len(ids) {
    err = errors.ErrResultWithDifferentLength
    spanUtil.RecordError(err, span)
    return nil, err
  }

  counts, ierr := sharedUtil.CastSliceErrsP(dbModels, func(from *model.Count) (entity.Count, error) {
    return from.ToDomain(itemType)
  })
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return nil, ierr
  }

  return counts, nil
}

func (r *reactionRepository) DeleteByItemId(ctx context.Context, itemType entity.ItemType, itemIds ...types.Id) error {
  ctx, span := r.tracer.Start(ctx, "ReactionRepository.DeleteByItemId")
  defer span.End()

  ids := sharedUtil.CastSlice(itemIds, sharedUtil.ToString[types.Id])
  res, err := r.db.NewSelect().
    Model(types.Nil[model.Reaction]()).
    Where("item_type = ? AND item_id IN (?)", itemType.Underlying(), bun.In(ids)).
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}

func (r *reactionRepository) DeleteByUserId(ctx context.Context, userId types.Id, items optional.Object[entity.Item]) error {
  ctx, span := r.tracer.Start(ctx, "ReactionRepository.DeleteByUserId")
  defer span.End()

  query := r.db.NewSelect().
    Model(types.Nil[model.Reaction]())

  if items.HasValue() {
    ids := sharedUtil.CastSlice(items.Value().Second, sharedUtil.ToString[types.Id])
    query = query.Where("user_id = ? AND item_type = ? AND item_id IN (?)", items.Value().First.Underlying(), bun.In(ids))
  } else {
    query = query.Where("user_id = ?", userId.String())
  }
  res, err := query.Exec(ctx)
  return repo.CheckResultWithSpan(res, err, span)
}
