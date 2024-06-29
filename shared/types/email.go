package types

import (
  "errors"
  "regexp"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9|._]+@[a-zA-Z0-9|-]+(\\.[a-zA-Z0-9]{2,})+$")

func EmailFromString(email string) (Email, error) {
  if !emailRegex.MatchString(email) {
    return "", ErrEmailMalformed
  }
  return Email(email), nil
}

type Email string

func (e Email) Underlying() string {
  return string(e)
}

func (e Email) String() string {
  return e.Underlying()
}

var ErrEmailMalformed = errors.New("email malformed")
