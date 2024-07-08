package pg

import (
  "context"
  "github.com/uptrace/bun"
  "go.opentelemetry.io/otel/trace"
  "nexa/services/mailer/internal/domain/entity"
  "nexa/services/mailer/internal/domain/repository"
  "nexa/services/mailer/internal/infra/repository/model"
  "nexa/services/mailer/util"
  "nexa/shared/types"
  sharedUtil "nexa/shared/util"
  "nexa/shared/util/repo"
  spanUtil "nexa/shared/util/span"
  "time"
)

func NewTag(db bun.IDB) repository.ITag {
  return &tagRepository{
    db:     db,
    tracer: util.GetTracer(),
  }
}

type tagRepository struct {
  db     bun.IDB
  tracer trace.Tracer
}

func (t *tagRepository) Get(ctx context.Context, query repo.QueryParameter) (repo.PaginatedResult[entity.Tag], error) {
  ctx, span := t.tracer.Start(ctx, "TagRepository.Get")
  defer span.End()

  var models []model.Tag
  count, err := t.db.NewSelect().
    Model(&models).
    Offset(int(query.Offset)).
    Limit(int(query.Limit)).
    OrderExpr("created_at DESC").
    ScanAndCount(ctx)

  result := repo.CheckPaginationResult(models, count, err)
  if result.IsError() {
    spanUtil.RecordError(result.Err, span)
    return repo.NewPaginatedResult[entity.Tag](nil, uint64(count)), result.Err
  }
  tags, ierr := sharedUtil.CastSliceErrsP(result.Data, repo.ToDomainErr[*model.Tag, entity.Tag])
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return repo.NewPaginatedResult[entity.Tag](nil, uint64(count)), ierr
  }
  return repo.NewPaginatedResult(tags, uint64(count)), nil
}

func (t *tagRepository) FindByIds(ctx context.Context, ids ...types.Id) ([]entity.Tag, error) {
  ctx, span := t.tracer.Start(ctx, "TagRepository.FindByIds")
  defer span.End()

  tagIds := sharedUtil.CastSlice(ids, sharedUtil.ToString[types.Id])

  var models []model.Tag
  err := t.db.NewSelect().
    Model(&models).
    Where("id IN (?)", bun.In(tagIds)).
    Distinct().
    OrderExpr("created_at DESC").
    Scan(ctx)

  result := repo.CheckSliceResult(models, err)
  if result.IsError() {
    spanUtil.RecordError(result.Err, span)
    return nil, result.Err
  }
  tags, ierr := sharedUtil.CastSliceErrsP(result.Data, repo.ToDomainErr[*model.Tag, entity.Tag])
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return nil, ierr
  }
  return tags, result.Err
}

func (t *tagRepository) FindByName(ctx context.Context, name string) (*entity.Tag, error) {
  ctx, span := t.tracer.Start(ctx, "TagRepository.FindByName")
  defer span.End()

  var models model.Tag
  err := t.db.NewSelect().
    Model(&models).
    Where("name = ?", name).
    OrderExpr("created_at DESC").
    Scan(ctx)

  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  result, err := models.ToDomain()
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  return &result, nil
}

func (t *tagRepository) Create(ctx context.Context, tag *entity.Tag) error {
  ctx, span := t.tracer.Start(ctx, "TagRepository.Create")
  defer span.End()

  models := model.FromTagDomain(tag, func(domain *entity.Tag, tag *model.Tag) {
    tag.CreatedAt = time.Now()
  })

  res, err := t.db.NewInsert().
    Model(&models).
    Returning("NULL").
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}

func (t *tagRepository) Patch(ctx context.Context, tag *entity.PatchedTag) error {
  ctx, span := t.tracer.Start(ctx, "TagRepository.Patch")
  defer span.End()

  models := model.FromPatchedTagDomain(tag, func(domain *entity.PatchedTag, tag *model.Tag) {
    tag.UpdatedAt = time.Now()
  })

  res, err := t.db.NewUpdate().
    Model(&models).
    WherePK().
    OmitZero().
    ExcludeColumn("id", "created_at").
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}

func (t *tagRepository) Remove(ctx context.Context, id types.Id) error {
  ctx, span := t.tracer.Start(ctx, "TagRepository.Remove")
  defer span.End()

  res, err := t.db.NewDelete().
    Model(types.Nil[model.Tag]()).
    Where("id = ?", id.String()).
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}
