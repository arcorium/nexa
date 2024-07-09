package external

import (
  "context"
  authNv1 "github.com/arcorium/nexa/proto/gen/go/authentication/v1"
  "github.com/arcorium/nexa/shared/types"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  "nexa/services/user/internal/domain/dto"
  "nexa/services/user/internal/domain/external"
  "nexa/services/user/util"
)

func NewAuthenticationClient(conn grpc.ClientConnInterface) external.IAuthenticationClient {
  return &authenticationClient{
    credClient:  authNv1.NewCredentialServiceClient(conn),
    tokenClient: authNv1.NewTokenServiceClient(conn),
    tracer:      util.GetTracer(),
  }
}

type authenticationClient struct {
  credClient  authNv1.CredentialServiceClient
  tokenClient authNv1.TokenServiceClient

  tracer trace.Tracer
}

func (a *authenticationClient) DeleteCredentials(ctx context.Context, userId types.Id) error {
  ctx, span := a.tracer.Start(ctx, "AuthenticationClient.DeleteCredentials")
  defer span.End()

  req := authNv1.LogoutAllRequest{
    UserId: userId.String(),
  }
  _, err := a.credClient.LogoutAll(ctx, &req)
  return err
}

func (a *authenticationClient) VerifyToken(ctx context.Context, verificationDTO *dto.TokenVerificationDTO) (types.Id, error) {
  ctx, span := a.tracer.Start(ctx, "AuthenticationClient.VerifyToken")
  defer span.End()

  req := &authNv1.TokenVerifyRequest{
    Token: verificationDTO.Token,
    Usage: util.TokenPurposeToUsage(verificationDTO.Purpose),
  }

  res, err := a.tokenClient.Verify(ctx, req)
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

func (a *authenticationClient) GenerateToken(ctx context.Context, dtos *dto.TokenGenerationDTO) (dto.TokenResponseDTO, error) {
  ctx, span := a.tracer.Start(ctx, "AuthenticationClient.GenerateToken")
  defer span.End()

  req := authNv1.TokenCreateRequest{
    UserId: dtos.UserId.String(),
    Usage:  util.TokenPurposeToUsage(dtos.Purpose),
  }

  token, err := a.tokenClient.Create(ctx, &req)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.TokenResponseDTO{}, err
  }

  return dto.TokenResponseDTO{
    Token:     token.Token,
    Purpose:   util.TokenUsageToPurpose(token.Usage),
    ExpiredAt: token.ExpiredAt.AsTime(),
  }, nil
}
