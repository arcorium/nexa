package handler

import (
  "context"
  "go.opentelemetry.io/otel/trace"
  emptypb "google.golang.org/protobuf/types/known/emptypb"
  authZv1 "nexa/proto/gen/go/authorization/v1"
  "nexa/services/authorization/internal/api/grpc/mapper"
  "nexa/services/authorization/internal/domain/service"
)

func NewAuthorization(authorization service.IAuthorization) AuthorizationHandler {
  return AuthorizationHandler{
    authZSvc: authorization,
  }
}

type AuthorizationHandler struct {
  authZv1.UnimplementedAuthorizationServiceServer

  authZSvc service.IAuthorization
}

func (a AuthorizationHandler) CheckUserPermission(ctx context.Context, request *authZv1.CheckUserRequest) (*emptypb.Empty, error) {
  span := trace.SpanFromContext(ctx)

  dtos := mapper.ToIsAuthorizationDTO(request)
  stat := a.authZSvc.IsAuthorized(ctx, &dtos)
  return nil, stat.ToGRPCErrorWithSpan(span)
}
