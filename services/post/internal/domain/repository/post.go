package repository

import (
  "context"
  "github.com/arcorium/nexa/shared/types"
  "github.com/arcorium/nexa/shared/util/repo"
  "nexa/services/post/internal/domain/dto"
  "nexa/services/post/internal/domain/entity"
)

type IPost interface {
  Get(ctx context.Context, expectedVisibility entity.Visibility, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Post], error)
  GetBookmarked(ctx context.Context, userId types.Id, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Post], error)
  FindById(ctx context.Context, ids ...types.Id) ([]entity.Post, error)
  GetEdited(ctx context.Context, postId types.Id) (entity.Post, error)
  FindByUserId(ctx context.Context, userId types.Id, expectedVisibility entity.Visibility, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Post], error)
  // DelsertBookmark is chaining operation that will delete the relation data when the data is duplicated or already there.
  // Otherwise, it will normally insert the relation data
  DelsertBookmark(ctx context.Context, userId, postId types.Id) error
  Create(ctx context.Context, post *entity.Post) error
  UpdateVisibility(ctx context.Context, userId types.Id, postId types.Id, visibility entity.Visibility) error
  Edit(ctx context.Context, post *entity.Post, flag dto.EditPostFlag) error
  // DeleteUsers works like Remove, but it needs the user id so the post should be belonged to it. if the postIds is nil it will remove all user's posts
  // When the postIds is nil (trying to delete all user posts), it will return all the ids
  DeleteUsers(ctx context.Context, userId types.Id, postIds ...types.Id) ([]types.Id, error)
  Delete(ctx context.Context, postIds ...types.Id) error
}
