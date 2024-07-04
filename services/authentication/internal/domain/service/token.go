package service

import (
  "context"
  "nexa/services/authentication/internal/domain/dto"
  "nexa/shared/status"
  "nexa/shared/types"
)

type IToken interface {
  // Request to create a token
  Request(ctx context.Context, dto *dto.TokenCreateDTO) (dto.TokenResponseDTO, status.Object)
  // Verify token and return the user id related to the token
  Verify(ctx context.Context, dto *dto.TokenVerifyDTO) (types.Id, status.Object)
  //AddUsage(ctx context.Context, usageDTO *dto.TokenAddUsageDTO) (types.UserId, status.Object)
  //RemoveUsage(ctx context.Context, usageId types.UserId) status.Object
  //UpdateUsage(ctx context.Context, usageDTO *dto.TokenUpdateUsageDTO) status.Object
  //FindUsage(ctx context.Context, usageId types.UserId) (dto.TokenUsageResponseDTO, status.Object)
  //FindAllUsages(ctx context.Context, elementDTO *sharedDto.PagedElementDTO) (sharedDto.PagedElementResult[dto.TokenUsageResponseDTO], status.Object)
}
