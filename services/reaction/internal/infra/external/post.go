package external

import (
  "context"
  postv1 "github.com/arcorium/nexa/proto/gen/go/post/v1"
  "github.com/arcorium/nexa/shared/types"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "github.com/sony/gobreaker"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  "nexa/services/reaction/config"
  "nexa/services/reaction/internal/domain/external"
  "nexa/services/reaction/util"
)

func NewPostClient(conn grpc.ClientConnInterface, conf *config.CircuitBreaker) external.IPostClient {
  cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
    Name:         "nexa-post",
    MaxRequests:  conf.MaxRequest,
    Interval:     conf.ResetInterval,
    Timeout:      conf.OpenStateTimeout,
    IsSuccessful: util.IsGrpcConnectivityError,
  })

  return &postClient{
    client: postv1.NewPostServiceClient(conn),
    tracer: util.GetTracer(),
    cb:     cb,
  }
}

type postClient struct {
  client postv1.PostServiceClient
  tracer trace.Tracer
  cb     *gobreaker.CircuitBreaker
}

func (c *postClient) ValidatePost(ctx context.Context, postId types.Id) (bool, error) {
  ctx, span := c.tracer.Start(ctx, "CommentClient.Validate")
  defer span.End()

  req := &postv1.FindPostByIdRequest{PostId: postId.String()}
  result, err := c.cb.Execute(func() (interface{}, error) {
    return c.client.FindById(ctx, req)
  })
  if err != nil {
    err = util.CastBreakerError(err)
    spanUtil.RecordError(err, span)
    return false, err
  }
  resp := result.(*postv1.FindPostByIdResponse)
  return resp.Post != nil, nil
}
