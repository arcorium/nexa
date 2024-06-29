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
  //return wrapper.Must(NewId())
}

func NullId() Id { return Id(uuid.UUID{}) }

type Id uuid.UUID

func (i Id) Underlying() uuid.UUID {
  return uuid.UUID(i)
}

func (i Id) String() string {
  return i.Underlying().String()
}

func (i Id) EqWithString(uuid string) bool {
  return i.Underlying().String() == uuid
}

func (i Id) Eq(other Id) bool {
  return i.Underlying().String() == other.Underlying().String()
}

var ErrMalformedUUID = errors.New("value has malformed format for an UUID")

const NullIdStr = "00000000-0000-0000-0000-000000000000"
