package entity

import (
  sharedErr "nexa/shared/errors"
  sharedJwt "nexa/shared/jwt"
  "nexa/shared/types"
  "time"
)

const (
  UsageVerifStr = "verif"
  UsageResetStr = "reset"
)

type TokenUsage uint8

const (
  TokenUsageVerification TokenUsage = iota
  TokenUsageResetPassword
  TokenUsageUnknown
)

func NewTokenUsage(val uint8) (TokenUsage, error) {
  usage := TokenUsage(val)
  if !usage.Valid() {
    return usage, sharedErr.ErrEnumOutOfBounds
  }
  return usage, nil
}

func UsageFromString(s string) (TokenUsage, error) {
  switch s {
  case UsageVerifStr:
    return TokenUsageVerification, nil
  case UsageResetStr:
    return TokenUsageResetPassword, nil
  }
  return TokenUsageUnknown, sharedErr.ErrEnumOutOfBounds
}

func (u TokenUsage) Underlying() uint8 {
  return uint8(u)
}

func (u TokenUsage) Valid() bool {
  return u.Underlying() < TokenUsageUnknown.Underlying()
}

func (u TokenUsage) String() string {
  switch u {
  case TokenUsageVerification:
    return UsageVerifStr
  case TokenUsageResetPassword:
    return UsageResetStr
  }
  return "unknown"
}

func NewToken(userId types.Id, usage TokenUsage, expiryTime time.Duration) Token {
  return Token{
    Token:     sharedJwt.GenerateRefreshToken(),
    UserId:    userId,
    Usage:     usage,
    ExpiredAt: time.Now().UTC().Add(expiryTime),
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

type JWTToken struct {
  Id    types.Id
  Token string
}

type PairTokens struct {
  Access  JWTToken
  Refresh JWTToken
}
