package repository

import (
	"context"
	"nexa/services/user/shared/domain/entity"
	"nexa/shared/types"
)

type IProfile interface {
	// Create make new profile based on user id
	Create(ctx context.Context, profile *entity.Profile) error
	// FindByIds find all profiles based on user ids
	FindByIds(ctx context.Context, userIds ...types.Id) ([]entity.Profile, error)
	// Update update all fields of profile based on user id
	Update(ctx context.Context, profile *entity.Profile) error
	// Patch update all non-zero fields of profile based on user id
	Patch(ctx context.Context, profile *entity.Profile) error
}
