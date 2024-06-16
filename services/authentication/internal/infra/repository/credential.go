package repository

import (
  "context"
  "github.com/uptrace/bun"
  "nexa/services/authentication/internal/domain/repository"
  "nexa/services/authentication/internal/infra/model"
  "nexa/services/authentication/shared/domain/entity"
  "nexa/shared/types"
  "nexa/shared/util"
  "nexa/shared/util/repo"
  "time"
)

func NewCredential(db bun.IDB) repository.ICredential {
  return &credentialRepository{db: db}
}

type credentialRepository struct {
  db bun.IDB
}

func (c *credentialRepository) Create(ctx context.Context, credential *entity.Credential) error {
  dbModel := model.FromCredentialModel(credential, func(domain *entity.Credential, credential *model.Credential) {
    credential.CreatedAt = time.Now()
  })

  res, err := c.db.NewInsert().
    Model(&dbModel).
    Exec(ctx)

  return repo.CheckResult(res, err)
}

func (c *credentialRepository) Patch(ctx context.Context, credential *entity.Credential) error {
  dbModel := model.FromCredentialModel(credential, func(domain *entity.Credential, credential *model.Credential) {
    credential.UpdatedAt = time.Now()
  })

  res, err := c.db.NewUpdate().
    Model(&dbModel).
    WherePK().
    OmitZero().
    Exec(ctx)

  return repo.CheckResult(res, err)
}

func (c *credentialRepository) Delete(ctx context.Context, accessTokenId ...types.Id) error {
  res, err := c.db.NewDelete().
    Model(util.Nil[model.Credential]()).
    Where("access_token_id IN (?)", bun.In(accessTokenId)).
    Exec(ctx)

  return repo.CheckResult(res, err)
}

func (c *credentialRepository) DeleteByUserId(ctx context.Context, userId types.Id) error {
  res, err := c.db.NewDelete().
    Model(util.Nil[model.Credential]()).
    Where("user_id = ?", userId.Underlying().String()).
    Exec(ctx)

  return repo.CheckResult(res, err)
}

func (c *credentialRepository) Find(ctx context.Context, accessTokenId types.Id) (entity.Credential, error) {
  var dbModel model.Credential

  err := c.db.NewSelect().
    Model(&dbModel).
    Where("access_token_id = ?", accessTokenId.Underlying().String()).
    Scan(ctx)

  return dbModel.ToDomain(), err
}

func (c *credentialRepository) FindByUserId(ctx context.Context, userId types.Id) ([]entity.Credential, error) {
  var dbModels []model.Credential

  err := c.db.NewSelect().
    Model(&dbModels).
    Where("user_id = ?", userId.Underlying().String()).
    Scan(ctx)

  credentials := util.CastSlice(dbModels, func(from *model.Credential) entity.Credential {
    return from.ToDomain()
  })

  return repo.CheckSliceResult(credentials, err).Value()
}

func (c *credentialRepository) FindAll(ctx context.Context, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Credential], error) {
  var dbModels []model.Credential

  count, err := c.db.NewSelect().
    Model(&dbModels).
    Limit(int(parameter.Limit)).
    Offset(int(parameter.Offset)).
    ScanAndCount(ctx)

  result := repo.CheckPaginationResult(dbModels, count, err)
  credentials := util.CastSlice(dbModels, func(from *model.Credential) entity.Credential {
    return from.ToDomain()
  })

  return repo.NewPaginatedResult(credentials, uint64(count)), result.Err
}
