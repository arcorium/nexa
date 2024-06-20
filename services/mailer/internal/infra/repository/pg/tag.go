package pg

import (
  "context"
  "github.com/uptrace/bun"
  "go.opentelemetry.io/otel/trace"
  domain "nexa/services/mailer/internal/domain/entity"
  "nexa/services/mailer/internal/domain/repository"
  "nexa/services/mailer/internal/infra/repository/model"
  "nexa/services/mailer/util"
  spanUtil "nexa/shared/span"
  "nexa/shared/types"
  sharedUtil "nexa/shared/util"
  "nexa/shared/util/repo"
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

func (t *tagRepository) FindAll(ctx context.Context, query repo.QueryParameter) (repo.PaginatedResult[domain.Tag], error) {
  ctx, span := t.tracer.Start(ctx, "TagRepository.FindAll")
  defer span.End()

  var models []model.Tag
  count, err := t.db.NewSelect().
    Model(&models).
    Offset(int(query.Offset)).
    Limit(int(query.Limit)).
    ScanAndCount(ctx)

  result := repo.CheckPaginationResultWithSpan(models, count, err, span)
  users := sharedUtil.CastSliceP(result.Data, func(from *model.Tag) domain.Tag {
    return from.ToDomain()
  })
  return repo.NewPaginatedResult(users, uint64(count)), result.Err
}

func (t *tagRepository) FindByIds(ctx context.Context, ids ...types.Id) ([]domain.Tag, error) {
  ctx, span := t.tracer.Start(ctx, "TagRepository.FindByIds")
  defer span.End()

  tagIds := sharedUtil.CastSlice(ids, func(from types.Id) string {
    return from.Underlying().String()
  })

  var models []model.Tag
  err := t.db.NewSelect().
    Model(&models).
    Where("id IN (?)", bun.In(tagIds)).
    Scan(ctx)

  res := repo.CheckSliceResultWithSpan(models, err, span)
  result := sharedUtil.CastSliceP(res.Data, func(from *model.Tag) domain.Tag {
    return from.ToDomain()
  })

  return result, res.Err
}

func (t *tagRepository) FindByName(ctx context.Context, name string) (*domain.Tag, error) {
  ctx, span := t.tracer.Start(ctx, "TagRepository.FindByName")
  defer span.End()

  var models model.Tag
  err := t.db.NewSelect().
    Model(&models).
    Where("id = ?", name).
    Scan(ctx)

  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  result := models.ToDomain()

  return &result, nil
}

func (t *tagRepository) Create(ctx context.Context, tag *domain.Tag) error {
  ctx, span := t.tracer.Start(ctx, "TagRepository.Create")
  defer span.End()

  models := model.FromTagDomain(tag, func(domain *domain.Tag, tag *model.Tag) {
    tag.CreatedAt = time.Now()
  })

  res, err := t.db.NewInsert().
    Model(&models).
    Returning("NULL").
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}

func (t *tagRepository) Patch(ctx context.Context, tag *domain.Tag) error {
  ctx, span := t.tracer.Start(ctx, "TagRepository.Patch")
  defer span.End()

  models := model.FromTagDomain(tag, func(domain *domain.Tag, tag *model.Tag) {
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
    Model(sharedUtil.Nil[model.Tag]()).
    Where("id = ?", id.Underlying().String()).
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}
