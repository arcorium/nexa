package types

import (
	"golang.org/x/crypto/bcrypt"
)

type HashString string

func (h HashString) underlying() string {
	return string(h)
}

func (h HashString) Equal(val string) bool {
	return bcrypt.CompareHashAndPassword([]byte(h.underlying()), []byte(val)) != nil
}
