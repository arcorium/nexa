package pg

import (
  "context"
  "github.com/uptrace/bun"
  "go.opentelemetry.io/otel/trace"
  domain "nexa/services/file_storage/internal/domain/entity"
  "nexa/services/file_storage/internal/domain/repository"
  "nexa/services/file_storage/internal/infra/repository/model"
  "nexa/services/file_storage/util"
  "nexa/shared/types"
  sharedUtil "nexa/shared/util"
  "nexa/shared/util/repo"
  spanUtil "nexa/shared/util/span"
  "time"
)

func NewFileMetadataRepository(db bun.IDB) repository.IFileMetadata {
  return &metadataRepository{
    db:     db,
    tracer: util.GetTracer(),
  }
}

type metadataRepository struct {
  db bun.IDB

  tracer trace.Tracer
}

func (f metadataRepository) FindByIds(ctx context.Context, ids ...types.Id) ([]domain.FileMetadata, error) {
  ctx, span := f.tracer.Start(ctx, "FileRepository.FindByIds")
  defer span.End()

  metadataIds := sharedUtil.CastSlice(ids, sharedUtil.ToString[types.Id])

  var models []model.FileMetadata

  err := f.db.NewSelect().
    Model(&models).
    Where("id IN (?)", bun.In(metadataIds)).
    Scan(ctx)

  res := repo.CheckSliceResultWithSpan(models, err, span)
  result, ierr := sharedUtil.CastSliceErrsP(res.Data, repo.ToDomainErr[*model.FileMetadata, domain.FileMetadata])
  if !ierr.IsNil() {
    spanUtil.RecordError(err, span)
    return nil, ierr
  }

  return result, res.Err
}

func (f metadataRepository) FindByNames(ctx context.Context, names ...string) ([]domain.FileMetadata, error) {
  ctx, span := f.tracer.Start(ctx, "FileRepository.FindByNames")
  defer span.End()

  var models []model.FileMetadata

  err := f.db.NewSelect().
    Model(&models).
    Where("filename IN (?)", bun.In(names)).
    Scan(ctx)

  res := repo.CheckSliceResultWithSpan(models, err, span)
  result, ierr := sharedUtil.CastSliceErrsP(res.Data, repo.ToDomainErr[*model.FileMetadata, domain.FileMetadata])
  if !ierr.IsNil() {
    spanUtil.RecordError(err, span)
    return nil, ierr
  }

  return result, res.Err
}

func (f metadataRepository) Create(ctx context.Context, metadata *domain.FileMetadata) error {
  ctx, span := f.tracer.Start(ctx, "FileRepository.Create")
  defer span.End()

  models := model.FromFileDomain(metadata, func(domain *domain.FileMetadata, metadata *model.FileMetadata) {
    metadata.CreatedAt = time.Now()
    metadata.UpdatedAt = time.Now()
  })

  res, err := f.db.NewInsert().
    Model(&models).
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}

func (f metadataRepository) Update(ctx context.Context, metadata *domain.FileMetadata) error {
  ctx, span := f.tracer.Start(ctx, "FileRepository.Update")
  defer span.End()

  models := model.FromFileDomain(metadata, func(domain *domain.FileMetadata, metadata *model.FileMetadata) {
    metadata.UpdatedAt = time.Now()
  })

  res, err := f.db.NewUpdate().
    Model(&models).
    OmitZero().
    WherePK().
    ExcludeColumn("id", "created_at").
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}

func (f metadataRepository) DeleteById(ctx context.Context, id types.Id) error {
  ctx, span := f.tracer.Start(ctx, "FileRepository.DeleteById")
  defer span.End()

  res, err := f.db.NewDelete().
    Model(types.Nil[model.FileMetadata]()).
    Where("id = ?", id.String()).
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}

func (f metadataRepository) DeleteByName(ctx context.Context, name string) error {
  ctx, span := f.tracer.Start(ctx, "FileRepository.DeleteByName")
  defer span.End()

  res, err := f.db.NewDelete().
    Model(types.Nil[model.FileMetadata]()).
    Where("id = ?", name).
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}
