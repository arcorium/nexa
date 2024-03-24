package types

import (
	"errors"
	"regexp"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9|._]+@[a-zA-Z0-9|-]+(\\.[a-zA-Z0-9]{2,})+$")

func EmailFromString(email string) Email {
	return Email(email)
}

type Email string

func (e Email) Underlying() string {
	return string(e)
}

func (e Email) Validate() error {
	if !emailRegex.MatchString(e.Underlying()) {
		return ErrEmailMalformed
	}
	return nil
}

var ErrEmailMalformed = errors.New("email malformed")
