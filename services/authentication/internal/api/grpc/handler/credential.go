package handler

import (
  "context"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  "google.golang.org/protobuf/types/known/emptypb"
  authNv1 "nexa/proto/gen/go/authentication/v1"
  "nexa/services/authentication/internal/api/grpc/mapper"
  "nexa/services/authentication/internal/domain/service"
  "nexa/services/authentication/util"
  sharedErr "nexa/shared/errors"
  "nexa/shared/types"
  sharedUtil "nexa/shared/util"
  spanUtil "nexa/shared/util/span"
)

func NewCredential(svc service.ICredential) CredentialHandler {
  return CredentialHandler{
    credService: svc,
    tracer:      util.GetTracer(),
  }
}

type CredentialHandler struct {
  authNv1.UnimplementedCredentialServiceServer
  credService service.ICredential

  tracer trace.Tracer
}

func (c *CredentialHandler) RegisterHandler(server *grpc.Server) {
  authNv1.RegisterCredentialServiceServer(server, c)
}

func (c *CredentialHandler) Login(ctx context.Context, req *authNv1.LoginRequest) (*authNv1.LoginResponse, error) {
  ctx, span := c.tracer.Start(ctx, "CredentialHandler.Login")
  defer span.End()

  dtos, err := mapper.ToLoginDTO(req)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  resp, stat := c.credService.Login(ctx, &dtos)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  return mapper.ToProtoLoginResponse(&resp), nil
}

func (c *CredentialHandler) Register(ctx context.Context, req *authNv1.RegisterRequest) (*emptypb.Empty, error) {
  ctx, span := c.tracer.Start(ctx, "CredentialHandler.Register")
  defer span.End()

  dtos, err := mapper.ToRegisterDTO(req)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  stat := c.credService.Register(ctx, &dtos)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (c *CredentialHandler) RefreshToken(ctx context.Context, req *authNv1.RefreshTokenRequest) (*authNv1.RefreshTokenResponse, error) {
  ctx, span := c.tracer.Start(ctx, "CredentialHandler.RefreshToken")
  defer span.End()

  dto, err := mapper.ToRefreshTokenDTO(req)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  accessToken, stat := c.credService.RefreshToken(ctx, &dto)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }
  return mapper.ToProtoRefreshTokenResponse(&accessToken), nil
}

func (c *CredentialHandler) GetCredentials(ctx context.Context, req *authNv1.GetCredentialsRequest) (*authNv1.GetCredentialsResponse, error) {
  ctx, span := c.tracer.Start(ctx, "CredentialHandler.GetCredentials")
  defer span.End()

  userId, err := types.IdFromString(req.UserId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, sharedErr.NewFieldError("user_id", err).ToGrpcError()
  }

  credentials, stat := c.credService.GetCredentials(ctx, userId)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  responses := sharedUtil.CastSliceP(credentials, mapper.ToProtoCredential)
  return &authNv1.GetCredentialsResponse{Creds: responses}, nil
}

func (c *CredentialHandler) Logout(ctx context.Context, req *authNv1.LogoutRequest) (*emptypb.Empty, error) {
  ctx, span := c.tracer.Start(ctx, "CredentialHandler.Logout")
  defer span.End()

  logoutDTO, err := mapper.ToLogoutDTO(req)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  stat := c.credService.Logout(ctx, &logoutDTO)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (c *CredentialHandler) LogoutAll(ctx context.Context, req *authNv1.LogoutAllRequest) (*emptypb.Empty, error) {
  ctx, span := c.tracer.Start(ctx, "CredentialHandler.LogoutAll")
  defer span.End()

  userId, err := types.IdFromString(req.UserId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, sharedErr.NewFieldError("user_id", err).ToGrpcError()
  }

  stats := c.credService.LogoutAll(ctx, userId)
  return nil, stats.ToGRPCErrorWithSpan(span)
}
