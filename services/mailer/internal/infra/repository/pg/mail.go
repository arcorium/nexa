package pg

import (
  "context"
  "github.com/uptrace/bun"
  "go.opentelemetry.io/otel/trace"
  domain "nexa/services/mailer/internal/domain/entity"
  "nexa/services/mailer/internal/domain/repository"
  "nexa/services/mailer/internal/infra/repository/model"
  "nexa/services/mailer/util"
  "nexa/shared/types"
  sharedUtil "nexa/shared/util"
  "nexa/shared/util/repo"
  "time"
)

func NewMail(db bun.IDB) repository.IMail {
  return &mailRepository{
    db:     db,
    tracer: util.GetTracer(),
  }
}

type mailRepository struct {
  db     bun.IDB
  tracer trace.Tracer
}

func (m *mailRepository) FindAll(ctx context.Context, query repo.QueryParameter) (repo.PaginatedResult[domain.Mail], error) {
  ctx, span := m.tracer.Start(ctx, "MailRepository.FindAll")
  defer span.End()

  var models []model.Mail
  count, err := m.db.NewSelect().
    Model(&models).
    Offset(int(query.Offset)).
    Limit(int(query.Limit)).
    ScanAndCount(ctx)

  // Checking and Mapping
  result := repo.CheckPaginationResultWithSpan(models, count, err, span)
  users, ierr := sharedUtil.CastSliceErrsP(result.Data, repo.ToDomainErr[*model.Mail, domain.Mail])
  if !ierr.IsNil() {
    return repo.PaginatedResult[domain.Mail]{}, ierr
  }

  return repo.NewPaginatedResult(users, uint64(count)), result.Err
}

func (m *mailRepository) FindByIds(ctx context.Context, ids ...types.Id) ([]domain.Mail, error) {
  ctx, span := m.tracer.Start(ctx, "MailRepository.FindByIds")
  defer span.End()

  mailIds := sharedUtil.CastSlice(ids, func(from types.Id) string {
    return from.Underlying().String()
  })

  var models []model.Mail

  err := m.db.NewSelect().
    Model(&models).
    Where("id IN (?)", bun.In(mailIds)).
    Scan(ctx)

  res := repo.CheckSliceResultWithSpan(models, err, span)

  result, ierr := sharedUtil.CastSliceErrsP(res.Data, repo.ToDomainErr[*model.Mail, domain.Mail])
  if !ierr.IsNil() {
    return nil, ierr
  }
  return result, res.Err
}

func (m *mailRepository) FindByTag(ctx context.Context, tag types.Id) ([]domain.Mail, error) {
  ctx, span := m.tracer.Start(ctx, "MailRepository.FindByTag")
  defer span.End()

  var models []model.Mail

  err := m.db.NewSelect().
    Model(&models).
    Relation("Tags").
    Where("tags.id = ?", tag.String()).
    Scan(ctx)

  res := repo.CheckSliceResultWithSpan(models, err, span)
  result, ierr := sharedUtil.CastSliceErrsP(res.Data, repo.ToDomainErr[*model.Mail, domain.Mail])
  if !ierr.IsNil() {
    return nil, ierr
  }
  return result, res.Err
}

func (m *mailRepository) Create(ctx context.Context, mails ...domain.Mail) error {
  ctx, span := m.tracer.Start(ctx, "MailRepository.Create")
  defer span.End()

  models := sharedUtil.CastSliceP(mails, func(mail *domain.Mail) model.Mail {
    return model.FromMailDomain(mail, func(domain *domain.Mail, mail *model.Mail) {
    })
  })

  res, err := m.db.NewInsert().
    Model(&models).
    Returning("NULL").
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}

func (m *mailRepository) AppendTags(ctx context.Context, mailId types.Id, tagIds []types.Id) error {
  ctx, span := m.tracer.Start(ctx, "MailRepository.AppendTags")
  defer span.End()

  models := sharedUtil.CastSlice(tagIds, func(tagId types.Id) model.MailTag {
    return model.MailTag{
      MailId: mailId.Underlying().String(),
      TagId:  tagId.Underlying().String(),
    }
  })

  res, err := m.db.NewInsert().
    Model(&models).
    Returning("NULL").
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}

func (m *mailRepository) RemoveTags(ctx context.Context, mailId types.Id, tagIds []types.Id) error {
  ctx, span := m.tracer.Start(ctx, "MailRepository.RemoveTags")
  defer span.End()

  ids := sharedUtil.CastSlice(tagIds, func(tagId types.Id) string {
    return tagId.Underlying().String()
  })

  res, err := m.db.NewDelete().
    Model(types.Nil[model.MailTag]()).
    Where("id = ? AND tag_id IN (?)", mailId.Underlying().String(), bun.In(ids)).
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}

func (m *mailRepository) Patch(ctx context.Context, mail *domain.Mail) error {
  ctx, span := m.tracer.Start(ctx, "MailRepository.Patch")
  defer span.End()

  models := model.FromMailDomain(mail, func(domain *domain.Mail, mail *model.Mail) {
    mail.UpdatedAt = time.Now()
  })

  res, err := m.db.NewUpdate().
    Model(&models).
    WherePK().
    OmitZero().
    ExcludeColumn("id").
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}

func (m *mailRepository) Remove(ctx context.Context, id types.Id) error {
  ctx, span := m.tracer.Start(ctx, "MailRepository.Remove")
  defer span.End()

  res, err := m.db.NewDelete().
    Model(types.Nil[model.Mail]()).
    Where("id = ?", id.String()).
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}

func (m *mailRepository) AppendMultipleTags(ctx context.Context, mailTags ...types.Pair[types.Id, []types.Id]) error {
  ctx, span := m.tracer.Start(ctx, "MailRepository.AppendMultipleTags")
  defer span.End()

  models := model.FromPairs(mailTags...)

  res, err := m.db.NewInsert().
    Model(&models).
    Returning("NULL").
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}
