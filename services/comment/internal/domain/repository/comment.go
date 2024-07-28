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
  FindById(ctx context.Context, showReply bool, commentId types.Id) ([]entity.Comment, error)
  FindByPostId(ctx context.Context, showReply bool, postId types.Id, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Comment], error)
  Count(ctx context.Context, commentIds ...types.Id) (uint64, error)
  Create(ctx context.Context, comment *entity.Comment) error
  UpdateContent(ctx context.Context, userId types.Id, commentId types.Id, content string) error
  DeleteByIds(ctx context.Context, commentIds ...types.Id) error
  DeleteUsers(ctx context.Context, userId types.Id, commentIds ...types.Id) ([]types.Id, error)
  DeleteByPostIds(ctx context.Context, postIds ...types.Id) ([]types.Id, error)
}
