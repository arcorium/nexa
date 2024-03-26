package types

import (
	"errors"
	"github.com/google/uuid"
)

func IdFromString(id string) Id {
	return Id([]byte(id))
}

func NewId() Id {
	return Id(uuid.New())
}

type Id uuid.UUID

func (i Id) Validate() error {
	_, err := uuid.Parse(i.Underlying().String())
	if err != nil {
		return ErrMalformedUUID
	}
	return nil
}

func (i Id) Underlying() uuid.UUID {
	return uuid.UUID(i)
}

var ErrMalformedUUID = errors.New("value has malformed format for an UUID")
