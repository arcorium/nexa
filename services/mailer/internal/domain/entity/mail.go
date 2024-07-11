package entity

import (
  "github.com/arcorium/nexa/shared/types"
  "time"
)

type MailBodyType uint8

const (
  BodyTypeHTML MailBodyType = iota
  BodyTypePlain
  BodyTypeUnknown
)

func (m MailBodyType) Underlying() uint8 {
  return uint8(m)
}

func (m MailBodyType) String() string {
  switch m {
  case BodyTypeHTML:
    return "text/html"
  case BodyTypePlain:
    return "text/plain"
  }
  return "text/unknown"
}

type Mail struct {
  Id        types.Id
  Subject   string
  Recipient types.Email
  Sender    types.Email

  BodyType MailBodyType
  Body     string

  Status Status

  SentAt      time.Time
  DeliveredAt time.Time

  Tags []Tag
}
