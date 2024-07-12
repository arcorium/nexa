package pg

import (
  "context"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  "github.com/arcorium/nexa/shared/util/repo"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "github.com/uptrace/bun"
  "go.opentelemetry.io/otel/trace"
  "nexa/services/authorization/internal/domain/entity"
  "nexa/services/authorization/internal/domain/repository"
  "nexa/services/authorization/internal/infra/repository/model"
  "nexa/services/authorization/util"
  "time"
)

func NewPermission(db bun.IDB) repository.IPermission {
  return &permissionRepository{
    db:     db,
    tracer: util.GetTracer(),
  }
}

type permissionRepository struct {
  db     bun.IDB
  tracer trace.Tracer
}

func (p *permissionRepository) FindByIds(ctx context.Context, ids ...types.Id) ([]entity.Permission, error) {
  ctx, span := p.tracer.Start(ctx, "PermissionRepository.FindById")
  defer span.End()

  permIds := sharedUtil.CastSlice(ids, sharedUtil.ToString[types.Id])

  var dbModels []model.Permission
  err := p.db.NewSelect().
    Model(&dbModels).
    Where("id IN (?)", bun.In(permIds)).
    OrderExpr("created_at").
    Distinct().
    Scan(ctx)

  result := repo.CheckSliceResult(dbModels, err)
  if result.IsError() {
    spanUtil.RecordError(result.Err, span)
    return nil, result.Err
  }

  permissions, ierr := sharedUtil.CastSliceErrsP(result.Data, repo.ToDomainErr[*model.Permission, entity.Permission])
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return nil, ierr
  }

  return permissions, nil
}

func (p *permissionRepository) FindByRoleIds(ctx context.Context, roleIds ...types.Id) ([]entity.Permission, error) {
  ctx, span := p.tracer.Start(ctx, "PermissionRepository.FindByRoleId")
  defer span.End()

  ids := sharedUtil.CastSlice(roleIds, sharedUtil.ToString[types.Id])

  var dbModels []model.RolePermission
  err := p.db.NewSelect().
    Model(&dbModels).
    Relation("Permission").
    Where("role_id IN (?)", bun.In(ids)).
    OrderExpr("permission_id, permission.created_at").
    DistinctOn("permission_id").
    Scan(ctx)

  result := repo.CheckSliceResult(dbModels, err)
  if result.IsError() {
    spanUtil.RecordError(result.Err, span)
    return nil, result.Err
  }

  permissions, ierr := sharedUtil.CastSliceErrsP(result.Data, func(from *model.RolePermission) (entity.Permission, error) {
    return from.Permission.ToDomain()
  })
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return nil, ierr
  }

  return permissions, nil
}

func (p *permissionRepository) Get(ctx context.Context, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Permission], error) {
  ctx, span := p.tracer.Start(ctx, "PermissionRepository.Get")
  defer span.End()

  var dbModels []model.Permission

  count, err := p.db.NewSelect().
    Model(&dbModels).
    Limit(int(parameter.Limit)).
    Offset(int(parameter.Offset)).
    OrderExpr("created_at").
    ScanAndCount(ctx)

  result := repo.CheckPaginationResult(dbModels, count, err)
  if result.IsError() {
    spanUtil.RecordError(result.Err, span)
    return repo.NewPaginatedResult[entity.Permission](nil, uint64(count)), result.Err
  }

  permissions, ierr := sharedUtil.CastSliceErrsP(result.Data, repo.ToDomainErr[*model.Permission, entity.Permission])
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return repo.NewPaginatedResult[entity.Permission](nil, uint64(count)), ierr
  }

  return repo.NewPaginatedResult(permissions, uint64(count)), nil
}

func (p *permissionRepository) Create(ctx context.Context, permission *entity.Permission) error {
  ctx, span := p.tracer.Start(ctx, "PermissionRepository.Create")
  defer span.End()

  dbModel := model.FromPermissionDomain(permission, func(domain *entity.Permission, permission *model.Permission) {
    permission.CreatedAt = time.Now()
  })

  res, err := p.db.NewInsert().
    Model(&dbModel).
    Returning("NULL").
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}

func (p *permissionRepository) Creates(ctx context.Context, permissions ...entity.Permission) error {
  ctx, span := p.tracer.Start(ctx, "PermissionRepository.Create")
  defer span.End()

  dbModels := sharedUtil.CastSliceP(permissions, func(perm *entity.Permission) model.Permission {
    return model.FromPermissionDomain(perm, func(domain *entity.Permission, permission *model.Permission) {
      permission.CreatedAt = time.Now()
    })
  })

  res, err := p.db.NewInsert().
    Model(&dbModels).
    Returning("NULL").
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}

func (p *permissionRepository) Delete(ctx context.Context, id types.Id) error {
  ctx, span := p.tracer.Start(ctx, "PermissionRepository.Delete")
  defer span.End()

  res, err := p.db.NewDelete().
    Model(types.Nil[model.Permission]()).
    Where("id = ?", id.String()).
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}
