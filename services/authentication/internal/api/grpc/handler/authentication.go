package handler

import (
  "context"
  authNv1 "github.com/arcorium/nexa/proto/gen/go/authentication/v1"
  sharedErr "github.com/arcorium/nexa/shared/errors"
  "github.com/arcorium/nexa/shared/jwt"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  "google.golang.org/protobuf/types/known/emptypb"
  "nexa/services/authentication/internal/api/grpc/mapper"
  "nexa/services/authentication/internal/domain/service"
  "nexa/services/authentication/util"
)

func NewCredential(svc service.IAuthentication) CredentialHandler {
  return CredentialHandler{
    credService: svc,
    tracer:      util.GetTracer(),
  }
}

type CredentialHandler struct {
  authNv1.UnimplementedAuthenticationServiceServer
  credService service.IAuthentication

  tracer trace.Tracer
}

func (c *CredentialHandler) RegisterHandler(server *grpc.Server) {
  authNv1.RegisterAuthenticationServiceServer(server, c)
}

func (c *CredentialHandler) Register(ctx context.Context, req *authNv1.RegisterRequest) (*emptypb.Empty, error) {
  ctx, span := c.tracer.Start(ctx, "CredentialHandler.Register")
  defer span.End()

  registerDTO, err := mapper.ToRegisterDTO(req)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  stat := c.credService.Register(ctx, &registerDTO)
  return nil, stat.ToGRPCErrorWithSpan(span)
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

  claims := types.Must(jwt.GetUserClaimsFromCtx(ctx))
  id := types.NewNullable(req.UserId)
  userId, err := types.IdFromString(id.ValueOr(claims.UserId))
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

  claims := types.Must(jwt.GetUserClaimsFromCtx(ctx))
  logoutDTO, err := mapper.ToLogoutDTO(claims, req)
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

  claims := types.Must(jwt.GetUserClaimsFromCtx(ctx))
  id := types.NewNullable(req.UserId)
  userId, err := types.IdFromString(id.ValueOr(claims.UserId))
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, sharedErr.NewFieldError("user_id", err).ToGrpcError()
  }

  stats := c.credService.LogoutAll(ctx, userId)
  return nil, stats.ToGRPCErrorWithSpan(span)
}
