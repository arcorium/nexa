package model

import (
  "github.com/uptrace/bun"
  domain "nexa/services/mailer/internal/domain/entity"
  "nexa/shared/types"
  "nexa/shared/util"
  "nexa/shared/util/repo"
  "nexa/shared/variadic"
  "nexa/shared/wrapper"
  "time"
)

type MailMapOption = repo.DataAccessModelMapOption[*domain.Mail, *Mail]

func FromMailDomain(domain *domain.Mail, opts ...MailMapOption) Mail {
  obj := Mail{
    Id:        domain.Id.Underlying().String(),
    Subject:   domain.Subject,
    Recipient: domain.Recipient.Underlying(),
    Sender:    domain.Sender.Underlying(),
    Status:    domain.Status.Underlying(),
  }

  variadic.New(opts...).DoAll(repo.MapOptionFunc(domain, &obj))

  return obj
}

type Mail struct {
  bun.BaseModel `bun:"table:mail"`

  Id        string `bun:",type:uuid,pk,nullzero"`
  Subject   string `bun:",nullzero"`
  Recipient string `bun:",nullzero,notnull"`
  Sender    string `bun:",nullzero,notnull"`
  Status    uint8  `bun:",type:uuid,nullzero,notnull"`

  SendAt      time.Time `bun:",nullzero,notnull"`
  DeliveredAt time.Time `bun:",nullzero"`
  UpdatedAt   time.Time `bun:",nullzero"`

  Tags []Tag `bun:"m2m:mail_tags,join:Mail=Tags"`
}

func (p *Mail) ToDomain() domain.Mail {
  tags := util.CastSlice(p.Tags, func(from *Tag) domain.Tag {
    return from.ToDomain()
  })

  return domain.Mail{
    Id:        wrapper.DropError(types.IdFromString(p.Id)),
    Subject:   p.Subject,
    Recipient: wrapper.DropError(types.EmailFromString(p.Recipient)),
    Sender:    wrapper.DropError(types.EmailFromString(p.Sender)),
    Status:    domain.Status(p.Status),
    Tags:      tags,
  }
}
