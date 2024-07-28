package pg

import (
  "context"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  "github.com/arcorium/nexa/shared/util/repo"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "github.com/uptrace/bun"
  "go.opentelemetry.io/otel/trace"
  entity "nexa/services/mailer/internal/domain/entity"
  "nexa/services/mailer/internal/domain/repository"
  "nexa/services/mailer/internal/infra/repository/model"
  "nexa/services/mailer/util"
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

func (m *mailRepository) Get(ctx context.Context, query repo.QueryParameter) (repo.PaginatedResult[entity.Mail], error) {
  ctx, span := m.tracer.Start(ctx, "MailRepository.Get")
  defer span.End()

  var models []model.Mail
  count, err := m.db.NewSelect().
    Model(&models).
    Offset(int(query.Offset)).
    Limit(int(query.Limit)).
    Relation("Tags").
    OrderExpr("sent_at DESC").
    ScanAndCount(ctx)

  // Checking and Mapping
  result := repo.CheckPaginationResult(models, count, err)
  if result.IsError() {
    spanUtil.RecordError(result.Err, span)
    return repo.NewPaginatedResult[entity.Mail](nil, uint64(count)), result.Err
  }

  users, ierr := sharedUtil.CastSliceErrsP(result.Data, repo.ToDomainErr[*model.Mail, entity.Mail])
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return repo.NewPaginatedResult[entity.Mail](nil, uint64(count)), ierr
  }

  return repo.NewPaginatedResult(users, uint64(count)), nil
}

func (m *mailRepository) FindByIds(ctx context.Context, ids ...types.Id) ([]entity.Mail, error) {
  ctx, span := m.tracer.Start(ctx, "MailRepository.FindByIds")
  defer span.End()

  mailIds := sharedUtil.CastSlice(ids, sharedUtil.ToString[types.Id])
  var models []model.Mail
  err := m.db.NewSelect().
    Model(&models).
    Where("id IN (?)", bun.In(mailIds)).
    Relation("Tags").
    Distinct().
    OrderExpr("sent_at DESC").
    Scan(ctx)

  res := repo.CheckSliceResult(models, err)
  if res.IsError() {
    spanUtil.RecordError(res.Err, span)
    return nil, res.Err
  }

  result, ierr := sharedUtil.CastSliceErrsP(res.Data, repo.ToDomainErr[*model.Mail, entity.Mail])
  if !ierr.IsNil() {
    return nil, ierr
  }
  return result, nil
}

func (m *mailRepository) FindByTag(ctx context.Context, tag types.Id) ([]entity.Mail, error) {
  ctx, span := m.tracer.Start(ctx, "MailRepository.FindByTag")
  defer span.End()

  var models []model.MailTag

  err := m.db.NewSelect().
    Model(&models).
    Relation("Mail.Tags").
    Where("tag_id = ?", tag.String()).
    OrderExpr("mail.sent_at DESC").
    Scan(ctx)

  res := repo.CheckSliceResult(models, err)
  if res.IsError() {
    spanUtil.RecordError(res.Err, span)
    return nil, res.Err
  }

  result, ierr := sharedUtil.CastSliceErrsP(res.Data, func(mailTag *model.MailTag) (entity.Mail, error) {
    return mailTag.Mail.ToDomain()
  })
  if !ierr.IsNil() {
    return nil, ierr
  }
  return result, nil
}

func (m *mailRepository) Create(ctx context.Context, mails ...entity.Mail) error {
  ctx, span := m.tracer.Start(ctx, "MailRepository.Create")
  defer span.End()

  models := sharedUtil.CastSliceP(mails, func(mail *entity.Mail) model.Mail {
    return model.FromMailDomain(mail, func(domain *entity.Mail, mail *model.Mail) {
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
    Where("mail_id = ? AND tag_id IN (?)", mailId.String(), bun.In(ids)).
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}

func (m *mailRepository) Patch(ctx context.Context, mail *entity.Mail) error {
  ctx, span := m.tracer.Start(ctx, "MailRepository.Patch")
  defer span.End()

  models := model.FromMailDomain(mail, func(domain *entity.Mail, mail *model.Mail) {
    mail.UpdatedAt = time.Now()
  })

  res, err := m.db.NewUpdate().
    Model(&models).
    WherePK().
    OmitZero().
    Column("status", "updated_at", "delivered_at").
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

func (m *mailRepository) AppendMultipleTags(ctx context.Context, mailTags ...repository.MailTags) error {
  ctx, span := m.tracer.Start(ctx, "MailRepository.AppendMultipleTags")
  defer span.End()

  models := model.FromMailTags(mailTags...)

  res, err := m.db.NewInsert().
    Model(&models).
    Returning("NULL").
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}
