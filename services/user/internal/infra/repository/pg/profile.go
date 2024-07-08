package pg

import (
  "context"
  "github.com/uptrace/bun"
  "go.opentelemetry.io/otel/trace"
  "nexa/services/user/internal/domain/entity"
  "nexa/services/user/internal/domain/repository"
  "nexa/services/user/internal/infra/repository/model"
  "nexa/services/user/util"
  "nexa/shared/types"
  sharedUtil "nexa/shared/util"
  "nexa/shared/util/repo"
  spanUtil "nexa/shared/util/span"
  "time"
)

func NewProfile(db bun.IDB) repository.IProfile {
  return &profileRepository{
    db:     db,
    tracer: util.GetTracer(),
  }
}

type profileRepository struct {
  db bun.IDB

  tracer trace.Tracer
}

func (p profileRepository) Create(ctx context.Context, profile *entity.Profile) error {
  ctx, span := p.tracer.Start(ctx, "ProfileRepository.Create")
  defer span.End()

  dbModel := model.FromProfileDomain(profile)

  res, err := p.db.NewInsert().
    Model(&dbModel).
    Returning("NULL").
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}

func (p profileRepository) FindByIds(ctx context.Context, userIds ...types.Id) ([]entity.Profile, error) {
  ctx, span := p.tracer.Start(ctx, "ProfileRepository.FindByIds")
  defer span.End()

  ids := sharedUtil.CastSlice(userIds, sharedUtil.ToString[types.Id])
  var dbModel []model.Profile

  err := p.db.NewSelect().
    Model(&dbModel).
    Where("id IN (?)", bun.In(ids)).
    Distinct().
    OrderExpr("updated_at DESC").
    Scan(ctx)

  result := repo.CheckSliceResultWithSpan(dbModel, err, span)
  if result.Err != nil {
    return nil, result.Err
  }

  profiles, ierr := sharedUtil.CastSliceErrsP(result.Data, repo.ToDomainErr[*model.Profile, entity.Profile])
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return nil, ierr
  }

  return profiles, nil
}

// NOTE: Currently not used, because each user can only have single profile
func (p profileRepository) FindByUserId(ctx context.Context, userId types.Id) (*entity.Profile, error) {
  ctx, span := p.tracer.Start(ctx, "ProfileRepository.FindByIds")
  defer span.End()

  var dbModel model.Profile

  err := p.db.NewSelect().
    Model(&dbModel).
    Where("user_id = ?", userId.String()).
    Distinct().
    OrderExpr("updated_at DESC").
    Scan(ctx)

  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  profile, err := dbModel.ToDomain()
  if err != nil {
    return nil, err
  }

  return &profile, nil
}

func (p profileRepository) Update(ctx context.Context, profile *entity.Profile) error {
  ctx, span := p.tracer.Start(ctx, "ProfileRepository.Update")
  defer span.End()

  dbModel := model.FromProfileDomain(profile, func(domain *entity.Profile, profile *model.Profile) {
    profile.UpdatedAt = time.Now()
  })

  res, err := p.db.NewUpdate().
    Model(&dbModel).
    WherePK().
    ExcludeColumn("user_id", "id").
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}

func (p profileRepository) Patch(ctx context.Context, profile *entity.PatchedProfile) error {
  ctx, span := p.tracer.Start(ctx, "ProfileRepository.Patch")
  defer span.End()

  dbModel := model.FromPatchedProfileDomain(profile, func(domain *entity.PatchedProfile, profile *model.Profile) {
    profile.UpdatedAt = time.Now()
  })

  res, err := p.db.NewUpdate().
    Model(&dbModel).
    WherePK().
    OmitZero().
    ExcludeColumn("user_id", "id").
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}
