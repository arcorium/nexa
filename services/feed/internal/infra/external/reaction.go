package external

import (
  "context"
  reactionv1 "github.com/arcorium/nexa/proto/gen/go/reaction/v1"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "github.com/sony/gobreaker"
  "go.opentelemetry.io/otel/attribute"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  "nexa/services/feed/config"
  "nexa/services/feed/internal/domain/dto"
  "nexa/services/feed/internal/domain/external"
  "nexa/services/feed/util"
)

func NewReactionClient(conn grpc.ClientConnInterface, conf *config.CircuitBreaker) external.IReactionClient {
  breaker := gobreaker.NewCircuitBreaker(gobreaker.Settings{
    Name:         "nexa-reaction",
    MaxRequests:  conf.MaxRequest,
    Interval:     conf.ResetInterval,
    Timeout:      conf.OpenStateTimeout,
    IsSuccessful: util.IsGrpcConnectivityError,
  })

  return &reactionClient{
    client: reactionv1.NewReactionServiceClient(conn),
    tracer: util.GetTracer(),
    cb:     breaker,
  }
}

type reactionClient struct {
  client reactionv1.ReactionServiceClient
  tracer trace.Tracer
  cb     *gobreaker.CircuitBreaker
}

func (c *reactionClient) GetPostReactionCounts(ctx context.Context, postIds ...types.Id) ([]dto.PostReactionCountDTO, error) {
  ids := sharedUtil.CastSlice(postIds, sharedUtil.ToString[types.Id])
  ctx, span := c.tracer.Start(ctx, "ReactionClient.GetPostReactionCounts", trace.WithAttributes(attribute.StringSlice("post_ids", ids)))
  defer span.End()

  req := &reactionv1.GetCountRequest{
    ItemType: reactionv1.Type_POST,
    ItemIds:  ids,
  }
  result, err := c.cb.Execute(func() (interface{}, error) {
    return c.client.GetCount(ctx, req)
  })
  if err != nil {
    spanUtil.RecordError(err, span)
    err = util.CastBreakerError(err)
    return nil, err
  }
  counts := result.(*reactionv1.GetCountResponse).Counts

  resp := sharedUtil.CastSlice(counts, func(count *reactionv1.Count) dto.PostReactionCountDTO {
    return dto.PostReactionCountDTO{
      TotalLikes:    count.TotalLikes,
      TotalDislikes: count.TotalDislikes,
    }
  })
  return resp, nil
}
