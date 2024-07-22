package service

import (
  "context"
  sharedDto "github.com/arcorium/nexa/shared/dto"
  "github.com/arcorium/nexa/shared/status"
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/relation/internal/domain/dto"
  "nexa/services/relation/internal/domain/entity"
)

type IFollow interface {
  Follow(ctx context.Context, targetUserId ...types.Id) status.Object
  Unfollow(ctx context.Context, targetUserId ...types.Id) status.Object
  GetFollowers(ctx context.Context, userId types.Id, pageDTO sharedDto.PagedElementDTO) (sharedDto.PagedElementResult[dto.FollowResponseDTO], status.Object)
  GetFollowings(ctx context.Context, userId types.Id, pageDTO sharedDto.PagedElementDTO) (sharedDto.PagedElementResult[dto.FollowResponseDTO], status.Object)
  GetStatus(ctx context.Context, userId types.Id, targetUserIds ...types.Id) ([]entity.FollowStatus, status.Object)
  GetUsersCount(ctx context.Context, userIds ...types.Id) ([]dto.FollowCountResponseDTO, status.Object)
  ClearUsers(ctx context.Context, userId types.Id) status.Object
}
