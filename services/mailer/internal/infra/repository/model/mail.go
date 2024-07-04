package model

import (
  "database/sql"
  "github.com/uptrace/bun"
  domain "nexa/services/mailer/internal/domain/entity"
  "nexa/shared/types"
  "nexa/shared/util"
  "nexa/shared/util/repo"
  "nexa/shared/variadic"
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
  bun.BaseModel `bun:"table:mail"`

  Id        string        `bun:",type:uuid,pk,nullzero"`
  Subject   string        `bun:",nullzero"`
  Recipient string        `bun:",nullzero,notnull"`
  Sender    string        `bun:",nullzero,notnull"`
  Status    sql.NullInt64 `bun:",type:smallint,nullzero,notnull"`

  SentAt      time.Time `bun:",nullzero,notnull"`
  DeliveredAt time.Time `bun:",nullzero"`
  UpdatedAt   time.Time `bun:",nullzero"`

  Tags []Tag `bun:"m2m:mail_tags,join:Mail=Tags"`
}

func (p *Mail) ToDomain() (domain.Mail, error) {
  tags, ierr := util.CastSliceErrsP(p.Tags, func(from *Tag) (domain.Tag, error) {
    return from.ToDomain()
  })
  if !ierr.IsNil() {
    return domain.Mail{}, ierr
  }

  mailId, err := types.IdFromString(p.Id)
  if err != nil {
    return domain.Mail{}, err
  }

  sender, err := types.EmailFromString(p.Sender)
  if err != nil {
    return domain.Mail{}, err
  }

  recipient, err := types.EmailFromString(p.Sender)
  if err != nil {
    return domain.Mail{}, err
  }

  status, err := domain.NewStatus(uint8(p.Status.Int64))
  if err != nil {
    return domain.Mail{}, err
  }

  return domain.Mail{
    Id:        mailId,
    Subject:   p.Subject,
    Recipient: recipient,
    Sender:    sender,
    Status:    status,
    Tags:      tags,
  }, nil
}
