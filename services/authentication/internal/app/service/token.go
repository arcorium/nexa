package service

import (
  "context"
  "go.opentelemetry.io/otel/trace"
  "nexa/services/authentication/internal/domain/dto"
  "nexa/services/authentication/internal/domain/entity"
  "nexa/services/authentication/internal/domain/mapper"
  "nexa/services/authentication/internal/domain/repository"
  "nexa/services/authentication/internal/domain/service"
  "nexa/services/authentication/util"
  "nexa/services/authentication/util/errors"
  sharedErr "nexa/shared/errors"
  spanUtil "nexa/shared/span"
  "nexa/shared/status"
  sharedUtil "nexa/shared/util"
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
  type Null = dto.TokenResponseDTO

  ctx, span := t.tracer.Start(ctx, "TokenService.Request")
  defer span.End()

  // Input validation
  var expiryDuration time.Duration
  if req.Usage == entity.TokenUsageVerification.Underlying() {
    expiryDuration = t.config.VerificationTokenExpiration
  } else if req.Usage == entity.TokenUsageResetPassword.Underlying() {
    expiryDuration = t.config.ResetTokenExpiration
  } else {
    err := sharedErr.ErrEnumOutOfBounds
    spanUtil.RecordError(err, span)
    return Null{}, status.ErrBadRequest(err)
  }

  domain, err := req.ToDomain(expiryDuration)
  if err != nil {
    spanUtil.RecordError(err, span)
    return Null{}, status.ErrBadRequest(err)
  }

  err = t.tokenRepo.Create(ctx, &domain)
  if err != nil {
    spanUtil.RecordError(err, span)
    return Null{}, status.FromRepository(err, status.NullCode)
  }

  responseDTO := mapper.ToTokenResponseDTO(&domain)
  return responseDTO, status.Created()
}

func (t *tokenService) Verify(ctx context.Context, dto *dto.TokenVerifyDTO) status.Object {
  ctx, span := t.tracer.Start(ctx, "TokenService.Verify")
  defer span.End()

  // Input validation
  if err := sharedUtil.ValidateStructCtx(ctx, &dto); err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrBadRequest(err)
  }

  // Get the token
  token, err := t.tokenRepo.Find(ctx, dto.Token)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }

  // Check the usage
  if token.Usage.Underlying() == dto.Usage {
    spanUtil.RecordError(errors.ErrTokenDifferentUsage, span)
    return status.ErrBadRequest(errors.ErrTokenDifferentUsage)
  }

  // Remove from database
  err = t.tokenRepo.Delete(ctx, token.Token)
  if err != nil {
    //return status.FromRepository(err, status.NullCode)
    spanUtil.RecordError(err, span)
    return status.ErrInternal(err)
  }

  // Check expiration time
  if token.IsExpired() {
    spanUtil.RecordError(errors.ErrTokenExpired, span)
    return status.ErrBadRequest(errors.ErrTokenExpired)
  }

  return status.Success()
}
