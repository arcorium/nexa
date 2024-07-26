package repository

import (
  "context"
  "github.com/arcorium/nexa/shared/types"
  "github.com/arcorium/nexa/shared/util/repo"
  "nexa/services/relation/internal/domain/entity"
)

type IFollow interface {
  Create(ctx context.Context, follow *entity.Follow) error
  Creates(ctx context.Context, follows []entity.Follow) error
  Delete(ctx context.Context, follows ...entity.Follow) error
  GetFollowers(ctx context.Context, userId types.Id, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Follow], error)
  GetFollowings(ctx context.Context, userId types.Id, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Follow], error)
  IsFollowing(ctx context.Context, userId types.Id, followeeIds ...types.Id) ([]bool, error)
  GetCounts(ctx context.Context, userIds ...types.Id) ([]entity.FollowCount, error)
  DeleteByUserId(ctx context.Context, deleteFollower bool, userId types.Id) error
}
