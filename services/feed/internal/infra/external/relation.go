package external

import (
  "context"
  "fmt"
  relationv1 "github.com/arcorium/nexa/proto/gen/go/relation/v1"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "github.com/redis/go-redis/v9"
  "github.com/sony/gobreaker"
  "go.opentelemetry.io/otel/attribute"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  "nexa/services/feed/config"
  "nexa/services/feed/internal/domain/dto"
  "nexa/services/feed/internal/domain/external"
  "nexa/services/feed/util"
  "time"
)

func NewRelationClient(conn grpc.ClientConnInterface, redis redis.UniversalClient, conf *config.CircuitBreaker) external.IRelationClient {
  breaker := gobreaker.NewCircuitBreaker(gobreaker.Settings{
    Name:         "nexa-relation",
    MaxRequests:  conf.MaxRequest,
    Interval:     conf.ResetInterval,
    Timeout:      conf.OpenStateTimeout,
    IsSuccessful: util.IsGrpcConnectivityError,
  })

  return &relationClient{
    client: relationv1.NewFollowServiceClient(conn),
    cache:  redis,
    cb:     breaker,
    tracer: util.GetTracer(),
  }
}

const ttl = time.Minute * 5

type relationClient struct {
  client relationv1.FollowServiceClient
  cache  redis.UniversalClient
  tracer trace.Tracer
  cb     *gobreaker.CircuitBreaker
}

func (r *relationClient) key(userId types.Id) string {
  return fmt.Sprintf("rel:%s", userId.String())
}

func (r *relationClient) keyCount(userId types.Id) string {
  return fmt.Sprintf("rel:%s:count", userId.String())
}

func (r *relationClient) getFromCache(ctx context.Context, userId types.Id) ([]string, error) {
  key := r.key(userId)
  countKey := r.keyCount(userId)

  // Get total followings
  get := r.cache.Get(ctx, countKey)
  if err := get.Err(); err != nil {
    return nil, err
  }
  followingCount, err := get.Int64()
  if err != nil {
    return nil, err
  }

  result := r.cache.SRandMemberN(ctx, key, followingCount)
  if err := result.Err(); err != nil {
    return nil, err
  }

  return result.Result()
}

func (r *relationClient) setToCache(ctx context.Context, userId types.Id, followingIds ...string) error {
  key := r.key(userId)
  countKey := r.keyCount(userId)

  cmds, _ := r.cache.TxPipelined(ctx, func(p redis.Pipeliner) error {
    p.SetNX(ctx, countKey, len(followingIds), ttl)
    p.SAdd(ctx, key, sharedUtil.CastSlice(followingIds, sharedUtil.ToAny[string]))
    p.Expire(ctx, key, ttl)
    return nil
  })

  for _, cmd := range cmds {
    if cmd.Err() != nil {
      return cmd.Err()
    }
  }

  return nil
}

func (r *relationClient) GetFollowings(ctx context.Context, userId types.Id) (dto.GetFollowingResponseDTO, error) {
  ctx, span := r.tracer.Start(ctx, "RelationClient.GetFollowings", trace.WithAttributes(attribute.String("user_id", userId.String())))
  defer span.End()

  //cacheResult, err := r.getFromCache(ctx, userId)
  //if err == nil {
  //  ids, ierr := sharedUtil.CastSliceErrs(cacheResult, types.IdFromString)
  //  return dto.GetFollowingResponseDTO{UserIds: ids}, ierr
  //}

  id := userId.String()
  req := relationv1.GetUserFolloweesRequest{
    UserId: &id,
  }

  result, err := r.cb.Execute(func() (interface{}, error) {
    return r.client.GetFollowees(ctx, &req)
  })
  if err != nil {
    spanUtil.RecordError(err, span)
    err = util.CastBreakerError(err)
    return dto.GetFollowingResponseDTO{}, err
  }

  ids := result.(*relationv1.GetUserFolloweesResponse).UserIds
  userIds, ierr := sharedUtil.CastSliceErrs(ids, types.IdFromString)
  if !ierr.IsNil() {
    if ierr.IsEmptySlice() {
      return dto.GetFollowingResponseDTO{}, nil
    }
    spanUtil.RecordError(ierr, span)
    return dto.GetFollowingResponseDTO{}, ierr
  }

  //err = r.setToCache(ctx, userId, ids...)
  //if err != nil {
  //  spanUtil.RecordError(err, span)
  //  return dto.GetFollowingResponseDTO{}, err
  //}

  return dto.GetFollowingResponseDTO{UserIds: userIds}, nil
}
