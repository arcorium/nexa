package pg

import (
  "context"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  "github.com/arcorium/nexa/shared/util/repo"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "github.com/uptrace/bun"
  "go.opentelemetry.io/otel/trace"
  entity "nexa/services/file_storage/internal/domain/entity"
  "nexa/services/file_storage/internal/domain/repository"
  "nexa/services/file_storage/internal/infra/repository/model"
  "nexa/services/file_storage/util"
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

func (f metadataRepository) FindByIds(ctx context.Context, ids ...types.Id) ([]entity.FileMetadata, error) {
  ctx, span := f.tracer.Start(ctx, "FileRepository.FindByIds")
  defer span.End()

  metadataIds := sharedUtil.CastSlice(ids, sharedUtil.ToString[types.Id])
  var models []model.FileMetadata
  err := f.db.NewSelect().
    Model(&models).
    Where("id IN (?)", bun.In(metadataIds)).
    Distinct().
    OrderExpr("created_at DESC").
    Scan(ctx)

  res := repo.CheckSliceResult(models, err)
  if res.IsError() {
    spanUtil.RecordError(res.Err, span)
    return nil, res.Err
  }

  result, ierr := sharedUtil.CastSliceErrsP(res.Data, repo.ToDomainErr[*model.FileMetadata, entity.FileMetadata])
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return nil, ierr
  }

  return result, nil
}

func (f metadataRepository) FindByNames(ctx context.Context, names ...string) ([]entity.FileMetadata, error) {
  ctx, span := f.tracer.Start(ctx, "FileRepository.FindByNames")
  defer span.End()

  var models []model.FileMetadata
  err := f.db.NewSelect().
    Model(&models).
    Where("filename IN (?)", bun.In(names)).
    OrderExpr("created_at DESC").
    Scan(ctx)

  res := repo.CheckSliceResult(models, err)
  if res.IsError() {
    spanUtil.RecordError(res.Err, span)
    return nil, res.Err
  }

  result, ierr := sharedUtil.CastSliceErrsP(res.Data, repo.ToDomainErr[*model.FileMetadata, entity.FileMetadata])
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return nil, ierr
  }

  return result, nil
}

func (f metadataRepository) Create(ctx context.Context, metadata *entity.FileMetadata) error {
  ctx, span := f.tracer.Start(ctx, "FileRepository.Create")
  defer span.End()

  models := model.FromFileDomain(metadata, func(domain *entity.FileMetadata, metadata *model.FileMetadata) {
    metadata.CreatedAt = time.Now()
    metadata.UpdatedAt = time.Now()
    metadata.FullPath = &domain.FullPath
  })

  res, err := f.db.NewInsert().
    Model(&models).
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}

func (f metadataRepository) Patch(ctx context.Context, metadata *entity.PatchedFileMetadata) error {
  ctx, span := f.tracer.Start(ctx, "FileRepository.Patch")
  defer span.End()

  models := model.FromPatchedDomain(metadata, func(domain *entity.PatchedFileMetadata, metadata *model.FileMetadata) {
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
