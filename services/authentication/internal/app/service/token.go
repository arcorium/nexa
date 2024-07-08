package service

import (
  "context"
  "database/sql"
  "go.opentelemetry.io/otel/trace"
  "nexa/services/authentication/internal/domain/dto"
  "nexa/services/authentication/internal/domain/entity"
  "nexa/services/authentication/internal/domain/mapper"
  "nexa/services/authentication/internal/domain/repository"
  "nexa/services/authentication/internal/domain/service"
  "nexa/services/authentication/util"
  "nexa/services/authentication/util/errors"
  "nexa/shared/status"
  "nexa/shared/types"
  spanUtil "nexa/shared/util/span"
  "time"
)

func NewToken(tokenRepo repository.IToken, conf TokenServiceConfig) service.IToken {
  return &tokenService{
    config:    conf,
    tokenRepo: tokenRepo,
    tracer:    util.GetTracer(),
  }
}

type TokenServiceConfig struct {
  VerificationTokenExpiration time.Duration
  ResetTokenExpiration        time.Duration
}

type tokenService struct {
  config    TokenServiceConfig
  tokenRepo repository.IToken
  tracer    trace.Tracer
}

func (t *tokenService) Request(ctx context.Context, req *dto.TokenCreateDTO) (dto.TokenResponseDTO, status.Object) {
  ctx, span := t.tracer.Start(ctx, "TokenService.Request")
  defer span.End()

  var expiryDuration time.Duration
  if req.Usage == entity.TokenUsageVerification {
    expiryDuration = t.config.VerificationTokenExpiration
  } else if req.Usage == entity.TokenUsageResetPassword {
    expiryDuration = t.config.ResetTokenExpiration
  }
  // NOTE: Doesn't need else, its already handled by handler

  domain := req.ToDomain(expiryDuration)

  err := t.tokenRepo.Create(ctx, &domain)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.TokenResponseDTO{}, status.FromRepository(err, status.NullCode)
  }

  responseDTO := mapper.ToTokenResponseDTO(&domain)
  return responseDTO, status.Created()
}

func (t *tokenService) Verify(ctx context.Context, verifyDTO *dto.TokenVerifyDTO) (types.Id, status.Object) {
  ctx, span := t.tracer.Start(ctx, "TokenService.Verify")
  defer span.End()

  // Get the token
  token, err := t.tokenRepo.Find(ctx, verifyDTO.Token)
  if err != nil {
    spanUtil.RecordError(err, span)
    if err == sql.ErrNoRows {
      return types.NullId(), status.ErrBadRequest(errors.ErrTokenNotFound)
    }
    return types.NullId(), status.FromRepository(err, status.NullCode)
  }

  // Check the usage
  if token.Usage != verifyDTO.Usage {
    spanUtil.RecordError(errors.ErrTokenDifferentUsage, span)
    return types.NullId(), status.ErrBadRequest(errors.ErrTokenDifferentUsage)
  }

  // Remove from database
  err = t.tokenRepo.Delete(ctx, token.Token)
  if err != nil {
    spanUtil.RecordError(err, span)
    return types.NullId(), status.FromRepository(err, status.NullCode)
  }

  // Check expiration time
  if token.IsExpired() {
    spanUtil.RecordError(errors.ErrTokenExpired, span)
    return types.NullId(), status.ErrBadRequest(errors.ErrTokenExpired)
  }

  return token.UserId, status.Success()
}
