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

func NewRelationClient(conn grpc.ClientConnInterface, conf *config.CircuitBreaker) external.IRelationClient {
  breaker := gobreaker.NewCircuitBreaker(gobreaker.Settings{
    Name:         "nexa-relation",
    MaxRequests:  conf.MaxRequest,
    Interval:     conf.ResetInterval,
    Timeout:      conf.OpenStateTimeout,
    IsSuccessful: util.IsGrpcConnectivityError,
  })

  return &relationClient{
    followClient: relationv1.NewFollowServiceClient(conn),
    blockClient:  relationv1.NewBlockServiceClient(conn),
    tracer:       util.GetTracer(),
    cb:           breaker,
  }
}

type relationClient struct {
  followClient relationv1.FollowServiceClient
  blockClient  relationv1.BlockServiceClient
  tracer       trace.Tracer

  cb *gobreaker.CircuitBreaker
}

func (f *relationClient) IsFollower(ctx context.Context, followerId types.Id, followedId types.Id) (bool, error) {
  ctx, span := f.tracer.Start(ctx, "RelationClient.IsFollower")
  defer span.End()

  req := &relationv1.GetRelationRequest{
    //UserId:         , // Use forwarded authorization
    OpponentUserIds: []string{followedId.String()},
  }
  res, err := f.cb.Execute(func() (interface{}, error) {
    return f.followClient.GetRelation(ctx, req)
  })
  if err != nil {
    spanUtil.RecordError(err, span)
    err = util.CastBreakerError(err)
    return false, err
  }

  statuses := res.(*relationv1.GetRelationResponse).Status
  return statuses[0] == relationv1.Relation_FOLLOWER || statuses[0] == relationv1.Relation_MUTUAL, nil
}

func (f *relationClient) IsBlocked(ctx context.Context, blockerId types.Id) (bool, error) {
  ctx, span := f.tracer.Start(ctx, "RelationClient.IsBlocked")
  defer span.End()

  req := &relationv1.IsUserBlockedRequest{TargetUserId: blockerId.String()}
  blocked, err := f.cb.Execute(func() (interface{}, error) {
    return f.blockClient.IsBlocked(ctx, req)
  })
  if err != nil {
    spanUtil.RecordError(err, span)
    err = util.CastBreakerError(err)
    return false, err
  }

  return blocked.(*relationv1.IsUserBlockedResponse).IsBlocked, nil
}
