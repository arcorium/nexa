package service

import (
  "context"
  sharedDto "github.com/arcorium/nexa/shared/dto"
  "github.com/arcorium/nexa/shared/status"
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/relation/internal/domain/dto"
)

type IBlock interface {
  Block(ctx context.Context, targetUserId types.Id) status.Object
  Unblock(ctx context.Context, targetUserId types.Id) status.Object
  // GetUsers return all blocked users based on userId
  GetUsers(ctx context.Context, pageDTO sharedDto.PagedElementDTO) (sharedDto.PagedElementResult[dto.BlockResponseDTO], status.Object)
  GetUsersCount(ctx context.Context, userIds types.Id) (dto.BlockCountResponseDTO, status.Object)
  ClearUsers(ctx context.Context, userId types.Id) status.Object
}