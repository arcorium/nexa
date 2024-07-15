package service

import (
  "context"
  "github.com/arcorium/nexa/shared/status"
  "github.com/arcorium/nexa/shared/types"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "go.opentelemetry.io/otel/trace"
  "nexa/services/token/internal/domain/dto"
  "nexa/services/token/internal/domain/entity"
  "nexa/services/token/internal/domain/mapper"
  "nexa/services/token/internal/domain/repository"
  "nexa/services/token/internal/domain/service"
  "nexa/services/token/util"
  "nexa/services/token/util/errors"
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
  LoginTokenExpiration        time.Duration
  GeneralTokenExpiration      time.Duration
}

type tokenService struct {
  config    TokenServiceConfig
  tokenRepo repository.IToken
  tracer    trace.Tracer
}

func (t *tokenService) getTokenUsageDuration(usage entity.TokenUsage) (time.Duration, error) {
  switch usage {
  case entity.TokenUsageEmailVerification:
    return t.config.VerificationTokenExpiration, nil
  case entity.TokenUsageResetPassword:
    return t.config.ResetTokenExpiration, nil
  case entity.TokenUsageLogin:
    return t.config.LoginTokenExpiration, nil
  case entity.TokenUsageGeneral:
    return t.config.GeneralTokenExpiration, nil
  }
  return 0, errors.ErrTokenUsageUnknown
}

func (t *tokenService) checkToken(token *entity.Token, verifyDTO *dto.TokenVerifyDTO) error {
  // Check the usage
  if token.Usage != verifyDTO.ExpectedUsage {
    return errors.ErrTokenDifferentUsage
  }

  // Check expiration time
  if token.IsExpired() {
    return errors.ErrTokenExpired
  }
  return nil
}

func (t *tokenService) Request(ctx context.Context, req *dto.TokenCreateDTO) (dto.TokenResponseDTO, status.Object) {
  ctx, span := t.tracer.Start(ctx, "TokenService.Request")
  defer span.End()

  expiryDuration, err := t.getTokenUsageDuration(req.Usage)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.TokenResponseDTO{}, status.ErrBadRequest(err)
  }

  domain := req.ToDomain(expiryDuration)
  err = t.tokenRepo.Upsert(ctx, &domain)
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
    return types.NullId(), status.FromRepositoryOverride(err, types.NewPair(status.BAD_REQUEST_ERROR, errors.ErrTokenNotFound))
  }

  // Remove from database
  err = t.tokenRepo.Delete(ctx, token.Token)
  if err != nil {
    spanUtil.RecordError(err, span)
    return types.NullId(), status.FromRepository(err, status.NullCode)
  }

  if err = t.checkToken(&token, verifyDTO); err != nil {
    spanUtil.RecordError(err, span)
    return types.NullId(), status.ErrBadRequest(err)
  }

  return token.UserId, status.Success()
}

func (t *tokenService) AuthVerify(ctx context.Context, verifyDTO *dto.TokenAuthVerifyDTO) status.Object {
  ctx, span := t.tracer.Start(ctx, "TokenService.AuthVerify")
  defer span.End()

  // Get the token
  token, err := t.tokenRepo.Find(ctx, verifyDTO.Token)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepositoryOverride(err, types.NewPair(status.BAD_REQUEST_ERROR, errors.ErrTokenNotFound))
  }

  // Check user
  if !token.UserId.Eq(verifyDTO.ExpectedUserId) {
    spanUtil.RecordError(errors.ErrTokenDifferentUsage, span)
    return status.ErrUnAuthorized(errors.ErrTokenDifferentUser)
  }

  // Remove from database
  err = t.tokenRepo.Delete(ctx, token.Token)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrInternal(err)
  }

  if err = t.checkToken(&token, &verifyDTO.TokenVerifyDTO); err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrBadRequest(err)
  }

  return status.Success()
}
