package external

import (
  "context"
  relationv1 "github.com/arcorium/nexa/proto/gen/go/reaction/v1"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "github.com/sony/gobreaker"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  "nexa/services/post/config"
  "nexa/services/post/internal/domain/dto"
  "nexa/services/post/internal/domain/external"
  "nexa/services/post/util"
)

func NewLike(conn grpc.ClientConnInterface, conf *config.CircuitBreaker) external.ILikeClient {
  breaker := gobreaker.NewCircuitBreaker(gobreaker.Settings{
    Name:         "nexa-like",
    MaxRequests:  conf.MaxRequest,
    Interval:     conf.ResetInterval,
    Timeout:      conf.OpenStateTimeout,
    IsSuccessful: util.IsGrpcConnectivityError,
  })

  return &likeClient{
    client: relationv1.NewReactionServiceClient(conn),
    trace:  util.GetTracer(),
    cb:     breaker,
  }
}

type likeClient struct {
  client relationv1.ReactionServiceClient
  trace  trace.Tracer
  cb     *gobreaker.CircuitBreaker
}

func (l *likeClient) GetPostCounts(ctx context.Context, postIds ...types.Id) ([]dto.LikeCountResponseDTO, error) {
  ctx, span := l.trace.Start(ctx, "LikeClient.GetPostCounts")
  defer span.End()

  req := relationv1.GetCountRequest{
    ItemType: relationv1.Type_POST,
    ItemIds:  sharedUtil.CastSlice(postIds, sharedUtil.ToString[types.Id]),
  }

  res, err := l.cb.Execute(func() (interface{}, error) {
    return l.client.GetCount(ctx, &req)
  })
  if err != nil {
    err = util.CastBreakerError(err)
    spanUtil.RecordError(err, span)
    return nil, err
  }

  return res.([]dto.LikeCountResponseDTO), nil
}

func (l *likeClient) DeletePostsLikes(ctx context.Context, postIds ...types.Id) error {
  ctx, span := l.trace.Start(ctx, "LikeClient.DeletePostsLikes")
  defer span.End()

  req := relationv1.DeleteItemsReactionsRequest{
    ItemType: relationv1.Type_POST,
    ItemIds:  sharedUtil.CastSlice(postIds, sharedUtil.ToString[types.Id]),
  }

  _, err := l.cb.Execute(func() (interface{}, error) {
    _, err := l.client.DeleteItems(ctx, &req)
    return nil, err
  })

  if err != nil {
    err = util.CastBreakerError(err)
    spanUtil.RecordError(err, span)
    return err
  }
  return nil
}
