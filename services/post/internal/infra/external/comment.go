package external

import (
  "context"
  commentv1 "github.com/arcorium/nexa/proto/gen/go/comment/v1"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "github.com/sony/gobreaker"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  "nexa/services/post/config"
  "nexa/services/post/internal/domain/external"
  "nexa/services/post/util"
)

func NewComment(conn grpc.ClientConnInterface, conf *config.CircuitBreaker) external.ICommentClient {
  breaker := gobreaker.NewCircuitBreaker(gobreaker.Settings{
    Name:         "nexa-comment",
    MaxRequests:  conf.MaxRequest,
    Interval:     conf.ResetInterval,
    Timeout:      conf.OpenStateTimeout,
    IsSuccessful: util.IsGrpcConnectivityError,
  })

  return &commentClient{
    client: commentv1.NewCommentServiceClient(conn),
    trace:  util.GetTracer(),
    cb:     breaker,
  }
}

type commentClient struct {
  client commentv1.CommentServiceClient
  trace  trace.Tracer
  cb     *gobreaker.CircuitBreaker
}

func (c *commentClient) GetPostCounts(ctx context.Context, postIds ...types.Id) ([]uint64, error) {
  ctx, span := c.trace.Start(ctx, "CommentClient.GetPostCounts")
  defer span.End()

  req := commentv1.GetCountsRequest{
    ItemType: commentv1.Type_POST_COMMENT,
    ItemIds:  sharedUtil.CastSlice(postIds, sharedUtil.ToString[types.Id]),
  }

  res, err := c.cb.Execute(func() (interface{}, error) {
    return c.client.GetCounts(ctx, &req)
  })
  if err != nil {
    err = util.CastBreakerError(err)
    spanUtil.RecordError(err, span)
    return nil, err
  }

  return res.([]uint64), nil
}

func (c *commentClient) DeletePostsComments(ctx context.Context, postIds ...types.Id) error {
  ctx, span := c.trace.Start(ctx, "CommentClient.DeletePostComments")
  defer span.End()

  req := commentv1.ClearPostsCommentsRequest{
    PostIds: sharedUtil.CastSlice(postIds, sharedUtil.ToString[types.Id]),
  }

  _, err := c.cb.Execute(func() (interface{}, error) {
    _, err := c.client.ClearPosts(ctx, &req)
    return nil, err
  })

  if err != nil {
    err = util.CastBreakerError(err)
    spanUtil.RecordError(err, span)
    return err
  }

  return nil
}
