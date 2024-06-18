package types

import (
  "errors"
  "github.com/google/uuid"
  "nexa/shared/wrapper"
)

func IdFromString(id string) (Id, error) {
  uid, err := uuid.Parse(id)
  if err != nil {
    return NullId(), ErrMalformedUUID
  }
  return Id(uid), nil
}

func NewId() (Id, error) {
  uid, err := uuid.NewRandom()
  return Id(uid), err
}

func NewId2() Id {
  return wrapper.DropError(NewId()) // TODO: Panic on error?
  //return wrapper.PanicDropError(NewId())
}

func NullId() Id { return Id(uuid.UUID{}) }

type Id uuid.UUID

func (i Id) Underlying() uuid.UUID {
  return uuid.UUID(i)
}

func (i Id) Equal(uuid string) bool {
  return i.Underlying().String() == uuid
}

var ErrMalformedUUID = errors.New("value has malformed format for an UUID")
