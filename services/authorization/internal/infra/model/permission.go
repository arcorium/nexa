package model

import (
	"github.com/uptrace/bun"
	"nexa/services/authorization/shared/domain/entity"
	"nexa/shared/types"
	"nexa/shared/util/repo"
	"nexa/shared/variadic"
	"time"
)

type PermissionMapOption = repo.DataAccessModelMapOption[*entity.Permission, *Permission]

func FromPermissionDomain(domain *entity.Permission, opts ...PermissionMapOption) Permission {
	permission := Permission{
		Id:         domain.Id.Underlying().String(),
		ResourceId: domain.Resource.Id.Underlying().String(),
		ActionId:   domain.Action.Id.Underlying().String(),
	}

	variadic.New(opts...).DoAtFirst(func(p *PermissionMapOption) {
		(*p)(domain, &permission)
	})

	return permission
}

type Permission struct {
	bun.BaseModel `bun:"permissions"`

	Id         string `bun:",nullzero,type:uuid,pk"`
	ResourceId string `bun:",nullzero,type:uuid,unique:resource_action_idx"`
	ActionId   string `bun:",nullzero,type:uuid,unique:resource_action_idx"`

	CreatedAt time.Time `bun:",nullzero,notnull"`

	Resource *Resource `bun:"rel:belongs-to,join=resource_id=id,on_delete:CASCADE"`
	Action   *Action   `bun:"rel:belongs-to,join=action_id=id,on_delete:CASCADE"`
}

func (p *Permission) ToDomain() entity.Permission {
	return entity.Permission{
		Id: types.IdFromString(p.Id),
		Resource: entity.Resource{
			Id:          types.IdFromString(p.ResourceId),
			Name:        p.Resource.Name,
			Description: p.Resource.Description,
		},
		Action: entity.Action{
			Id:          types.IdFromString(p.ActionId),
			Name:        p.Action.Name,
			Description: p.Action.Description,
		},
	}
}
