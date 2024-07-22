package external

import (
  "context"
  relationv1 "github.com/arcorium/nexa/proto/gen/go/relation/v1"
  "github.com/arcorium/nexa/shared/types"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "github.com/sony/gobreaker"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  "nexa/services/post/config"
  "nexa/services/post/internal/domain/external"
  "nexa/services/post/util"
)

func NewFollow(conn grpc.ClientConnInterface, conf *config.CircuitBreaker) external.IFollowClient {
  breaker := gobreaker.NewCircuitBreaker(gobreaker.Settings{
    Name:         "nexa-follow",
    MaxRequests:  conf.MaxRequest,
    Interval:     conf.ResetInterval,
    Timeout:      conf.OpenStateTimeout,
    IsSuccessful: util.IsGrpcConnectivityError,
  })

  return &followClient{
    client: relationv1.NewFollowServiceClient(conn),
    tracer: util.GetTracer(),
    cb:     breaker,
  }
}

type followClient struct {
  client relationv1.FollowServiceClient
  tracer trace.Tracer

  cb *gobreaker.CircuitBreaker
}

func (f *followClient) IsFollower(ctx context.Context, followerId types.Id, followedId types.Id) (bool, error) {
  ctx, span := f.tracer.Start(ctx, "FollowClient.IsFollower")
  defer span.End()

  req := relationv1.GetFollowStatusRequest{
    UserId:         followerId.String(),
    OpponentUserId: []string{followedId.String()},
  }
  res, err := f.cb.Execute(func() (interface{}, error) {
    return f.client.GetFollowStatus(ctx, &req)
  })

  if err != nil {
    err = util.CastBreakerError(err)
    spanUtil.RecordError(err, span)
    return false, err
  }

  statuses := res.(*relationv1.GetFollowStatusResponse).Status
  return statuses[0] == relationv1.FollowStatus_FOLLOWER || statuses[0] == relationv1.FollowStatus_MUTUAL, nil
}
