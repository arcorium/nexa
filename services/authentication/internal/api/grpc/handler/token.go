package handler

import (
  "context"
  authNv1 "github.com/arcorium/nexa/proto/gen/go/authentication/v1"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  "nexa/services/authentication/internal/api/grpc/mapper"
  "nexa/services/authentication/internal/domain/service"
  "nexa/services/authentication/util"
)

func NewToken(token service.IToken) TokenHandler {
  return TokenHandler{
    tokenSvc: token,
    tracer:   util.GetTracer(),
  }
}

type TokenHandler struct {
  authNv1.UnimplementedTokenServiceServer

  tokenSvc service.IToken
  tracer   trace.Tracer
}

func (t *TokenHandler) RegisterHandler(server *grpc.Server) {
  authNv1.RegisterTokenServiceServer(server, t)
}

func (t *TokenHandler) Create(ctx context.Context, req *authNv1.TokenCreateRequest) (*authNv1.Token, error) {
  ctx, span := t.tracer.Start(ctx, "TokenHandler.Create")
  defer span.End()

  dtos, err := mapper.ToCreateTokenDTO(req)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  result, stats := t.tokenSvc.Request(ctx, &dtos)
  if stats.IsError() {
    spanUtil.RecordError(stats.Error, span)
    return nil, stats.ToGRPCError()
  }
  return mapper.ToProtoToken(&result), nil
}

func (t *TokenHandler) Verify(ctx context.Context, req *authNv1.TokenVerifyRequest) (*authNv1.TokenVerifyResponse, error) {
  ctx, span := t.tracer.Start(ctx, "TokenHandler.Verify")
  defer span.End()

  dtos, err := mapper.ToTokenVerifyDTO(req)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  userId, stats := t.tokenSvc.Verify(ctx, &dtos)
  if stats.IsError() {
    spanUtil.RecordError(stats.Error, span)
    return nil, stats.ToGRPCError()
  }

  return &authNv1.TokenVerifyResponse{UserId: userId.String()}, nil
}
