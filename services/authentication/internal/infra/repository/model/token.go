package model

import (
  "github.com/uptrace/bun"
  "nexa/services/authentication/internal/domain/entity"
  "nexa/shared/types"
  "nexa/shared/util/repo"
  "nexa/shared/variadic"
  "time"
)

type TokenMapOption = repo.DataAccessModelMapOption[*entity.Token, *Token]

func FromTokenDomain(domain *entity.Token, opts ...TokenMapOption) Token {
  token := Token{
    Token:   domain.Token,
    UserId:  domain.UserId.Underlying().String(),
    UsageId: domain.Usage.Id.Underlying().String(),
  }

  variadic.New(opts...).
    DoAll(repo.MapOptionFunc(domain, &token))

  return token
}

type Token struct {
  bun.BaseModel `bun:"table:tokens"`

  Token     string    `bun:",nullzero,pk"`
  UserId    string    `bun:",nullzero,notnull,type:uuid,unique:creds_usage_idx"`
  UsageId   string    `bun:",nullzero,notnull,type:uuid,unique:creds_usage_idx"`
  ExpiredAt time.Time `bun:",nullzero,notnull"`

  UpdatedAt time.Time `bun:",nullzero"`
  CreatedAt time.Time `bun:",nullzero,notnull"`

  TokenUsage *TokenUsage `bun:"rel:belongs-to,join:usage_id=id,on_delete:CASCADE"`
}

func (t *Token) ToDomain() entity.Token {
  return entity.Token{
    Token:     t.Token,
    UserId:    types.IdFromString(t.UserId),
    Usage:     t.TokenUsage.ToDomain(),
    ExpiredAt: t.ExpiredAt,
  }
}
