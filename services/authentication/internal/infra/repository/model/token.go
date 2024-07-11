package model

import (
  "github.com/arcorium/nexa/shared/types"
  "github.com/arcorium/nexa/shared/util/repo"
  "github.com/arcorium/nexa/shared/variadic"
  "github.com/uptrace/bun"
  "nexa/services/authentication/internal/domain/entity"
  "time"
)

type TokenMapOption = repo.DataAccessModelMapOption[*entity.Token, *Token]

func FromTokenDomain(ent *entity.Token, opts ...TokenMapOption) Token {
  token := Token{
    Token:     ent.Token,
    UserId:    ent.UserId.String(),
    Usage:     ent.Usage.Underlying(),
    ExpiredAt: ent.ExpiredAt,
  }

  variadic.New(opts...).
    DoAll(repo.MapOptionFunc(ent, &token))

  return token
}

// Token the data is read only and shouldn't edit
type Token struct {
  bun.BaseModel `bun:"table:tokens"`

  Token     string    `bun:",nullzero,pk"`
  UserId    string    `bun:",nullzero,notnull,type:uuid,unique:creds_usage_idx"`
  Usage     uint8     `bun:",notnull,unique:creds_usage_idx"`
  ExpiredAt time.Time `bun:",nullzero,notnull"`

  //UpdatedAt time.Time `bun:",nullzero"`
  CreatedAt time.Time `bun:",nullzero,notnull"`
}

func (t *Token) ToDomain() (entity.Token, error) {
  userId, err := types.IdFromString(t.UserId)
  if err != nil {
    return entity.Token{}, err
  }

  usage, err := entity.NewTokenUsage(t.Usage)
  if err != nil {
    return entity.Token{}, err
  }

  return entity.Token{
    Token:     t.Token,
    UserId:    userId,
    Usage:     usage,
    ExpiredAt: t.ExpiredAt,
  }, nil
}
