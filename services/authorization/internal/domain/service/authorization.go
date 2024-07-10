package service

import (
  "context"
  "github.com/arcorium/nexa/shared/status"
  "nexa/services/authorization/internal/domain/dto"
)

type IAuthorization interface {
  IsAuthorized(ctx context.Context, authDto *dto.IsAuthorizationDTO) status.Object
}
