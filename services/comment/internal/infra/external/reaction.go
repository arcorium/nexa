package external

import (
  "context"
  reactionv1 "github.com/arcorium/nexa/proto/gen/go/reaction/v1"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "github.com/sony/gobreaker"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  "nexa/services/comment/config"
  "nexa/services/comment/internal/domain/dto"
  "nexa/services/comment/internal/domain/external"
  "nexa/services/comment/util"
)

func NewReaction(conn grpc.ClientConnInterface, conf *config.CircuitBreaker) external.IReactionClient {
  breaker := gobreaker.NewCircuitBreaker(gobreaker.Settings{
    Name:         "nexa-comment",
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

func (r *reactionClient) DeleteComments(ctx context.Context, commentIds ...types.Id) error {
  ctx, span := r.tracer.Start(ctx, "ReactionClient.DeleteComments")
  defer span.End()

  ids := sharedUtil.CastSlice(commentIds, sharedUtil.ToString[types.Id])
  req := reactionv1.DeleteItemsReactionsRequest{
    ItemType: reactionv1.Type_COMMENT,
    ItemIds:  ids,
  }
  _, err := r.cb.Execute(func() (interface{}, error) {
    return r.client.DeleteItems(ctx, &req)
  })
  if err != nil {
    spanUtil.RecordError(err, span)
    return err
  }
  return nil
}

func (r *reactionClient) GetCommentsCounts(ctx context.Context, commentIds ...types.Id) ([]dto.ReactionCountDTO, error) {
  ctx, span := r.tracer.Start(ctx, "ReactionClient.GetCommentsCount")
  defer span.End()

  ids := sharedUtil.CastSlice(commentIds, sharedUtil.ToString[types.Id])
  req := reactionv1.GetCountRequest{
    ItemType: reactionv1.Type_COMMENT,
    ItemIds:  ids,
  }
  res, err := r.cb.Execute(func() (interface{}, error) {
    return r.client.GetCount(ctx, &req)
  })
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  resp := sharedUtil.CastSlice(res.(*reactionv1.GetCountResponse).Counts, func(from *reactionv1.Count) dto.ReactionCountDTO {
    return dto.ReactionCountDTO{
      TotalLikes:    from.TotalLikes,
      TotalDislikes: from.TotalDislikes,
    }
  })
  return resp, nil
}
