package external

import (
  "context"
  tokenv1 "github.com/arcorium/nexa/proto/gen/go/token/v1"
  "github.com/arcorium/nexa/shared/types"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "github.com/sony/gobreaker"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  "nexa/services/authentication/config"
  "nexa/services/authentication/internal/domain/dto"
  "nexa/services/authentication/internal/domain/external"
  "nexa/services/authentication/util"
)

func NewTokenClient(conn grpc.ClientConnInterface, conf *config.CircuitBreaker) external.ITokenClient {
  breaker := gobreaker.NewCircuitBreaker(gobreaker.Settings{
    Name:         "nexa-token",
    MaxRequests:  conf.MaxRequest,
    Interval:     conf.ResetInterval,
    Timeout:      conf.OpenStateTimeout,
    IsSuccessful: util.IsGrpcConnectivityError,
  })

  return &tokenClient{
    client: tokenv1.NewTokenServiceClient(conn),
    tracer: util.GetTracer(),
    cb:     breaker,
  }
}

type tokenClient struct {
  client tokenv1.TokenServiceClient
  tracer trace.Tracer
  cb     *gobreaker.CircuitBreaker
}

func (a *tokenClient) Verify(ctx context.Context, verificationDTO *dto.TokenVerificationDTO) (types.Id, error) {
  ctx, span := a.tracer.Start(ctx, "TokenClient.Verify")
  defer span.End()

  req := &tokenv1.VerifyTokenRequest{
    Token: verificationDTO.Token,
    Usage: util.TokenPurposeToUsage(verificationDTO.Purpose),
  }

  result, err := a.cb.Execute(func() (interface{}, error) {
    return a.client.Verify(ctx, req)
  })
  if err != nil {
    spanUtil.RecordError(err, span)
    return types.NullId(), err
  }

  resp := result.(*tokenv1.VerifyTokenResponse)
  userId, err := types.IdFromString(resp.UserId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return userId, err
  }

  return userId, nil
}

func (a *tokenClient) Generate(ctx context.Context, genDTO *dto.TokenGenerationDTO) (dto.TokenResponseDTO, error) {
  ctx, span := a.tracer.Start(ctx, "TokenClient.Generate")
  defer span.End()

  req := tokenv1.CreateTokenRequest{
    UserId: genDTO.UserId.String(),
    Usage:  util.TokenPurposeToUsage(genDTO.Usage),
  }

  result, err := a.cb.Execute(func() (interface{}, error) {
    return a.client.Create(ctx, &req)
  })
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.TokenResponseDTO{}, err
  }

  resp := result.(*tokenv1.Token)
  userId, err := types.IdFromString(resp.UserId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.TokenResponseDTO{}, err
  }

  return dto.TokenResponseDTO{
    Token:     resp.Token,
    UserId:    userId,
    Usage:     util.TokenUsageToPurpose(resp.Usage),
    ExpiredAt: resp.ExpiredAt.AsTime(),
  }, nil
}
