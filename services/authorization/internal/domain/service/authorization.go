package service

import (
  "context"
  "nexa/services/authorization/internal/domain/dto"
  "nexa/shared/status"
)

type IAuthorization interface {
  IsAuthorized(ctx context.Context, authDto *dto.IsAuthorizationDTO) status.Object
}
