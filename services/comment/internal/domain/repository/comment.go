package repository

import (
  "context"
  "github.com/arcorium/nexa/shared/types"
  "github.com/arcorium/nexa/shared/util/repo"
  "nexa/services/comment/internal/domain/entity"
)

type IComment interface {
  GetReplies(ctx context.Context, showReply bool, commentId types.Id, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Comment], error)
  GetReplyCounts(ctx context.Context, commentIds ...types.Id) ([]entity.Count, error)
  GetPostCounts(ctx context.Context, postIds ...types.Id) ([]entity.Count, error)
  FindByPostId(ctx context.Context, showReply bool, postId types.Id, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Comment], error)
  Create(ctx context.Context, comment *entity.Comment) error
  UpdateContent(ctx context.Context, commentId types.Id, content string) error
  DeleteByIds(ctx context.Context, commentIds ...types.Id) error
  DeleteUsers(ctx context.Context, userId types.Id, commentIds ...types.Id) ([]types.Id, error)
  DeleteByPostIds(ctx context.Context, postIds ...types.Id) ([]types.Id, error)
}
