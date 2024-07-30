package service

import (
  "context"
  "github.com/arcorium/nexa/shared/status"
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/feed/internal/domain/dto"
)

type IFeed interface {
  GetUserFeed(ctx context.Context, userId types.Id, limit uint64) ([]dto.PostResponseDTO, status.Object)
}
