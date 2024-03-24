package types

import (
	"github.com/google/uuid"
)

func IdFromString(id string) Id {
	return Id([]byte(id))
}

func NewId() Id {
	return Id(uuid.New())
}

type Id uuid.UUID

func (i Id) Validate() bool {
	_, err := uuid.Parse(i.Underlying().String())
	return err == nil
}

func (i Id) Underlying() uuid.UUID {
	return uuid.UUID(i)
}
