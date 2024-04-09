package types

import (
  "errors"
  "golang.org/x/crypto/bcrypt"
)

// PasswordFromString create password by hashing the parameter with bcrypt.DefaultCost
func PasswordFromString(password string) (Password, error) {
  hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
  if err != nil {
    return "", ErrHashPlainString
  }
  return Password(hashed), nil
}

// Password hashed string
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

var ErrHashPlainString = errors.New("failed to hash plain text")
var ErrPasswordDifferent = errors.New("password compared is different")
