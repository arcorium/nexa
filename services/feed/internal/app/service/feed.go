package service

import (
  "context"
  "errors"
  sharedErr "github.com/arcorium/nexa/shared/errors"
  sharedJwt "github.com/arcorium/nexa/shared/jwt"
  "github.com/arcorium/nexa/shared/status"
  "github.com/arcorium/nexa/shared/types"
  authUtil "github.com/arcorium/nexa/shared/util/auth"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "go.opentelemetry.io/otel/attribute"
  "go.opentelemetry.io/otel/trace"
  "nexa/services/feed/constant"
  "nexa/services/feed/internal/domain/dto"
  "nexa/services/feed/internal/domain/external"
  "nexa/services/feed/internal/domain/service"
  "nexa/services/feed/util"
)

func NewFeed(relationClient external.IRelationClient, postClient external.IPostClient) service.IFeed {
  return &feedService{
    relClient:  relationClient,
    postClient: postClient,
    trace:      util.GetTracer(),
  }
}

type feedService struct {
  relClient  external.IRelationClient
  postClient external.IPostClient

  trace trace.Tracer
}

func (f *feedService) getUserClaims(ctx context.Context) (types.Id, error) {
  claims, err := sharedJwt.GetUserClaimsFromCtx(ctx)
  if err != nil {
    return types.NullId(), err
  }

  userId, err := types.IdFromString(claims.UserId)
  if err != nil {
    return types.NullId(), err
  }

  return userId, nil
}

func (f *feedService) checkPermission(ctx context.Context, targetId types.Id, permission string) error {
  // Validate permission
  claims, err := sharedJwt.GetUserClaimsFromCtx(ctx)
  if err != nil {
    return sharedErr.ErrUnauthenticated
  }

  if !targetId.EqWithString(claims.UserId) {
    // Need permission to update other users
    if !authUtil.ContainsPermission(claims.Roles, permission) {
      return sharedErr.ErrUnauthorizedPermission
    }
  }

  return nil
}

func (f *feedService) GetUserFeed(ctx context.Context, userId types.Id, limit uint64) ([]dto.PostResponseDTO, status.Object) {
  ctx, span := f.trace.Start(ctx, "FeedService.GetUserFeed", trace.WithAttributes(attribute.String("user_id", userId.String())))
  defer span.End()

  // Check permissions
  if err := f.checkPermission(ctx, userId, constant.FEED_PERMISSIONS[constant.FEED_GET_ARB]); err != nil {
    spanUtil.RecordError(err, span)
    return nil, status.ErrUnAuthorized(err)
  }

  // Get user followings
  followings, err := f.relClient.GetFollowings(ctx, userId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, status.ErrExternal(err)
  }

  if len(followings.UserIds) == 0 {
    spanUtil.RecordError(errors.New("user has no followings"), span)
    return nil, status.ErrNotFound()
  }

  postResp, err := f.postClient.GetUsers(ctx, limit, followings.UserIds...)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, status.ErrExternal(err)
  }

  return postResp.Posts, status.Success()
}
