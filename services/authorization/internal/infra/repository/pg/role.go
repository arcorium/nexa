package pg

import (
	"context"
	"github.com/uptrace/bun"
	"nexa/services/authorization/internal/domain/repository"
	"nexa/services/authorization/internal/infra/model"
	"nexa/services/authorization/shared/domain/entity"
	"nexa/shared/types"
	"nexa/shared/util"
	"nexa/shared/util/repo"
	"time"
)

func NewRole(db bun.IDB) repository.IRole {
	return &roleRepository{db: db}
}

type roleRepository struct {
	db bun.IDB
}

func (r *roleRepository) FindByIds(ctx context.Context, id ...types.Id) ([]entity.Role, error) {
	var dbModels []model.Role

	err := r.db.NewSelect().
		Model(&dbModels).
		Where("id IN (?)", bun.In(id)).
		Relation("Permissions").
		Scan(ctx)

	result := repo.CheckSliceResult(dbModels, err)
	roles := util.CastSlice(dbModels, func(from *model.Role) entity.Role {
		return from.ToDomain()
	})
	return roles, result.Err
}

func (r *roleRepository) FindAll(ctx context.Context, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Role], error) {
	var dbModels []model.Role

	count, err := r.db.NewSelect().
		Model(&dbModels).
		Offset(int(parameter.Offset)).
		Limit(int(parameter.Limit)).
		ScanAndCount(ctx)

	result := repo.CheckPaginationResult(dbModels, count, err)
	roles := util.CastSlice(dbModels, func(from *model.Role) entity.Role {
		return from.ToDomain()
	})
	return repo.NewPaginatedResult(roles, uint64(count)), result.Err
}

func (r *roleRepository) Create(ctx context.Context, role *entity.Role) error {
	dbModel := model.FromRoleDomain(role, func(domain *entity.Role, role *model.Role) {
		role.CreatedAt = time.Now()
	})

	res, err := r.db.NewInsert().
		Model(&dbModel).
		Exec(ctx)

	return repo.CheckResult(res, err)
}

func (r *roleRepository) Patch(ctx context.Context, role *entity.Role) error {
	dbModel := model.FromRoleDomain(role, func(domain *entity.Role, role *model.Role) {
		role.UpdatedAt = time.Now()
	})

	res, err := r.db.NewUpdate().
		Model(&dbModel).
		WherePK().
		OmitZero().
		Exec(ctx)

	return repo.CheckResult(res, err)
}

func (r *roleRepository) Delete(ctx context.Context, id types.Id) error {
	res, err := r.db.NewDelete().
		Model(util.Nil[model.Role]()).
		Where("id = ?", id.Underlying().String()).
		Exec(ctx)

	return repo.CheckResult(res, err)
}

func (r *roleRepository) AddPermissions(ctx context.Context, roleId types.Id, permissionIds ...types.Id) error {
	permissionModels := util.CastSlice(permissionIds, func(from *types.Id) model.RolePermission {
		return model.RolePermission{
			RoleId:       roleId.Underlying().String(),
			PermissionId: from.Underlying().String(),
			CreatedAt:    time.Now(),
		}
	})

	res, err := r.db.NewInsert().
		Model(&permissionModels).
		Exec(ctx)

	return repo.CheckResult(res, err)
}

func (r *roleRepository) RemovePermissions(ctx context.Context, roleId types.Id, permissionIds ...types.Id) error {
	idModels := util.CastSlice(permissionIds, func(from *types.Id) string {
		return from.Underlying().String()
	})

	res, err := r.db.NewDelete().
		Model(util.Nil[model.RolePermission]()).
		Where("role_id = ? AND permission_id IN (?)", roleId.Underlying().String(), bun.In(idModels)).
		Exec(ctx)

	return repo.CheckResult(res, err)
}

func (r *roleRepository) AddUser(ctx context.Context, userId types.Id, roleIds ...types.Id) error {
	userRoleModels := util.CastSlice(roleIds, func(from *types.Id) model.UserRole {
		return model.UserRole{
			UserId:    userId.Underlying().String(),
			RoleId:    from.Underlying().String(),
			CreatedAt: time.Now(),
		}
	})

	res, err := r.db.NewInsert().
		Model(&userRoleModels).
		Exec(ctx)

	return repo.CheckResult(res, err)
}

func (r *roleRepository) RemoveUser(ctx context.Context, userId types.Id, roleIds ...types.Id) error {
	idModels := util.CastSlice(roleIds, func(from *types.Id) string {
		return from.Underlying().String()
	})

	res, err := r.db.NewDelete().
		Model(util.Nil[model.UserRole]()).
		Where("role_id IN (?) AND user_id = ?", idModels, userId.Underlying().String()).
		Exec(ctx)

	return repo.CheckResult(res, err)
}
