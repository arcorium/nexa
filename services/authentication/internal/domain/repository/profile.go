package repository

import (
  "context"
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/authentication/internal/domain/entity"
)

type IProfile interface {
  // Create make new profile based on user id
  Create(ctx context.Context, profile *entity.Profile) error
  // FindByIds find all profiles based on the id
  FindByIds(ctx context.Context, userIds ...types.Id) ([]entity.Profile, error)
  // FindByUserId find profiles based on the user id
  FindByUserId(ctx context.Context, userId types.Id) (*entity.Profile, error)
  // Update update all fields of profile based on user id
  Update(ctx context.Context, profile *entity.Profile) error
  // Patch update all non-zero fields of profile based on user id
  Patch(ctx context.Context, profile *entity.PatchedProfile) error
  //// Delete delete profiles based on user ids
  //Delete(ctx context.Context, userIds ...types.Id) error
}
