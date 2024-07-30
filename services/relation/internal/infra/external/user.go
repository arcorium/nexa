package external

import (
  "context"
  authNv1 "github.com/arcorium/nexa/proto/gen/go/authentication/v1"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "github.com/sony/gobreaker"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  "nexa/services/relation/config"
  "nexa/services/relation/internal/domain/external"
  "nexa/services/relation/util"
)

func NewUserClient(conn grpc.ClientConnInterface, conf *config.CircuitBreaker) external.IUserClient {
  cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
    Name:         "nexa-user",
    MaxRequests:  conf.MaxRequest,
    Interval:     conf.ResetInterval,
    Timeout:      conf.OpenStateTimeout,
    IsSuccessful: util.IsGrpcConnectivityError,
  })

  return &userClient{
    client: authNv1.NewUserServiceClient(conn),
    tracer: util.GetTracer(),
    cb:     cb,
  }
}

type userClient struct {
  client authNv1.UserServiceClient
  tracer trace.Tracer
  cb     *gobreaker.CircuitBreaker
}

func (u *userClient) ValidateUsers(ctx context.Context, userIds ...types.Id) (bool, error) {
  ctx, span := u.tracer.Start(ctx, "UserClient.Validate")
  defer span.End()

  req := authNv1.FindUsersByIdsRequest{Ids: sharedUtil.CastSlice(userIds, sharedUtil.ToString[types.Id])}
  result, err := u.cb.Execute(func() (interface{}, error) {
    return u.client.FindByIds(ctx, &req)
  })
  if err != nil {
    spanUtil.RecordError(err, span)
    err = util.CastBreakerError(err)
    return false, err
  }
  resp := result.(*authNv1.FindUserByIdsResponse)
  return len(resp.Users) == len(userIds), nil
}
