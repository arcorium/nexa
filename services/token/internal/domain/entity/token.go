package entity

import (
  sharedErr "github.com/arcorium/nexa/shared/errors"
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/token/util"
  "time"
)

func NewTokenType(val uint8) (TokenType, error) {
  tokenType := TokenType(val)
  if !tokenType.Valid() {
    return TokenTypeUnknown, sharedErr.ErrEnumOutOfBounds
  }
  return tokenType, nil
}

type TokenType uint8

const (
  TokenTypeString TokenType = iota
  TokenTypePIN
  TokenTypeAlphanumericPIN
  TokenTypeUnknown
)

func (t TokenType) Underlying() uint8 {
  return uint8(t)
}

func (t TokenType) Valid() bool {
  return t.Underlying() < TokenTypeUnknown.Underlying()
}

func NewTokenUsage(val uint8) (TokenUsage, error) {
  usage := TokenUsage(val)
  if !usage.Valid() {
    return TokenUsageUnknown, sharedErr.ErrEnumOutOfBounds
  }
  return usage, nil
}

type TokenUsage uint8

const (
  TokenUsageEmailVerification TokenUsage = iota
  TokenUsageResetPassword
  TokenUsageLogin
  TokenUsageGeneral // Could be used for something that doesn't covered above
  TokenUsageUnknown
)

func (t TokenUsage) Underlying() uint8 {
  return uint8(t)
}

func (t TokenUsage) Valid() bool {
  return t.Underlying() < TokenUsageUnknown.Underlying()
}

func NewToken(userId types.Id, expiryTime time.Duration, usage TokenUsage, tokenType TokenType, tokenLength uint16) Token {
  var token string
  if tokenType == TokenTypeString {
    token = util.RandomString(tokenLength)
  } else if tokenType == TokenTypePIN {
    token = util.RandomPIN(false, tokenLength)
  } else if tokenType == TokenTypeAlphanumericPIN {
    token = util.RandomPIN(true, tokenLength)
  }

  return Token{
    Token:     token,
    UserId:    userId,
    Usage:     usage,
    ExpiredAt: time.Now().Add(expiryTime),
  }
}

type Token struct {
  Token     string
  UserId    types.Id
  Usage     TokenUsage
  ExpiredAt time.Time
}

func (t *Token) IsExpired() bool {
  return t.ExpiredAt.Before(time.Now())
}
