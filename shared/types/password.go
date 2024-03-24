package types

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

// PasswordFromString make password type from string, the parameter could be hash or plain text
// For plain text call Hash() method before calling Equal() method
func PasswordFromString(password string) Password {
	return Password(password)
}

type Password string

func (h Password) Underlying() string {
	return string(h)
}

func (h Password) Equal(plainPassword string) error {
	if bcrypt.CompareHashAndPassword([]byte(h.Underlying()), []byte(plainPassword)) != nil {
		return ErrPasswordDifferent
	}
	return nil
}

func (h Password) Hash() (Password, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(h.Underlying()), bcrypt.DefaultCost)
	if err != nil {
		return "", ErrHashPlainString
	}
	return Password(hashed), nil
}

var ErrHashPlainString = errors.New("failed to hash plain text")
var ErrPasswordDifferent = errors.New("password compared is different")
