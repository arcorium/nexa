package entity

import (
	"nexa/shared/types"
	"time"
)

type Token struct {
	Token     string
	UserId    types.Id
	Usage     TokenUsage
	ExpiredAt time.Time
}

func (t *Token) IsExpired() bool {
	return t.ExpiredAt.After(time.Now())
}
