package model

import (
	"github.com/uptrace/bun"
	"nexa/services/authentication/shared/domain/entity"
	"nexa/shared/types"
	"nexa/shared/util/repo"
	"nexa/shared/variadic"
	"time"
)

type TokenUsageMapOption = repo.DataAccessModelMapOption[*entity.TokenUsage, *TokenUsage]

func FromTokenUsageDomain(domain *entity.TokenUsage, opts ...TokenUsageMapOption) TokenUsage {
	obj := TokenUsage{
		Id:          domain.Id.Underlying().String(),
		Name:        domain.Name,
		Description: domain.Description,
	}

	variadic.New(opts...).
		DoAll(repo.MapOptionFunc(domain, &obj))

	return obj
}

type TokenUsage struct {
	bun.BaseModel `bun:"table:verification_usages"`

	Id          string `bun:",nullzero,type:uuid,pk"`
	Name        string `bun:",nullzero,notnull,unique"`
	Description string `bun:",nullzero"`

	UpdatedAt time.Time `bun:",nullzero"`
	CreatedAt time.Time `bun:",nullzero,notnull"`
}

func (o *TokenUsage) ToDomain() entity.TokenUsage {
	return entity.TokenUsage{
		Id:          types.IdFromString(o.Id),
		Name:        o.Name,
		Description: o.Description,
	}
}
