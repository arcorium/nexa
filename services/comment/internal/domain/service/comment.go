package service

import (
  "context"
  sharedDto "github.com/arcorium/nexa/shared/dto"
  "github.com/arcorium/nexa/shared/status"
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/comment/internal/domain/dto"
  "nexa/services/comment/internal/domain/entity"
)

type IComment interface {
  Create(ctx context.Context, commentDTO *dto.CreateCommentDTO) (types.Id, status.Object)
  Edit(ctx context.Context, commentDTO *dto.EditCommentDTO) status.Object
  Delete(ctx context.Context, commentIds ...types.Id) status.Object
  FindById(ctx context.Context, findDTO *dto.FindCommentByIdDTO) (dto.CommentResponseDTO, status.Object)
  GetPosts(ctx context.Context, dto *dto.GetPostsCommentsDTO, elementDTO sharedDto.PagedElementDTO) (sharedDto.PagedElementResult[dto.CommentResponseDTO], status.Object)
  GetReplies(ctx context.Context, repliesDTO *dto.GetCommentsRepliesDTO, elementDTO sharedDto.PagedElementDTO) (sharedDto.PagedElementResult[dto.CommentResponseDTO], status.Object)
  GetCounts(ctx context.Context, itemType entity.ItemType, itemIds ...types.Id) ([]uint64, status.Object)
  IsExists(ctx context.Context, commentIds ...types.Id) (bool, status.Object)
  ClearPosts(ctx context.Context, postIds ...types.Id) status.Object
  ClearUsers(ctx context.Context, userId types.Id) status.Object
}
