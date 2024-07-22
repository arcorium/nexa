package service

import (
  "context"
  sharedDto "github.com/arcorium/nexa/shared/dto"
  "github.com/arcorium/nexa/shared/status"
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/post/internal/domain/dto"
  "nexa/services/post/internal/domain/entity"
)

type IPost interface {
  GetAll(ctx context.Context, pageDTO *sharedDto.PagedElementDTO) (sharedDto.PagedElementResult[dto.PostResponseDTO], status.Object)
  GetEdited(ctx context.Context, postId types.Id) (dto.EditedPostResponseDTO, status.Object)
  FindById(ctx context.Context, id types.Id) (dto.PostResponseDTO, status.Object)
  FindByUserId(ctx context.Context, userId types.Id, pageDTO *sharedDto.PagedElementDTO) (sharedDto.PagedElementResult[dto.PostResponseDTO], status.Object)
  Create(ctx context.Context, createDTO *dto.CreatePostDTO) (types.Id, status.Object)
  UpdateVisibility(ctx context.Context, id types.Id, newVisibility entity.Visibility) status.Object
  Edit(ctx context.Context, editDTO *dto.EditPostDTO) status.Object
  Delete(ctx context.Context, id types.Id) status.Object
  ToggleBookmark(ctx context.Context, postId types.Id) status.Object
  GetBookmarked(ctx context.Context, userId types.Id, elementDTO *sharedDto.PagedElementDTO) (sharedDto.PagedElementResult[dto.PostResponseDTO], status.Object)
  ClearUsers(ctx context.Context, userId types.Id) status.Object
}
