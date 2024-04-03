package model

import (
	"github.com/uptrace/bun"
	"nexa/services/authorization/shared/domain/entity"
	"nexa/shared/types"
	"nexa/shared/util/repo"
	"nexa/shared/variadic"
	"time"
)

type ActionMapOption = repo.DataAccessModelMapOption[*entity.Action, *Action]

func FromActionDomain(domain *entity.Action, opts ...ActionMapOption) Action {
	action := Action{
		BaseModel:   bun.BaseModel{},
		Id:          domain.Id.Underlying().String(),
		Name:        domain.Name,
		Description: domain.Description,
	}

	variadic.New(opts...).DoAtFirst(func(a *ActionMapOption) {
		(*a)(domain, &action)
	})

	return action
}

type Action struct {
	bun.BaseModel `bun:"table:actions"`

	Id          string `bun:",nullzero,type:uuid,pk"`
	Name        string `bun:",nullzero,unique"`
	Description string `bun:",nullzero"`

	UpdatedAt time.Time `bun:",nullzero"`
	CreatedAt time.Time `bun:",nullzero,notnull"`
}

func (a *Action) ToDomain() entity.Action {
	return entity.Action{
		Id:          types.IdFromString(a.Id),
		Name:        a.Name,
		Description: a.Description,
	}
}
