package factory

import (
	"github.com/uptrace/bun"
	"nexa/services/authorization/internal/domain/repository"
	"nexa/services/authorization/internal/infra/repository/pg"
)

func NewPostgresRepositories(db bun.IDB) PGRepositoryFactory {
	return PGRepositoryFactory{
		Action:     pg.NewAction(db),
		Permission: pg.NewPermission(db),
		Resource:   pg.NewResource(db),
		Role:       pg.NewRole(db),
	}
}

type PGRepositoryFactory struct {
	Action     repository.IAction
	Permission repository.IPermission
	Resource   repository.IResource
	Role       repository.IRole
}
