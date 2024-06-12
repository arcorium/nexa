package pg

import (
  "context"
  "github.com/uptrace/bun"
  "nexa/services/user/internal/domain/entity"
  "nexa/services/user/internal/domain/repository"
  "nexa/services/user/internal/infra/model"
  "nexa/shared/types"
  "nexa/shared/util"
  "nexa/shared/util/repo"
  "time"
)

func NewProfile(db bun.IDB) repository.IProfile {
  return &profileRepository{db: db}
}

type profileRepository struct {
  db bun.IDB
}

func (p profileRepository) Create(ctx context.Context, profile *entity.Profile) error {
  dbModel := model.FromProfileDomain(profile)

  res, err := p.db.NewInsert().
    Model(&dbModel).
    Returning("NULL").
    Exec(ctx)

  return repo.CheckResult(res, err)
}

func (p profileRepository) FindByIds(ctx context.Context, userIds ...types.Id) ([]entity.Profile, error) {
  var dbModel []model.Profile

  ids := util.CastSlice(userIds, func(from *types.Id) string {
    return from.Underlying().String()
  })

  err := p.db.NewSelect().
    Model(&dbModel).
    Where("user_id IN (?)", bun.In(ids)).
    Scan(ctx)

  result := repo.CheckSliceResult(dbModel, err)
  profiles := util.CastSlice(result.Data, func(from *model.Profile) entity.Profile {
    return from.ToDomain()
  })
  return profiles, result.Err
}

func (p profileRepository) Update(ctx context.Context, profile *entity.Profile) error {
  dbModel := model.FromProfileDomain(profile, func(domain *entity.Profile, profile *model.Profile) {
    profile.UpdatedAt = time.Now()
  })

  res, err := p.db.NewUpdate().
    Model(&dbModel).
    WherePK().
    ExcludeColumn("user_id").
    Exec(ctx)

  return repo.CheckResult(res, err)
}

func (p profileRepository) Patch(ctx context.Context, profile *entity.Profile) error {
  dbModel := model.FromProfileDomain(profile, func(domain *entity.Profile, profile *model.Profile) {
    profile.UpdatedAt = time.Now()
  })

  res, err := p.db.NewUpdate().
    Model(&dbModel).
    WherePK().
    OmitZero().
    ExcludeColumn("user_id").
    Exec(ctx)

  return repo.CheckResult(res, err)
}
