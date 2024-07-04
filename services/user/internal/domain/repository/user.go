package repository

import (
  "context"
  "nexa/services/user/internal/domain/entity"
  "nexa/shared/types"
  "nexa/shared/util/repo"
)

type IUser interface {
  // Create make new user
  Create(ctx context.Context, user *entity.User) error
  // Update update all fields of user based on the id
  Update(ctx context.Context, user *entity.User) error
  // Patch update all non-zero fields of user based on the id
  Patch(ctx context.Context, user *entity.User) error
  // FindByIds get all users based on specified user ids
  FindByIds(ctx context.Context, userIds ...types.Id) ([]entity.User, error)
  // FindByEmails get all users based on specified emails
  FindByEmails(ctx context.Context, emails ...types.Email) ([]entity.User, error)
  // FindAllUsers get all registered users
  Get(ctx context.Context, query repo.QueryParameter) (repo.PaginatedResult[entity.User], error)
  // Delete delete user based on the id
  Delete(ctx context.Context, ids ...types.Id) error
}
