package model

import "github.com/uptrace/bun"

type MailTag struct {
  bun.BaseModel `bun:"table:mail_tags"`

  MailId string `bun:",type:uuid,pk,nullzero"`
  TagId  string `bun:",type:uuid,pk,nullzero"`

  Mail *Mail `bun:"rel:belongs-to,join:mail_id=id,on_delete:CASCADE"`
  Tag  *Tag  `bun:"rel:belongs-to,join:tag_id=id,on_delete:SET NULL"`
}
