package external

import (
  "context"
  "github.com/arcorium/nexa/proto/gen/go/common"
  postv1 "github.com/arcorium/nexa/proto/gen/go/post/v1"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "github.com/redis/go-redis/v9"
  "github.com/sony/gobreaker"
  "go.opentelemetry.io/otel/attribute"
  "go.opentelemetry.io/otel/trace"
  "golang.org/x/exp/slices"
  "google.golang.org/grpc"
  "google.golang.org/grpc/codes"
  "google.golang.org/grpc/status"
  "nexa/services/feed/config"
  "nexa/services/feed/internal/domain/dto"
  "nexa/services/feed/internal/domain/entity"
  "nexa/services/feed/internal/domain/external"
  "nexa/services/feed/util"
  "time"
)

func NewPostClient(conn grpc.ClientConnInterface, redis redis.UniversalClient, reactClient external.IReactionClient, commentClient external.ICommentClient, conf *config.CircuitBreaker) external.IPostClient {
  breaker := gobreaker.NewCircuitBreaker(gobreaker.Settings{
    Name:         "nexa-post",
    MaxRequests:  conf.MaxRequest,
    Interval:     conf.ResetInterval,
    Timeout:      conf.OpenStateTimeout,
    IsSuccessful: util.IsGrpcConnectivityError,
  })

  return &postClient{
    client:      postv1.NewPostServiceClient(conn),
    reactClient: reactClient,
    commClient:  commentClient,
    cache:       redis,
    tracer:      util.GetTracer(),
    cb:          breaker,
  }
}

type postClient struct {
  client      postv1.PostServiceClient
  reactClient external.IReactionClient
  commClient  external.ICommentClient
  cache       redis.UniversalClient
  tracer      trace.Tracer
  cb          *gobreaker.CircuitBreaker
}

func toPostResponseDTO(post *postv1.Post, totalComment, likeCount, dislikeCount int64) dto.PostResponseDTO {
  var parent *dto.PostResponseDTO
  if post.ParentPost != nil {
    temp := toPostResponseDTO(post, -1, -1, -1)
    parent = &temp
  }

  var lastEditedAt time.Time
  if post.LastEdited != nil {
    lastEditedAt = post.LastEdited.AsTime()
  }

  var createdAt time.Time
  if post.CreatedAt != nil {
    createdAt = post.CreatedAt.AsTime()
  }

  taggedUserIds, _ := sharedUtil.CastSliceErrs(post.TaggedUserIds, types.IdFromString)

  return dto.PostResponseDTO{
    Id:            types.Must(types.IdFromString(post.Id)),
    Parent:        parent,
    CreatorId:     types.Must(types.IdFromString(post.CreatorId)),
    Content:       post.Content,
    Visibility:    types.Must(entity.NewVisibility(uint8(post.Visibility))),
    TotalLikes:    likeCount,
    TotalDislikes: dislikeCount,
    TotalComments: totalComment,
    TaggedUserIds: taggedUserIds,
    MediaUrls:     post.MediaUrls,
    LastEditedAt:  lastEditedAt,
    CreatedAt:     createdAt,
  }
}

func (p *postClient) GetUsers(ctx context.Context, limit uint64, userIds ...types.Id) (dto.GetUsersPostResponseDTO, error) {
  ids := sharedUtil.CastSlice(userIds, sharedUtil.ToString[types.Id])
  ctx, span := p.tracer.Start(ctx, "PostClient.GetUsers", trace.WithAttributes(attribute.StringSlice("user_ids", ids)))
  defer span.End()

  limits := int64(limit)
  isNoLimit := limits == 0
  var posts []*postv1.Post
  _, err := p.cb.Execute(func() (interface{}, error) {
    for _, id := range ids {
      req := &postv1.FindUserPostRequest{
        UserId: &id,
        Details: &common.PagedElementInput{
          Element: uint64(limits),
          Page:    1,
        },
      }

      // TODO: Get data from cache

      post, err := p.client.FindUsers(ctx, req)
      if err != nil {
        s, ok := status.FromError(err)
        if !ok || s.Code() != codes.NotFound {
          return nil, err
        }
        continue
      }

      // TODO: Set cache

      limits -= int64(post.Details.Element)
      limits = max(limits, 0)
      posts = append(posts, post.Posts...)
      if limits <= 0 && !isNoLimit {
        break
      }
    }
    return nil, nil
  })
  if err != nil {
    spanUtil.RecordError(err, span)
    err = util.CastBreakerError(err)
    return dto.GetUsersPostResponseDTO{}, err
  }

  // Sort DESC
  slices.SortFunc(posts, func(a, b *postv1.Post) bool {
    return a.CreatedAt.AsTime().After(b.CreatedAt.AsTime())
  })

  postIds := sharedUtil.CastSlice(posts, func(from *postv1.Post) types.Id {
    return types.Must(types.IdFromString(from.Id))
  })
  // Get comment counts
  commentCounts, err := p.commClient.GetPostCommentCounts(ctx, postIds...)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.GetUsersPostResponseDTO{}, err
  }

  // Get reaction counts
  reactionCounts, err := p.reactClient.GetPostReactionCounts(ctx, postIds...)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.GetUsersPostResponseDTO{}, err
  }

  resp := make([]dto.PostResponseDTO, 0, len(posts))
  for i := range posts {
    currentReactCount := &reactionCounts[i]
    casted := toPostResponseDTO(posts[i],
      int64(commentCounts[i]),
      int64(currentReactCount.TotalLikes),
      int64(currentReactCount.TotalDislikes))

    resp = append(resp, casted)
  }
  return dto.GetUsersPostResponseDTO{Posts: resp}, nil
}
