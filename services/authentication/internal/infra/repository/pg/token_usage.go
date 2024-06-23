package pg

import (
  "context"
  "github.com/uptrace/bun"
  "nexa/services/authentication/internal/domain/entity"
  "nexa/services/authentication/internal/domain/repository"
  "nexa/services/authentication/internal/infra/repository/model"
  "nexa/shared/types"
  "nexa/shared/util"
  "nexa/shared/util/repo"
  "time"
)

func NewTokenUsage(db bun.IDB) repository.ITokenUsage {
  return &tokenUsageRepository{db: db}
}

type tokenUsageRepository struct {
  db bun.IDB
}

func (t *tokenUsageRepository) Create(ctx context.Context, usage *entity.TokenUsage) (types.Id, error) {
  dbModel := model.FromTokenUsageDomain(usage, func(domain *entity.TokenUsage, usage *model.TokenUsage) {
    usage.CreatedAt = time.Now()
  })

  res, err := t.db.NewInsert().
    Model(&dbModel).
    Exec(ctx)

  return usage.Id, repo.CheckResult(res, err)
}

func (t *tokenUsageRepository) Patch(ctx context.Context, usage *entity.TokenUsage) error {
  dbModel := model.FromTokenUsageDomain(usage, func(domain *entity.TokenUsage, usage *model.TokenUsage) {
    usage.UpdatedAt = time.Now()
  })

  res, err := t.db.NewUpdate().
    Model(&dbModel).
    WherePK().
    OmitZero().
    Exec(ctx)

  return repo.CheckResult(res, err)
}

func (t *tokenUsageRepository) Delete(ctx context.Context, id types.Id) error {
  res, err := t.db.NewDelete().
    Model(util.Nil[model.TokenUsage]()).
    Where("id = ?", id.Underlying().String()).
    Exec(ctx)

  return repo.CheckResult(res, err)
}

func (t *tokenUsageRepository) Find(ctx context.Context, id types.Id) (entity.TokenUsage, error) {
  var dbModel model.TokenUsage

  err := t.db.NewSelect().
    Model(&dbModel).
    Where("id = ?", id.Underlying().String()).
    Scan(ctx)

  return dbModel.ToDomain(), err
}

func (t *tokenUsageRepository) FindAll(ctx context.Context, parameter repo.QueryParameter) (repo.PaginatedResult[entity.TokenUsage], error) {
  var dbModels []model.TokenUsage

  count, err := t.db.NewSelect().
    Model(&dbModels).
    Limit(int(parameter.Limit)).
    Offset(int(parameter.Offset)).
    ScanAndCount(ctx)

  result := repo.CheckPaginationResult(dbModels, count, err)
  tokenUsages := util.CastSlice(result.Data, func(from *model.TokenUsage) entity.TokenUsage {
    return from.ToDomain()
  })

  return repo.NewPaginatedResult(tokenUsages, uint64(count)), result.Err
}
