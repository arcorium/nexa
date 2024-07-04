package pg

import (
  "context"
  "github.com/uptrace/bun"
  "go.opentelemetry.io/otel/trace"
  "nexa/services/authorization/internal/domain/entity"
  "nexa/services/authorization/internal/domain/repository"
  "nexa/services/authorization/internal/infra/repository/model"
  "nexa/services/authorization/util"
  "nexa/shared/types"
  sharedUtil "nexa/shared/util"
  "nexa/shared/util/repo"
  spanUtil "nexa/shared/util/span"
  "time"
)

func NewRole(db bun.IDB) repository.IRole {
  return &roleRepository{
    db:     db,
    tracer: util.GetTracer(),
  }
}

type roleRepository struct {
  db bun.IDB

  tracer trace.Tracer
}

func (r *roleRepository) FindByIds(ctx context.Context, ids ...types.Id) ([]entity.Role, error) {
  ctx, span := r.tracer.Start(ctx, "RoleRepository.FindByIds")
  defer span.End()

  roleIds := sharedUtil.CastSlice(ids, sharedUtil.ToString[types.Id])
  var dbModels []model.Role

  err := r.db.NewSelect().
    Model(&dbModels).
    Where("id IN (?)", bun.In(roleIds)).
    Relation("Permissions").
    Scan(ctx)

  result := repo.CheckSliceResultWithSpan(dbModels, err, span)
  roles, ierr := sharedUtil.CastSliceErrsP(dbModels, repo.ToDomainErr[*model.Role, entity.Role])
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return nil, ierr
  }

  return roles, result.Err
}

func (r *roleRepository) FindByUserId(ctx context.Context, userId types.Id) ([]entity.Role, error) {
  ctx, span := r.tracer.Start(ctx, "RoleRepository.FindByUserId")
  defer span.End()

  var dbModels []model.UserRole

  err := r.db.NewSelect().
    Model(&dbModels).
    Where("user_id = ?", userId.String()).
    Relation("Role").
    Relation("Role.Permissions").
    Scan(ctx)

  result := repo.CheckSliceResultWithSpan(dbModels, err, span)
  roles, ierr := sharedUtil.CastSliceErrsP(dbModels, func(userRole *model.UserRole) (entity.Role, error) {
    return userRole.Role.ToDomain()
  })
  if !ierr.IsNil() {
    spanUtil.RecordError(err, span)
    return nil, ierr
  }

  return roles, result.Err
}

func (r *roleRepository) FindByName(ctx context.Context, name string) (entity.Role, error) {
  ctx, span := r.tracer.Start(ctx, "RoleRepository.FindByName")
  defer span.End()

  var dbModels model.Role

  err := r.db.NewSelect().
    Model(&dbModels).
    Where("name = ?", name).
    Relation("Role").
    Relation("Role.Permissions"). // NOTE: Is it needed?
    Scan(ctx)

  if err != nil {
    spanUtil.RecordError(err, span)
    return entity.Role{}, err
  }

  return dbModels.ToDomain()
}

func (r *roleRepository) FindAll(ctx context.Context, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Role], error) {
  ctx, span := r.tracer.Start(ctx, "RoleRepository.FindAll")
  defer span.End()

  var dbModels []model.Role

  count, err := r.db.NewSelect().
    Model(&dbModels).
    Relation("Permissions").
    Offset(int(parameter.Offset)).
    Limit(int(parameter.Limit)).
    ScanAndCount(ctx)

  result := repo.CheckPaginationResult(dbModels, count, err)
  roles, ierr := sharedUtil.CastSliceErrsP(dbModels, repo.ToDomainErr[*model.Role, entity.Role])
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return repo.NewPaginatedResult[entity.Role](nil, 0), ierr
  }

  return repo.NewPaginatedResult(roles, uint64(count)), result.Err
}

func (r *roleRepository) Create(ctx context.Context, role *entity.Role) error {
  ctx, span := r.tracer.Start(ctx, "RoleRepository.Create")
  defer span.End()

  dbModel := model.FromRoleDomain(role, func(domain *entity.Role, role *model.Role) {
    role.CreatedAt = time.Now()
  })

  res, err := r.db.NewInsert().
    Model(&dbModel).
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}

func (r *roleRepository) Patch(ctx context.Context, role *entity.Role) error {
  ctx, span := r.tracer.Start(ctx, "RoleRepository.Patch")
  defer span.End()

  dbModel := model.FromRoleDomain(role, func(domain *entity.Role, role *model.Role) {
    role.UpdatedAt = time.Now()
  })

  res, err := r.db.NewUpdate().
    Model(&dbModel).
    WherePK().
    OmitZero().
    ExcludeColumn("created_at").
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}

func (r *roleRepository) Delete(ctx context.Context, id types.Id) error {
  ctx, span := r.tracer.Start(ctx, "RoleRepository.Delete")
  defer span.End()

  res, err := r.db.NewDelete().
    Model(types.Nil[model.Role]()).
    Where("id = ?", id.String()).
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}

func (r *roleRepository) AddPermissions(ctx context.Context, roleId types.Id, permissionIds ...types.Id) error {
  ctx, span := r.tracer.Start(ctx, "RoleRepository.AddPermissions")
  defer span.End()

  permissionModels := sharedUtil.CastSliceP(permissionIds, func(permId *types.Id) model.RolePermission {
    return model.RolePermission{
      RoleId:       roleId.String(),
      PermissionId: permId.String(),
      CreatedAt:    time.Now(),
    }
  })

  res, err := r.db.NewInsert().
    Model(&permissionModels).
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}

func (r *roleRepository) RemovePermissions(ctx context.Context, roleId types.Id, permissionIds ...types.Id) error {
  ctx, span := r.tracer.Start(ctx, "RoleRepository.RemovePermissions")
  defer span.End()

  permIds := sharedUtil.CastSlice(permissionIds, sharedUtil.ToString[types.Id])

  res, err := r.db.NewDelete().
    Model(types.Nil[model.RolePermission]()).
    Where("role_id = ? AND permission_id IN (?)", roleId.String(), bun.In(permIds)).
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}

func (r *roleRepository) AddUser(ctx context.Context, userId types.Id, roleIds ...types.Id) error {
  ctx, span := r.tracer.Start(ctx, "RoleRepository.AddUser")
  defer span.End()

  userRoleModels := sharedUtil.CastSlice(roleIds, func(roleId types.Id) model.UserRole {
    return model.UserRole{
      UserId:    userId.String(),
      RoleId:    roleId.String(),
      CreatedAt: time.Now(),
    }
  })

  res, err := r.db.NewInsert().
    Model(&userRoleModels).
    Returning("NULL").
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}

func (r *roleRepository) RemoveUser(ctx context.Context, userId types.Id, roleIds ...types.Id) error {
  ctx, span := r.tracer.Start(ctx, "RoleRepository.RemoveUser")
  defer span.End()

  idModels := sharedUtil.CastSlice(roleIds, sharedUtil.ToString[types.Id])

  res, err := r.db.NewDelete().
    Model(types.Nil[model.UserRole]()).
    Where("role_id IN (?) AND user_id = ?", bun.In(idModels), userId.String()).
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}
