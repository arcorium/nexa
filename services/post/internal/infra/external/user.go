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
  "nexa/services/post/config"
  "nexa/services/post/internal/domain/external"
  "nexa/services/post/util"
)

func NewUser(conn grpc.ClientConnInterface, conf *config.CircuitBreaker) external.IUserClient {
  breaker := gobreaker.NewCircuitBreaker(gobreaker.Settings{
    Name:         "nexa-user",
    MaxRequests:  conf.MaxRequest,
    Interval:     conf.ResetInterval,
    Timeout:      conf.OpenStateTimeout,
    IsSuccessful: util.IsGrpcConnectivityError,
  })

  return &userClient{
    client: authNv1.NewUserServiceClient(conn),
    tracer: util.GetTracer(),
    cb:     breaker,
  }
}

type userClient struct {
  client authNv1.UserServiceClient
  tracer trace.Tracer

  cb *gobreaker.CircuitBreaker
}

func (u *userClient) GetUserNames(ctx context.Context, userIds ...types.Id) ([]string, error) {
  ctx, span := u.tracer.Start(ctx, "UserClient.GetUserNames")
  defer span.End()

  ids := sharedUtil.CastSlice(userIds, sharedUtil.ToString[types.Id])

  res, err := u.cb.Execute(func() (interface{}, error) {
    return u.client.FindByIds(ctx, &authNv1.FindUsersByIdsRequest{
      Ids: ids,
    })
  })

  if err != nil {
    err = util.CastBreakerError(err)
    spanUtil.RecordError(err, span)
    return nil, err
  }

  // Get only the username
  resp := sharedUtil.CastSlice(res.(*authNv1.FindUserByIdsResponse).Users, func(user *authNv1.User) string {
    return user.Username
  })

  return resp, nil
}
