package pg

import (
  "context"
  "github.com/uptrace/bun"
  "go.opentelemetry.io/otel/trace"
  "nexa/services/authorization/internal/domain/entity"
  "nexa/services/authorization/internal/domain/repository"
  "nexa/services/authorization/internal/infra/model"
  "nexa/services/authorization/util"
  spanUtil "nexa/shared/span"
  "nexa/shared/types"
  sharedUtil "nexa/shared/util"
  "nexa/shared/util/repo"
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

  permIds := sharedUtil.CastSlice(ids, func(permId types.Id) string {
    return permId.String()
  })

  var dbModels []model.Permission
  err := p.db.NewSelect().
    Model(&dbModels).
    Where("id IN (?)", bun.In(permIds)).
    Scan(ctx)

  result := repo.CheckSliceResultWithSpan(dbModels, err, span)

  permissions, ierr := sharedUtil.CastSliceErrsP(result.Data, func(perm *model.Permission) (entity.Permission, error) {
    id, err := types.IdFromString(perm.Id)
    if err != nil {
      return entity.Permission{}, err
    }

    return entity.Permission{
      Id:   id,
      Code: perm.Code,
    }, nil
  })

  if !ierr.IsNil() {
    spanUtil.RecordError(err, span)
    return nil, ierr
  }

  return permissions, result.Err
}

func (p *permissionRepository) FindByRoleIds(ctx context.Context, roleIds ...types.Id) ([]entity.Permission, error) {
  ctx, span := p.tracer.Start(ctx, "PermissionRepository.FindByRoleId")
  defer span.End()

  ids := sharedUtil.CastSlice(roleIds, func(roleId types.Id) string {
    return roleId.String()
  })

  var dbModels []model.RolePermission
  err := p.db.NewSelect().
    Model(&dbModels).
    Relation("Permission").
    Where("role_id IN (?)", bun.In(ids)).
    Scan(ctx)

  result := repo.CheckSliceResultWithSpan(dbModels, err, span)
  permissions, ierr := sharedUtil.CastSliceErrsP(result.Data, func(from *model.RolePermission) (entity.Permission, error) {
    return from.Permission.ToDomain()
  })
  if !ierr.IsNil() {
    spanUtil.RecordError(err, span)
    return nil, ierr
  }

  return permissions, result.Err
}

func (p *permissionRepository) FindAll(ctx context.Context, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Permission], error) {
  ctx, span := p.tracer.Start(ctx, "PermissionRepository.FindAll")
  defer span.End()

  var dbModels []model.Permission

  count, err := p.db.NewSelect().
    Model(&dbModels).
    Limit(int(parameter.Limit)).
    Offset(int(parameter.Offset)).
    ScanAndCount(ctx)

  result := repo.CheckPaginationResult(dbModels, count, err)
  permissions, ierr := sharedUtil.CastSliceErrsP(result.Data, func(from *model.Permission) (entity.Permission, error) {
    return from.ToDomain()
  })

  if !ierr.IsNil() {
    spanUtil.RecordError(err, span)
    return repo.NewPaginatedResult[entity.Permission](nil, 0), ierr
  }

  return repo.NewPaginatedResult(permissions, uint64(count)), result.Err
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

func (p *permissionRepository) Delete(ctx context.Context, id types.Id) error {
  ctx, span := p.tracer.Start(ctx, "PermissionRepository.Delete")
  defer span.End()

  res, err := p.db.NewDelete().
    Model(types.Nil[model.Permission]()).
    Where("id = ?", id.String()).
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}
