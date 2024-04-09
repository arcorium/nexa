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

func NewAction(db bun.IDB) repository.IAction {
  return &actionRepository{db: db}
}

type actionRepository struct {
  db bun.IDB
}

func (a *actionRepository) FindById(ctx context.Context, id types.Id) (entity.Action, error) {
  dbModel := model.Action{}
  err := a.db.NewSelect().
    Model(&dbModel).
    Where("id = ?", id.Underlying().String()).
    Scan(ctx)

  if err != nil {
    return entity.Action{}, err
  }
  return dbModel.ToDomain(), nil
}

func (a *actionRepository) FindAll(ctx context.Context, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Action], error) {
  var dbModels []model.Action
  count, err := a.db.NewSelect().
    Model(&dbModels).
    Limit(int(parameter.Limit)).
    Offset(int(parameter.Offset)).
    ScanAndCount(ctx)

  result := repo.CheckPaginationResult(dbModels, count, err)
  actions := util.CastSlice(result.Data, func(from *model.Action) entity.Action {
    return from.ToDomain()
  })
  return repo.NewPaginatedResult(actions, uint64(count)), result.Err
}

func (a *actionRepository) Create(ctx context.Context, action *entity.Action) error {
  dbModel := model.FromActionDomain(action, func(domain *entity.Action, action *model.Action) {
    action.CreatedAt = time.Now()
  })

  res, err := a.db.NewInsert().
    Model(&dbModel).
    Exec(ctx)

  return repo.CheckResult(res, err)
}

func (a *actionRepository) Patch(ctx context.Context, action *entity.Action) error {
  dbModel := model.FromActionDomain(action, func(domain *entity.Action, action *model.Action) {
    action.UpdatedAt = time.Now()
  })

  res, err := a.db.NewUpdate().
    Model(&dbModel).
    WherePK().
    OmitZero().
    Exec(ctx)

  return repo.CheckResult(res, err)
}

func (a *actionRepository) DeleteById(ctx context.Context, id types.Id) error {
  res, err := a.db.NewDelete().
    Where("id = ?", id.Underlying().String()).
    Exec(ctx)

  return repo.CheckResult(res, err)
}
