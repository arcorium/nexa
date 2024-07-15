package handler

import (
  "context"
  tokenv1 "github.com/arcorium/nexa/proto/gen/go/token/v1"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  "google.golang.org/protobuf/types/known/emptypb"
  "nexa/services/token/internal/api/grpc/mapper"
  "nexa/services/token/internal/domain/service"
  "nexa/services/token/util"
)

func NewToken(token service.IToken) TokenHandler {
  return TokenHandler{
    tokenSvc: token,
    tracer:   util.GetTracer(),
  }
}

type TokenHandler struct {
  tokenv1.UnimplementedTokenServiceServer

  tokenSvc service.IToken
  tracer   trace.Tracer
}

func (t *TokenHandler) RegisterHandler(server *grpc.Server) {
  tokenv1.RegisterTokenServiceServer(server, t)
}

func (t *TokenHandler) Create(ctx context.Context, req *tokenv1.CreateTokenRequest) (*tokenv1.Token, error) {
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

func (t *TokenHandler) Verify(ctx context.Context, req *tokenv1.VerifyTokenRequest) (*tokenv1.VerifyTokenResponse, error) {
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

  return &tokenv1.VerifyTokenResponse{UserId: userId.String()}, nil
}

func (t *TokenHandler) AuthVerify(ctx context.Context, req *tokenv1.VerifyAuthTokenRequest) (*emptypb.Empty, error) {
  ctx, span := t.tracer.Start(ctx, "TokenHandler.AuthVerify")
  defer span.End()

  dtos, err := mapper.ToTokenAuthVerifyDTO(req)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  stats := t.tokenSvc.AuthVerify(ctx, &dtos)
  return nil, stats.ToGRPCErrorWithSpan(span)
}
