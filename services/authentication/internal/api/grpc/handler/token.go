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
)

func NewToken(token service.IToken) TokenHandler {
  return TokenHandler{
    tokenSvc: token,
  }
}

type TokenHandler struct {
  authNv1.UnimplementedTokenServiceServer

  tokenSvc service.IToken
}

func (t *TokenHandler) RegisterHandler(server *grpc.Server) {
  authNv1.RegisterTokenServiceServer(server, t)
}

func (t *TokenHandler) Create(ctx context.Context, req *authNv1.TokenCreateRequest) (*authNv1.Token, error) {
  span := trace.SpanFromContext(ctx)

  dto := mapper.ToCreateTokenDTO(req)
  result, stats := t.tokenSvc.Request(ctx, &dto)
  if stats.IsError() {
    spanUtil.RecordError(stats.Error, span)
    return nil, stats.ToGRPCError()
  }
  return mapper.ToProtoToken(&result), nil
}

func (t *TokenHandler) Verify(ctx context.Context, req *authNv1.TokenVerifyRequest) (*emptypb.Empty, error) {
  span := trace.SpanFromContext(ctx)

  dto := mapper.ToTokenVerifyDTO(req)
  stats := t.tokenSvc.Verify(ctx, &dto)
  if stats.IsError() {
    spanUtil.RecordError(stats.Error, span)
    return nil, stats.ToGRPCError()
  }

  return &emptypb.Empty{}, nil
}
