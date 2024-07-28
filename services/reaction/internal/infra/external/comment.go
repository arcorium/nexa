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
  "nexa/services/reaction/config"
  "nexa/services/reaction/internal/domain/external"
  "nexa/services/reaction/util"
)

func NewCommentClient(conn grpc.ClientConnInterface, conf *config.CircuitBreaker) external.ICommentClient {
  cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
    Name:         "nexa-comment",
    MaxRequests:  conf.MaxRequest,
    Interval:     conf.ResetInterval,
    Timeout:      conf.OpenStateTimeout,
    IsSuccessful: util.IsGrpcConnectivityError,
  })

  return &commentClient{
    client: commentv1.NewCommentServiceClient(conn),
    tracer: util.GetTracer(),
    cb:     cb,
  }
}

type commentClient struct {
  client commentv1.CommentServiceClient
  tracer trace.Tracer
  cb     *gobreaker.CircuitBreaker
}

func (c *commentClient) Validate(ctx context.Context, commentIds ...types.Id) (bool, error) {
  ctx, span := c.tracer.Start(ctx, "CommentClient.Validate")
  defer span.End()

  ids := sharedUtil.CastSlice(commentIds, sharedUtil.ToString[types.Id])
  req := &commentv1.IsCommentExistRequest{CommentIds: ids}
  result, err := c.cb.Execute(func() (interface{}, error) {
    return c.client.IsExist(ctx, req)
  })
  if err != nil {
    err = util.CastBreakerError(err)
    spanUtil.RecordError(err, span)
    return false, err
  }
  resp := result.(*commentv1.IsCommentExistResponse)
  return resp.IsExists, nil
}
