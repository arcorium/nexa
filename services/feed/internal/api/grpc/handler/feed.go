package handler

import (
  "context"
  feedv1 "github.com/arcorium/nexa/proto/gen/go/feed/v1"
  sharedErr "github.com/arcorium/nexa/shared/errors"
  sharedJwt "github.com/arcorium/nexa/shared/jwt"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  "nexa/services/feed/internal/api/grpc/mapper"
  "nexa/services/feed/internal/domain/service"
  "nexa/services/feed/util"
)

func NewFeed(feed service.IFeed) FeedHandler {
  return FeedHandler{
    svc:    feed,
    tracer: util.GetTracer(),
  }
}

type FeedHandler struct {
  feedv1.UnimplementedFeedServiceServer
  svc    service.IFeed
  tracer trace.Tracer
}

func (f *FeedHandler) Register(server *grpc.Server) {
  feedv1.RegisterFeedServiceServer(server, f)
}

func (f *FeedHandler) GetUserFeed(ctx context.Context, request *feedv1.GetUserFeedRequest) (*feedv1.GetUserFeedResponse, error) {
  ctx, span := f.tracer.Start(ctx, "FeedHandler.GetUserFeed")
  defer span.End()

  claims := types.Must(sharedJwt.GetUserClaimsFromCtx(ctx))
  id := types.NewNullable(request.UserId)
  userId, err := types.IdFromString(id.ValueOr(claims.UserId))
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, sharedErr.NewFieldError("user_id", err).ToGrpcError()
  }

  posts, stat := f.svc.GetUserFeed(ctx, userId, types.NewNullable(request.Limit).ValueOr(0))
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  resp := sharedUtil.CastSliceP(posts, mapper.ToProtoPostWithCount)
  return &feedv1.GetUserFeedResponse{Posts: resp, Element: uint64(len(posts))}, nil
}
