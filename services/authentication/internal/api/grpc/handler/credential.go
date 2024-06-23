package handler

import (
  "context"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  "google.golang.org/grpc/codes"
  grpcStatus "google.golang.org/grpc/status"
  "google.golang.org/protobuf/types/known/emptypb"
  authv1 "nexa/proto/gen/go/authentication/v1"
  "nexa/services/authentication/constant"
  "nexa/services/authentication/internal/api/grpc/mapper"
  "nexa/services/authentication/internal/domain/service"
  "nexa/shared/auth"
  "nexa/shared/jwt"
  spanUtil "nexa/shared/span"
  "nexa/shared/status"
  "nexa/shared/types"
  "nexa/shared/util"
)

func NewCredential(svc service.ICredential, tracer trace.Tracer) CredentialHandler {
  return CredentialHandler{
    credService: svc,
    tracer:      tracer,
  }
}

type CredentialHandler struct {
  authv1.UnimplementedCredentialServiceServer
  credService service.ICredential

  tracer trace.Tracer
}

func (c *CredentialHandler) RegisterHandler(server *grpc.Server) {
  authv1.RegisterCredentialServiceServer(server, c)
}

func (c *CredentialHandler) Login(ctx context.Context, input *authv1.LoginRequest) (*authv1.LoginResponse, error) {
  span := trace.SpanFromContext(ctx)

  dto := mapper.ToLoginDTO(input)
  if err := util.ValidateStruct(ctx, &dto); err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  accessToken, stat := c.credService.Login(ctx, &dto)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  return &authv1.LoginResponse{AccessToken: accessToken}, nil
}

func (c *CredentialHandler) Register(ctx context.Context, input *authv1.RegisterRequest) (*emptypb.Empty, error) {
  span := trace.SpanFromContext(ctx)

  dto := mapper.ToRegisterDTO(input)
  if err := util.ValidateStruct(ctx, &dto); err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  stat := c.credService.Register(ctx, &dto)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }
  return &emptypb.Empty{}, nil
}

func (c *CredentialHandler) RefreshToken(ctx context.Context, input *authv1.RefreshTokenRequest) (*authv1.RefreshTokenResponse, error) {
  span := trace.SpanFromContext(ctx)

  dto := mapper.ToRefreshTokenDTO(input)
  if err := util.ValidateStruct(ctx, &dto); err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  accessToken, stat := c.credService.RefreshToken(ctx, &dto)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }
  return &authv1.RefreshTokenResponse{
    Type:        constant.TOKEN_TYPE,
    AccessToken: accessToken,
  }, nil
}

func (c *CredentialHandler) GetCredentials(ctx context.Context, input *authv1.GetCredentialsRequest) (*authv1.GetCredentialsResponse, error) {
  span := trace.SpanFromContext(ctx)

  claims, err := jwt.GetClaimsFromCtx(ctx)
  if err != nil {
    return nil, grpcStatus.Errorf(codes.Unauthenticated, err.Error())
  }

  // Roles check if the user trying to get other user credentials
  if input.UserId != claims.UserId {
    if !auth.ContainsPermissions(claims.Permissions, constant.CRED_READ_OTHERS) {
      return nil, grpcStatus.Errorf(codes.PermissionDenied, "you dont have permission to access this resource")
    }
  }

  c.credService.GetCurrentCredentials()
  credentials, stats := c.credService.GetCurrentCredentials(ctx)
  if stats.IsError() {
    return nil, stats.ToGRPCError()
  }
  return &authv1.GetCredentialsOutput{Creds: util.CastSlice(credentials, mapper.ToProtoCredential)}, nil
}

func (c *CredentialHandler) Logout(ctx context.Context, input *authv1.LogoutRequest) (*emptypb.Empty, error) {
  ids := util.CastSlice2(input.CredIds, types.IdFromString)
  // Validate
  for _, v := range ids {
    if err := v.Validate(); err != nil {
      stats := status.ErrFieldValidation(err)
      return nil, stats.ToGRPCError()
    }
  }

  stats := c.credService.Logout(ctx, ids...)
  if stats.IsError() {
    return nil, stats.ToGRPCError()
  }
  return &emptypb.Empty{}, nil
}

func (c *CredentialHandler) LogoutAll(ctx context.Context, request *authv1.LogoutAllRequest) (*emptypb.Empty, error) {
  // NOTE: User id is placed on context by interceptor, doesn't need to check here
  stats := c.credService.LogoutAll(ctx)
  if stats.IsError() {
    return nil, stats.ToGRPCError()
  }
  return &emptypb.Empty{}, nil
}
