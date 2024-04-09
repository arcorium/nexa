package service

import (
	"context"
	"nexa/services/authentication/internal/domain/dto"
	sharedDto "nexa/shared/dto"
	"nexa/shared/status"
	"nexa/shared/types"
)

type IToken interface {
	Request(ctx context.Context, dto *dto.TokenRequestDTO) (dto.TokenRequestResponseDTO, status.Object)
	Verify(ctx context.Context, dto *dto.TokenVerifyDTO) status.Object
	AddUsage(ctx context.Context, usageDTO *dto.TokenAddUsageDTO) (types.Id, status.Object)
	RemoveUsage(ctx context.Context, usageId types.Id) status.Object
	UpdateUsage(ctx context.Context, usageDTO *dto.TokenUpdateUsageDTO) status.Object
	FindUsage(ctx context.Context, usageId types.Id) (dto.TokenUsageResponseDTO, status.Object)
	FindAllUsages(ctx context.Context, elementDTO *sharedDto.PagedElementDTO) (sharedDto.PagedElementResult[dto.TokenUsageResponseDTO], status.Object)
}
