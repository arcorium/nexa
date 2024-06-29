package model

import (
  "database/sql"
  "github.com/uptrace/bun"
  "nexa/services/authentication/internal/domain/entity"
  "nexa/shared/types"
  "nexa/shared/util/repo"
  "nexa/shared/variadic"
  "nexa/shared/wrapper"
  "time"
)

type TokenMapOption = repo.DataAccessModelMapOption[*entity.Token, *Token]

func FromTokenDomain(domain *entity.Token, opts ...TokenMapOption) Token {
  token := Token{
    Token:     domain.Token,
    UserId:    domain.UserId.Underlying().String(),
    Usage:     sql.NullInt64{Int64: int64(domain.Usage.Underlying()), Valid: true},
    ExpiredAt: domain.ExpiredAt,
  }

  variadic.New(opts...).
    DoAll(repo.MapOptionFunc(domain, &token))

  return token
}

type Token struct {
  bun.BaseModel `bun:"table:tokens"`

  Token     string        `bun:",nullzero,pk"`
  UserId    string        `bun:",nullzero,notnull,type:uuid,unique:creds_usage_idx"`
  Usage     sql.NullInt64 `bun:",type:smallint,notnull,unique:creds_usage_idx"`
  ExpiredAt time.Time     `bun:",nullzero,notnull"`

  UpdatedAt time.Time `bun:",nullzero"`
  CreatedAt time.Time `bun:",nullzero,notnull"`
}

func (t *Token) ToDomain() entity.Token {
  return entity.Token{
    Token:     t.Token,
    UserId:    wrapper.DropError(types.IdFromString(t.UserId)),
    Usage:     entity.TokenUsage(t.Usage.Int64),
    ExpiredAt: t.ExpiredAt,
  }
}
