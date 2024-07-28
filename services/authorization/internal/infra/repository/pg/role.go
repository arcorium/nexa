package pg

import (
  "context"
  "errors"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  "github.com/arcorium/nexa/shared/util/repo"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "github.com/uptrace/bun"
  "github.com/uptrace/bun/driver/pgdriver"
  "go.opentelemetry.io/otel/trace"
  "nexa/services/authorization/internal/domain/entity"
  "nexa/services/authorization/internal/domain/repository"
  "nexa/services/authorization/internal/infra/repository/model"
  "nexa/services/authorization/util"
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
    Distinct().
    OrderExpr("created_at DESC").
    Scan(ctx)

  result := repo.CheckSliceResult(dbModels, err)
  if result.IsError() {
    spanUtil.RecordError(result.Err, span)
    return nil, result.Err
  }

  roles, ierr := sharedUtil.CastSliceErrsP(dbModels, repo.ToDomainErr[*model.Role, entity.Role])
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return nil, ierr
  }

  return roles, nil
}

func (r *roleRepository) FindByUserId(ctx context.Context, userId types.Id) ([]entity.Role, error) {
  ctx, span := r.tracer.Start(ctx, "RoleRepository.FindByUserId")
  defer span.End()

  var dbModels []model.UserRole

  err := r.db.NewSelect().
    Model(&dbModels).
    Where("user_id = ?", userId.String()).
    Relation("Role.Permissions").
    OrderExpr("created_at DESC").
    Scan(ctx)

  result := repo.CheckSliceResultWithSpan(dbModels, err, span)
  if result.IsError() {
    spanUtil.RecordError(result.Err, span)
    return nil, result.Err
  }

  roles, ierr := sharedUtil.CastSliceErrsP(dbModels, func(userRole *model.UserRole) (entity.Role, error) {
    return userRole.Role.ToDomain()
  })
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return nil, ierr
  }

  return roles, nil
}

func (r *roleRepository) FindByName(ctx context.Context, name string) (entity.Role, error) {
  ctx, span := r.tracer.Start(ctx, "RoleRepository.FindByName")
  defer span.End()

  var dbModels model.Role
  err := r.db.NewSelect().
    Model(&dbModels).
    Where("name = ?", name).
    Relation("Permissions").
    Scan(ctx)

  if err != nil {
    spanUtil.RecordError(err, span)
    return entity.Role{}, err
  }

  return dbModels.ToDomain()
}

func (r *roleRepository) Get(ctx context.Context, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Role], error) {
  ctx, span := r.tracer.Start(ctx, "RoleRepository.Get")
  defer span.End()

  var dbModels []model.Role
  count, err := r.db.NewSelect().
    Model(&dbModels).
    Relation("Permissions").
    Offset(int(parameter.Offset)).
    Limit(int(parameter.Limit)).
    OrderExpr("created_at DESC").
    ScanAndCount(ctx)

  result := repo.CheckPaginationResult(dbModels, count, err)
  if result.IsError() {
    spanUtil.RecordError(result.Err, span)
    return repo.NewPaginatedResult[entity.Role](nil, uint64(count)), result.Err
  }

  roles, ierr := sharedUtil.CastSliceErrsP(dbModels, repo.ToDomainErr[*model.Role, entity.Role])
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return repo.NewPaginatedResult[entity.Role](nil, uint64(count)), ierr
  }

  return repo.NewPaginatedResult(roles, uint64(count)), nil
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

func (r *roleRepository) Patch(ctx context.Context, role *entity.PatchedRole) error {
  ctx, span := r.tracer.Start(ctx, "RoleRepository.Patch")
  defer span.End()

  dbModel := model.FromPatchedRoleDomain(role, func(domain *entity.PatchedRole, role *model.Role) {
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
    Returning("NULL").
    Exec(ctx)

  pgErr, ok := err.(pgdriver.Error)
  if ok {
    if pgErr.IntegrityViolation() {
      err = errors.Join(pgErr, errors.New(pgErr.Field(68)))
      spanUtil.RecordError(err, span)
      return err
    }
  }

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

func (r *roleRepository) ClearPermission(ctx context.Context, roleId types.Id) error {
  ctx, span := r.tracer.Start(ctx, "RoleRepository.ClearPermission")
  defer span.End()

  res, err := r.db.NewDelete().
    Model(types.Nil[model.RolePermission]()).
    Where("role_id = ?", roleId.String()).
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

  pgErr, ok := err.(pgdriver.Error)
  if ok {
    if pgErr.IntegrityViolation() {
      err = errors.Join(pgErr, errors.New(pgErr.Field(68)))
      spanUtil.RecordError(err, span)
      return err
    }
  }

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

func (r *roleRepository) ClearUser(ctx context.Context, userId types.Id) error {
  ctx, span := r.tracer.Start(ctx, "RoleRepository.ClearUser")
  defer span.End()

  res, err := r.db.NewDelete().
    Model(types.Nil[model.UserRole]()).
    Where("user_id = ?", userId.String()).
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}
