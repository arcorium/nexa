package types

import "regexp"

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9|._]+@[a-zA-Z0-9|-]+(\\.[a-zA-Z0-9]{2,})+$")

type Email string

func (e Email) underlying() string {
	return string(e)
}

func (e Email) Validate() bool {
	return emailRegex.MatchString(e.underlying())
}
