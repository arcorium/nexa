package external

import (
  "context"
  tokenv1 "github.com/arcorium/nexa/proto/gen/go/token/v1"
  "github.com/arcorium/nexa/shared/types"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  "nexa/services/authentication/internal/domain/dto"
  "nexa/services/authentication/internal/domain/external"
  "nexa/services/authentication/util"
)

func NewTokenClient(conn grpc.ClientConnInterface) external.ITokenClient {
  return &tokenClient{
    client: tokenv1.NewTokenServiceClient(conn),
    tracer: util.GetTracer(),
  }
}

type tokenClient struct {
  client tokenv1.TokenServiceClient

  tracer trace.Tracer
}

func (a *tokenClient) Verify(ctx context.Context, verificationDTO *dto.TokenVerificationDTO) (types.Id, error) {
  ctx, span := a.tracer.Start(ctx, "TokenClient.Verify")
  defer span.End()

  req := &tokenv1.VerifyTokenRequest{
    Token: verificationDTO.Token,
    Usage: util.TokenPurposeToUsage(verificationDTO.Purpose),
  }

  res, err := a.client.Verify(ctx, req)
  if err != nil {
    spanUtil.RecordError(err, span)
    return types.NullId(), err
  }

  userId, err := types.IdFromString(res.UserId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return userId, err
  }

  return userId, nil
}

func (a *tokenClient) Generate(ctx context.Context, dtos *dto.TokenGenerationDTO) (dto.TokenResponseDTO, error) {
  ctx, span := a.tracer.Start(ctx, "TokenClient.Generate")
  defer span.End()

  req := tokenv1.CreateTokenRequest{
    UserId: dtos.UserId.String(),
    Usage:  util.TokenPurposeToUsage(dtos.Usage),
  }

  token, err := a.client.Create(ctx, &req)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.TokenResponseDTO{}, err
  }

  userId, err := types.IdFromString(token.UserId)
  if err != nil {
    return dto.TokenResponseDTO{}, err
  }

  return dto.TokenResponseDTO{
    Token:     token.Token,
    UserId:    userId,
    Usage:     util.TokenUsageToPurpose(token.Usage),
    ExpiredAt: token.ExpiredAt.AsTime(),
  }, nil
}
