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

func NewPermission(db bun.IDB) repository.IPermission {
	return &permissionRepository{db: db}
}

type permissionRepository struct {
	db bun.IDB
}

func (p *permissionRepository) FindById(ctx context.Context, id types.Id) (entity.Permission, error) {
	var dbModel model.Permission
	err := p.db.NewSelect().
		Model(&dbModel).
		Where("id = ?", id.Underlying().String()).
		Scan(ctx)

	if err != nil {
		return entity.Permission{}, err
	}

	return dbModel.ToDomain(), nil
}

func (p *permissionRepository) FindByUserId(ctx context.Context, userId types.Id) ([]entity.Permission, error) {
	var dbModels []model.Permission
	err := p.db.NewSelect().
		Model(util.Nil[model.UserRole]()).
		Relation("Role").
		Relation("Role.Permissions").
		Where("user_id = ?", userId.Underlying().String()).
		Scan(ctx, &dbModels)

	result := repo.CheckSliceResult(dbModels, err)
	permissions := util.CastSlice(result.Data, func(from *model.Permission) entity.Permission {
		return from.ToDomain()
	})
	return permissions, result.Err
}

func (p *permissionRepository) FindAll(ctx context.Context, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Permission], error) {
	var dbModels []model.Permission

	count, err := p.db.NewSelect().
		Model(&dbModels).
		Limit(int(parameter.Limit)).
		Offset(int(parameter.Offset)).
		ScanAndCount(ctx)

	result := repo.CheckPaginationResult(dbModels, count, err)
	permissions := util.CastSlice(result.Data, func(from *model.Permission) entity.Permission {
		return from.ToDomain()
	})

	return repo.NewPaginatedResult(permissions, uint64(count)), result.Err
}

func (p *permissionRepository) Create(ctx context.Context, permission *entity.Permission) error {
	dbModel := model.FromPermissionDomain(permission, func(domain *entity.Permission, permission *model.Permission) {
		permission.CreatedAt = time.Now()
	})

	res, err := p.db.NewInsert().
		Model(&dbModel).
		Exec(ctx)

	return repo.CheckResult(res, err)
}

func (p *permissionRepository) Delete(ctx context.Context, id types.Id) error {
	res, err := p.db.NewDelete().
		Model(util.Nil[model.Permission]()).
		Where("id = ?", id.Underlying().String()).
		Exec(ctx)

	return repo.CheckResult(res, err)
}
