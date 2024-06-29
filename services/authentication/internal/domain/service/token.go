package service

import (
  "context"
  "nexa/services/authentication/internal/domain/dto"
  "nexa/shared/status"
)

type IToken interface {
  Request(ctx context.Context, dto *dto.TokenCreateDTO) (dto.TokenResponseDTO, status.Object)
  Verify(ctx context.Context, dto *dto.TokenVerifyDTO) status.Object
  //AddUsage(ctx context.Context, usageDTO *dto.TokenAddUsageDTO) (types.UserId, status.Object)
  //RemoveUsage(ctx context.Context, usageId types.UserId) status.Object
  //UpdateUsage(ctx context.Context, usageDTO *dto.TokenUpdateUsageDTO) status.Object
  //FindUsage(ctx context.Context, usageId types.UserId) (dto.TokenUsageResponseDTO, status.Object)
  //FindAllUsages(ctx context.Context, elementDTO *sharedDto.PagedElementDTO) (sharedDto.PagedElementResult[dto.TokenUsageResponseDTO], status.Object)
}
