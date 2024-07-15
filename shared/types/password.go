package types

import (
  "errors"
  "golang.org/x/crypto/bcrypt"
)

// PasswordFromString create password by hashing the parameter with bcrypt.DefaultCost
func PasswordFromString(password string) Password {
  return Password(password)
}

// Password hashed string
type Password string

func (h Password) Underlying() string {
  return string(h)
}

func (h Password) String() string {
  return h.Underlying()
}

func (h Password) Hash() (HashedPassword, error) {
  hashed, err := bcrypt.GenerateFromPassword([]byte(h.Underlying()), bcrypt.DefaultCost)
  if err != nil {
    return "", ErrHashPlainString
  }
  return HashedPassword(hashed), nil
}

type HashedPassword string

func (h HashedPassword) Underlying() string {
  return string(h)
}

func (h HashedPassword) String() string {
  return h.Underlying()
}

func (h HashedPassword) EqWithString(password string) error {
  if bcrypt.CompareHashAndPassword([]byte(h.Underlying()), []byte(password)) != nil {
    return ErrPasswordDifferent
  }
  return nil
}

func (h HashedPassword) Eq(password Password) error {
  return h.EqWithString(password.String())
}

var ErrHashPlainString = errors.New("failed to hash plain text")
var ErrPasswordDifferent = errors.New("password compared is different")
