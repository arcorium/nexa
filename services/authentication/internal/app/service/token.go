package service

import (
  "context"
  "nexa/services/authentication/config"
  "nexa/services/authentication/internal/domain/dto"
  "nexa/services/authentication/internal/domain/mapper"
  "nexa/services/authentication/internal/domain/repository"
  "nexa/services/authentication/internal/domain/service"
  "nexa/services/authentication/shared/errors"
  appUtil "nexa/services/authentication/util"
  sharedDto "nexa/shared/dto"
  "nexa/shared/status"
  "nexa/shared/types"
  "nexa/shared/util"
)

func NewToken(tokenRepo repository.IToken, usageRepo repository.ITokenUsage) service.IToken {
  return &tokenService{tokenRepo: tokenRepo, usageRepo: usageRepo}
}

type tokenService struct {
  tokenRepo repository.IToken
  usageRepo repository.ITokenUsage

  cfg *config.ServerConfig
}

func (t *tokenService) Request(ctx context.Context, dt *dto.TokenRequestDTO) (dto.TokenRequestResponseDTO, status.Object) {
  claims := appUtil.GetUserClaims(ctx)

  domain := dt.ToEntity(claims.UserId)
  domain.ExpiredAt = domain.ExpiredAt.Add(t.cfg.TokenExpiration)

  err := t.tokenRepo.Create(ctx, &domain)
  if err != nil {
    return dto.TokenRequestResponseDTO{}, status.FromRepository(err, status.NullCode)
  }

  return dto.TokenRequestResponseDTO{
    Token: domain.Token,
  }, status.Success()
}

func (t *tokenService) Verify(ctx context.Context, dto *dto.TokenVerifyDTO) status.Object {
  // Get the token
  token, err := t.tokenRepo.Find(ctx, dto.Token)
  if err != nil {
    return status.FromRepository(err, status.NullCode)
  }
  // Check the usage
  if !token.Usage.Id.Equal(dto.UsageId) {
    return status.ErrBadRequest(errors.ErrTokenDifferentUsage)
  }
  // Check expiration time
  if token.IsExpired() {
    return status.ErrBadRequest(errors.ErrTokenExpired)
  }
  // Remove from database
  err = t.tokenRepo.Delete(ctx, token.Token)
  if err != nil {
    return status.FromRepository(err, status.NullCode)
  }
  return status.Success()
}

func (t *tokenService) AddUsage(ctx context.Context, usageDTO *dto.TokenAddUsageDTO) (types.Id, status.Object) {
  entity := usageDTO.ToEntity()

  id, err := t.usageRepo.Create(ctx, &entity) // TODO: Remove id return from repository
  if err != nil {
    return types.Id{}, status.FromRepository(err, status.NullCode)
  }

  return id, status.Created()
}

func (t *tokenService) RemoveUsage(ctx context.Context, usageId types.Id) status.Object {
  err := t.usageRepo.Delete(ctx, usageId)
  if err != nil {
    return status.FromRepository(err, status.NullCode)
  }

  return status.Deleted()
}

func (t *tokenService) UpdateUsage(ctx context.Context, usageDTO *dto.TokenUpdateUsageDTO) status.Object {
  entity := usageDTO.ToEntity()

  err := t.usageRepo.Patch(ctx, &entity)
  if err != nil {
    return status.FromRepository(err, status.NullCode)
  }

  return status.Updated()
}

func (t *tokenService) FindUsage(ctx context.Context, usageId types.Id) (dto.TokenUsageResponseDTO, status.Object) {
  usage, err := t.usageRepo.Find(ctx, usageId)
  if err != nil {
    return dto.TokenUsageResponseDTO{}, status.FromRepository(err, status.NullCode)
  }

  return mapper.ToTokenUsageResponse(&usage), status.Success()
}

func (t *tokenService) FindAllUsages(ctx context.Context, elementDTO *sharedDto.PagedElementDTO) (sharedDto.PagedElementResult[dto.TokenUsageResponseDTO], status.Object) {
  result, err := t.usageRepo.FindAll(ctx, elementDTO.ToQueryParam())
  if err != nil {
    return sharedDto.PagedElementResult[dto.TokenUsageResponseDTO]{}, status.FromRepository(err, status.NullCode)
  }
  usages := util.CastSlice(result.Data, mapper.ToTokenUsageResponse)

  return sharedDto.NewPagedElementOutput2(usages, elementDTO, result.Total), status.Success()
}
