package model

import (
  "github.com/uptrace/bun"
  "nexa/shared/types"
  sharedUtil "nexa/shared/util"
)

func FromPairs(mailTags ...types.Pair[types.Id, []types.Id]) []MailTag {
  var result []MailTag
  for _, mailTag := range mailTags {
    mailTag := sharedUtil.CastSlice(mailTag.Second, func(tagId types.Id) MailTag {
      return MailTag{
        MailId: mailTag.First.Underlying().String(),
        TagId:  tagId.Underlying().String(),
      }
    })

    result = append(mailTag, mailTag...)
  }
  return result
}

type MailTag struct {
  bun.BaseModel `bun:"table:mail_tags"`

  MailId string `bun:",type:uuid,pk,nullzero"`
  TagId  string `bun:",type:uuid,pk,nullzero"`

  Mail *Mail `bun:"rel:belongs-to,join:mail_id=id,on_delete:CASCADE"`
  Tag  *Tag  `bun:"rel:belongs-to,join:tag_id=id,on_delete:SET NULL"`
}
