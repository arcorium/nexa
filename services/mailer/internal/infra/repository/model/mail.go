package model

import (
  "database/sql"
  sharedErr "github.com/arcorium/nexa/shared/errors"
  "github.com/arcorium/nexa/shared/types"
  "github.com/arcorium/nexa/shared/util"
  "github.com/arcorium/nexa/shared/util/repo"
  "github.com/arcorium/nexa/shared/variadic"
  "github.com/uptrace/bun"
  domain "nexa/services/mailer/internal/domain/entity"
  "time"
)

type MailMapOption = repo.DataAccessModelMapOption[*domain.Mail, *Mail]

func FromMailDomain(domain *domain.Mail, opts ...MailMapOption) Mail {
  mail := Mail{
    Id:          domain.Id.String(),
    Subject:     domain.Subject,
    Recipient:   domain.Recipient.Underlying(),
    Sender:      domain.Sender.Underlying(),
    Status:      sql.NullInt64{Int64: int64(domain.Status.Underlying()), Valid: domain.Status.Valid()},
    SentAt:      domain.SentAt,
    DeliveredAt: domain.DeliveredAt,
  }

  variadic.New(opts...).DoAll(repo.MapOptionFunc(domain, &mail))

  return mail
}

type Mail struct {
  bun.BaseModel `bun:"table:mails"`

  Id        string        `bun:",type:uuid,pk,nullzero"`
  Subject   string        `bun:",nullzero"`
  Recipient string        `bun:",nullzero,notnull"`
  Sender    string        `bun:",nullzero,notnull"`
  Status    sql.NullInt64 `bun:",type:smallint,nullzero,notnull"`

  SentAt      time.Time `bun:",nullzero,notnull"`
  DeliveredAt time.Time `bun:",nullzero"`
  UpdatedAt   time.Time `bun:",nullzero"`

  Tags []Tag `bun:"m2m:mail_tags,join:Mail=Tag"`
}

func (p *Mail) ToDomain() (domain.Mail, error) {
  var tags []domain.Tag
  if p.Tags != nil {
    var ierr sharedErr.IndicesError
    tags, ierr = util.CastSliceErrsP(p.Tags, repo.ToDomainErr[*Tag, domain.Tag])
    if !ierr.IsNil() {
      return domain.Mail{}, ierr
    }
  }

  mailId, err := types.IdFromString(p.Id)
  if err != nil {
    return domain.Mail{}, err
  }

  sender, err := types.EmailFromString(p.Sender)
  if err != nil {
    return domain.Mail{}, err
  }

  recipient, err := types.EmailFromString(p.Recipient)
  if err != nil {
    return domain.Mail{}, err
  }

  status, err := domain.NewStatus(uint8(p.Status.Int64))
  if err != nil {
    return domain.Mail{}, err
  }

  return domain.Mail{
    Id:          mailId,
    Subject:     p.Subject,
    Recipient:   recipient,
    Sender:      sender,
    Status:      status,
    SentAt:      p.SentAt,
    DeliveredAt: p.DeliveredAt,
    Tags:        tags,
  }, nil
}
