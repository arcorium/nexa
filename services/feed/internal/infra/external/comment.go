package external

import (
  "context"
  commentv1 "github.com/arcorium/nexa/proto/gen/go/comment/v1"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "github.com/sony/gobreaker"
  "go.opentelemetry.io/otel/attribute"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  "nexa/services/feed/config"
  "nexa/services/feed/internal/domain/external"
  "nexa/services/feed/util"
)

func NewCommentClient(conn grpc.ClientConnInterface, conf *config.CircuitBreaker) external.ICommentClient {
  breaker := gobreaker.NewCircuitBreaker(gobreaker.Settings{
    Name:         "nexa-comment",
    MaxRequests:  conf.MaxRequest,
    Interval:     conf.ResetInterval,
    Timeout:      conf.OpenStateTimeout,
    IsSuccessful: util.IsGrpcConnectivityError,
  })

  return &commentClient{
    client: commentv1.NewCommentServiceClient(conn),
    tracer: util.GetTracer(),
    cb:     breaker,
  }
}

type commentClient struct {
  client commentv1.CommentServiceClient
  tracer trace.Tracer
  cb     *gobreaker.CircuitBreaker
}

func (c *commentClient) GetPostCommentCounts(ctx context.Context, postIds ...types.Id) ([]uint64, error) {
  ids := sharedUtil.CastSlice(postIds, sharedUtil.ToString[types.Id])
  ctx, span := c.tracer.Start(ctx, "CommentClient.GetPostCommentCounts", trace.WithAttributes(attribute.StringSlice("post_ids", ids)))
  defer span.End()

  req := &commentv1.GetCountsRequest{
    ItemType: commentv1.Type_POST_COMMENT,
    ItemIds:  ids,
  }

  result, err := c.cb.Execute(func() (interface{}, error) {
    return c.client.GetCounts(ctx, req)
  })
  if err != nil {
    spanUtil.RecordError(err, span)
    err = util.CastBreakerError(err)
    return nil, err
  }

  return result.(*commentv1.GetCountsResponse).Total, nil
}
