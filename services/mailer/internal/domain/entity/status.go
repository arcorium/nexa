package entity

import (
  sharedErr "nexa/shared/errors"
)

type Status uint8

func NewStatus(status uint8) (Status, error) {
  s := Status(status)
  if !s.Valid() {
    return s, sharedErr.ErrEnumOutOfBounds
  }
  return s, nil
}

const (
  StatusPending Status = iota
  StatusSending
  StatusDelivered
  StatusFailed
)

func (s Status) Underlying() uint8 {
  return uint8(s)
}

func (s Status) Valid() bool {
  return s <= StatusFailed
}

func (s Status) String() string {
  switch s {
  case StatusPending:
    return "Pending"
  case StatusSending:
    return "Sending"
  case StatusDelivered:
    return "Delivered"
  case StatusFailed:
    return "Failed"
  }
  return "Unknown"
}
