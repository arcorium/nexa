package handler

import (
  "context"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  "google.golang.org/protobuf/types/known/emptypb"
  authNv1 "nexa/proto/gen/go/authentication/v1"
  "nexa/services/authentication/internal/api/grpc/mapper"
  "nexa/services/authentication/internal/domain/service"
  spanUtil "nexa/shared/span"
  sharedUtil "nexa/shared/util"
)

func NewCredential(svc service.ICredential) CredentialHandler {
  return CredentialHandler{
    credService: svc,
  }
}

type CredentialHandler struct {
  authNv1.UnimplementedCredentialServiceServer
  credService service.ICredential
}

func (c *CredentialHandler) RegisterHandler(server *grpc.Server) {
  authNv1.RegisterCredentialServiceServer(server, c)
}

func (c *CredentialHandler) Login(ctx context.Context, input *authNv1.LoginRequest) (*authNv1.LoginResponse, error) {
  span := trace.SpanFromContext(ctx)

  dto := mapper.ToLoginDTO(input)
  resp, stat := c.credService.Login(ctx, &dto)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  return mapper.ToProtoLoginResponse(&resp), nil
}

func (c *CredentialHandler) Register(ctx context.Context, input *authNv1.RegisterRequest) (*emptypb.Empty, error) {
  span := trace.SpanFromContext(ctx)

  dto := mapper.ToRegisterDTO(input)
  stat := c.credService.Register(ctx, &dto)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }
  return &emptypb.Empty{}, nil
}

func (c *CredentialHandler) RefreshToken(ctx context.Context, input *authNv1.RefreshTokenRequest) (*authNv1.RefreshTokenResponse, error) {
  span := trace.SpanFromContext(ctx)

  dto := mapper.ToRefreshTokenDTO(input)
  accessToken, stat := c.credService.RefreshToken(ctx, &dto)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }
  return mapper.ToProtoRefreshTokenResponse(&accessToken), nil
}

func (c *CredentialHandler) GetCredentials(ctx context.Context, input *authNv1.GetCredentialsRequest) (*authNv1.GetCredentialsResponse, error) {
  span := trace.SpanFromContext(ctx)

  credentials, stat := c.credService.GetCredentials(ctx, input.UserId)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  responses := sharedUtil.CastSliceP(credentials, mapper.ToProtoCredential)
  return &authNv1.GetCredentialsResponse{Creds: responses}, nil
}

func (c *CredentialHandler) Logout(ctx context.Context, input *authNv1.LogoutRequest) (*emptypb.Empty, error) {
  span := trace.SpanFromContext(ctx)

  logoutDTO := mapper.ToLogoutDTO(input)
  stat := c.credService.Logout(ctx, &logoutDTO)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }
  return &emptypb.Empty{}, nil
}

func (c *CredentialHandler) LogoutAll(ctx context.Context, request *authNv1.LogoutAllRequest) (*emptypb.Empty, error) {
  span := trace.SpanFromContext(ctx)

  stats := c.credService.LogoutAll(ctx, request.UserId)
  if stats.IsError() {
    spanUtil.RecordError(stats.Error, span)
    return nil, stats.ToGRPCError()
  }
  return &emptypb.Empty{}, nil
}
