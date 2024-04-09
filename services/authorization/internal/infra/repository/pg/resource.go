package pg

import (
  "context"
  "github.com/uptrace/bun"
  "nexa/services/authorization/internal/domain/entity"
  "nexa/services/authorization/internal/domain/repository"
  "nexa/services/authorization/internal/infra/model"
  "nexa/shared/types"
  "nexa/shared/util"
  "nexa/shared/util/repo"
  "time"
)

func NewResource(db bun.IDB) repository.IResource {
  return &resourceRepository{db: db}
}

type resourceRepository struct {
  db bun.IDB
}

func (r *resourceRepository) FindById(ctx context.Context, id types.Id) (entity.Resource, error) {
  var dbModel model.Resource
  err := r.db.NewSelect().
    Model(&dbModel).
    Scan(ctx)

  if err != nil {
    return entity.Resource{}, err
  }
  return dbModel.ToDomain(), nil
}

func (r *resourceRepository) FindAll(ctx context.Context, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Resource], error) {
  var dbModels []model.Resource

  count, err := r.db.NewSelect().
    Model(&dbModels).
    ScanAndCount(ctx)

  result := repo.CheckPaginationResult(dbModels, count, err)
  resources := util.CastSlice(result.Data, func(from *model.Resource) entity.Resource {
    return from.ToDomain()
  })

  return repo.NewPaginatedResult(resources, uint64(count)), result.Err
}

func (r *resourceRepository) Create(ctx context.Context, resource *entity.Resource) error {
  dbModel := model.FromResourceDomain(resource, func(domain *entity.Resource, resource *model.Resource) {
    resource.CreatedAt = time.Now()
  })

  res, err := r.db.NewInsert().
    Model(&dbModel).
    Exec(ctx)

  return repo.CheckResult(res, err)
}

func (r *resourceRepository) Patch(ctx context.Context, resource *entity.Resource) error {
  dbModel := model.FromResourceDomain(resource, func(domain *entity.Resource, resource *model.Resource) {
    resource.UpdatedAt = time.Now()
  })
  res, err := r.db.NewUpdate().
    Model(&dbModel).
    WherePK().
    OmitZero().
    Exec(ctx)

  return repo.CheckResult(res, err)
}

func (r *resourceRepository) DeleteById(ctx context.Context, id types.Id) error {
  res, err := r.db.NewDelete().
    Model(util.Nil[model.Resource]()).
    Where("id = ?", id.Underlying().String()).
    Exec(ctx)

  return repo.CheckResult(res, err)
}
