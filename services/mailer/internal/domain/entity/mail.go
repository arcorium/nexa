package entity

import (
  "nexa/shared/types"
)

type MailBodyType uint8

const (
  BodyTypeHTML MailBodyType = iota
  BodyTypePlain
)

func (m MailBodyType) String() string {
  switch m {
  case BodyTypeHTML:
    return "text/html"
  case BodyTypePlain:
    return "text/plain"
  }
  return "text/plain"
}

type Mail struct {
  Id        types.Id
  Subject   string
  Recipient types.Email
  Sender    types.Email

  BodyType MailBodyType
  Body     string

  Status Status
  Tags   []Tag
}
