package handler

import (
  "context"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  "google.golang.org/protobuf/types/known/emptypb"
  authZv1 "nexa/proto/gen/go/authorization/v1"
  "nexa/services/authorization/internal/api/grpc/mapper"
  "nexa/services/authorization/internal/domain/service"
  "nexa/services/authorization/util"
  spanUtil "nexa/shared/util/span"
)

func NewAuthorization(authorization service.IAuthorization) AuthorizationHandler {
  return AuthorizationHandler{
    authZSvc: authorization,
    tracer:   util.GetTracer(),
  }
}

type AuthorizationHandler struct {
  authZv1.UnimplementedAuthorizationServiceServer

  authZSvc service.IAuthorization
  tracer   trace.Tracer
}

func (a *AuthorizationHandler) Register(server *grpc.Server) {
  authZv1.RegisterAuthorizationServiceServer(server, a)
}

func (a *AuthorizationHandler) CheckUserPermission(ctx context.Context, request *authZv1.CheckUserRequest) (*emptypb.Empty, error) {
  ctx, span := a.tracer.Start(ctx, "AuthorizationHandler.CheckUserPermission")
  defer span.End()

  dtos, err := mapper.ToIsAuthorizationDTO(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }
  stat := a.authZSvc.IsAuthorized(ctx, &dtos)
  return nil, stat.ToGRPCErrorWithSpan(span)
}
