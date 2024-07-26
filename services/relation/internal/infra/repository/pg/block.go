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
  "nexa/services/relation/util/errs"
)

func NewBlock(db bun.IDB) repository.IBlock {
  return &blockRepository{
    db:     db,
    tracer: util.GetTracer(),
  }
}

type blockRepository struct {
  db     bun.IDB
  tracer trace.Tracer
}

type idInput struct {
  Id string `bun:",type:uuid"`
}

func (f *blockRepository) Create(ctx context.Context, block *entity.Block) error {
  ctx, span := f.tracer.Start(ctx, "BlockRepository.Create")
  defer span.End()

  dbModel := model.FromBlockDomain(block)
  res, err := f.db.NewInsert().
    Model(&dbModel).
    Returning("NULL").
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}

func (f *blockRepository) Delete(ctx context.Context, block *entity.Block) error {
  ctx, span := f.tracer.Start(ctx, "BlockRepository.Delete")
  defer span.End()

  dbModel := model.FromBlockDomain(block)
  res, err := f.db.NewDelete().
    Model(&dbModel).
    WherePK().
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}

func (f *blockRepository) Delsert(ctx context.Context, block *entity.Block) error {
  ctx, span := f.tracer.Start(ctx, "BlockRepository.Delsert")
  defer span.End()

  dbModel := model.FromBlockDomain(block)
  exists, err := f.db.NewSelect().
    Model(&dbModel).
    WherePK().
    Exists(ctx)

  if err != nil {
    spanUtil.RecordError(err, span)
    return err
  }

  if exists {
    return f.Delete(ctx, block)
  }
  return f.Create(ctx, block)
}

func (f *blockRepository) GetBlocked(ctx context.Context, userId types.Id, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Block], error) {
  ctx, span := f.tracer.Start(ctx, "BlockRepository.GetBlocked")
  defer span.End()

  var dbModels []model.Block
  count, err := f.db.NewSelect().
    Model(&dbModels).
    Where("blocker_id = ?", userId.String()).
    Limit(int(parameter.Limit)).
    Offset(int(parameter.Offset)).
    OrderExpr("created_at DESC").
    ScanAndCount(ctx)

  result := repo.CheckSliceResult(dbModels, err)
  if result.IsError() {
    spanUtil.RecordError(result.Err, span)
    return repo.NewPaginatedResult[entity.Block](nil, uint64(count)), result.Err
  }

  resp, ierr := sharedUtil.CastSliceErrsP(dbModels, repo.ToDomainErr[*model.Block, entity.Block])
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return repo.NewPaginatedResult[entity.Block](nil, uint64(count)), ierr
  }

  return repo.NewPaginatedResult(resp, uint64(count)), nil
}

func (f *blockRepository) GetCounts(ctx context.Context, userIds ...types.Id) ([]entity.BlockCount, error) {
  ctx, span := f.tracer.Start(ctx, "BlockRepository.GetCounts")
  defer span.End()

  var dbModels []model.BlockedCount
  input := sharedUtil.CastSlice(userIds, func(id types.Id) idInput {
    return idInput{Id: id.String()}
  })

  err := f.db.NewSelect().
    With("result", f.db.NewValues(&input).WithOrder()).
    Model(types.Nil[model.Block]()).
    ColumnExpr("result.id as user_id").
    ColumnExpr("COUNT(blocker_id) as total_blocked").
    GroupExpr("result.id, result._order").
    OrderExpr("result._order").
    Join("RIGHT JOIN result ON block.blocker_id = result.id").
    Scan(ctx, &dbModels)

  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  if len(dbModels) != len(userIds) {
    // TODO: It should be internal error
    err = errs.ErrResultWithDifferentLength
    spanUtil.RecordError(err, span)
    return nil, err
  }

  counts, ierr := sharedUtil.CastSliceErrsP(dbModels, repo.ToDomainErr[*model.BlockedCount, entity.BlockCount])
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return nil, ierr
  }
  return counts, nil
}

func (f *blockRepository) IsBlocked(ctx context.Context, blockerId, targetId types.Id) (bool, error) {
  ctx, span := f.tracer.Start(ctx, "BlockRepository.IsBlocked")
  defer span.End()

  count, err := f.db.NewSelect().
    Model(types.Nil[model.Block]()).
    Where("blocker_id = ? AND blocked_id = ?", blockerId.String(), targetId.String()).
    Count(ctx)

  if err != nil {
    spanUtil.RecordError(err, span)
    return false, err
  }

  return count >= 1, nil
}

// DeleteByUserId delete blocked user and user that blocks this user
func (f *blockRepository) DeleteByUserId(ctx context.Context, deleteBlocker bool, userId types.Id) error {
  ctx, span := f.tracer.Start(ctx, "BlockRepository.DeleteByUserId")
  defer span.End()

  query := f.db.NewDelete().
    Model(types.Nil[model.Block]())
  if deleteBlocker {
    // Delete other user that blocked this user
    query = query.
      Where("blocker_id = ? OR blocked_id = ?", userId.String(), userId.String())
  } else {
    // Only delete user blocked list
    query = query.
      Where("blocker_id = ? ", userId.String())
  }

  res, err := query.Exec(ctx)
  return repo.CheckResultWithSpan(res, err, span)
}
