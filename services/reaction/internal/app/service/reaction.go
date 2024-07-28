package service

import (
  "context"
  "errors"
  sharedDto "github.com/arcorium/nexa/shared/dto"
  sharedErr "github.com/arcorium/nexa/shared/errors"
  sharedJwt "github.com/arcorium/nexa/shared/jwt"
  "github.com/arcorium/nexa/shared/optional"
  "github.com/arcorium/nexa/shared/status"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  authUtil "github.com/arcorium/nexa/shared/util/auth"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "go.opentelemetry.io/otel/trace"
  "nexa/services/reaction/constant"
  "nexa/services/reaction/internal/domain/dto"
  "nexa/services/reaction/internal/domain/entity"
  "nexa/services/reaction/internal/domain/external"
  "nexa/services/reaction/internal/domain/mapper"
  "nexa/services/reaction/internal/domain/repository"
  "nexa/services/reaction/internal/domain/service"
  "nexa/services/reaction/util"
  "nexa/services/reaction/util/errs"
)

func NewReaction(reaction repository.IReaction, postClient external.IPostClient, commentClient external.ICommentClient) service.IReaction {
  return &reactionService{
    repo:          reaction,
    postClient:    postClient,
    commentClient: commentClient,
    tracer:        util.GetTracer(),
  }
}

type reactionService struct {
  repo          repository.IReaction
  postClient    external.IPostClient
  commentClient external.ICommentClient

  tracer trace.Tracer
}

func (r *reactionService) getUserClaims(ctx context.Context) (types.Id, error) {
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

func (r *reactionService) checkPermission(ctx context.Context, targetId types.Id, permission string) error {
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

func (r *reactionService) isExists(ctx context.Context, itemType entity.ItemType, itemId types.Id) error {
  var exists bool
  var err error
  if itemType == entity.ItemPost {
    exists, err = r.postClient.Validate(ctx, itemId)
  } else if itemType == entity.ItemComment {
    exists, err = r.commentClient.Validate(ctx, itemId)
  }
  if err != nil {
    return err
  }
  return sharedUtil.Ternary(exists, nil, errs.ErrItemNotFound)
}

func (r *reactionService) Like(ctx context.Context, itemType entity.ItemType, itemId types.Id) status.Object {
  ctx, span := r.tracer.Start(ctx, "ReactionService.Like")
  defer span.End()

  userId, err := r.getUserClaims(ctx)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrUnAuthorized(err)
  }

  // Validate item
  if err := r.isExists(ctx, itemType, itemId); err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrExternal(err)
  }

  domain := entity.NewLikeReaction(userId, itemType, itemId)
  err = r.repo.Delsert(ctx, &domain)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepositoryOverride(err, types.NewPair(status.INTERNAL_SERVER_ERROR, errors.New("something wrong in Delsert")))
  }

  return status.Success()
}

func (r *reactionService) Dislike(ctx context.Context, itemType entity.ItemType, itemId types.Id) status.Object {
  ctx, span := r.tracer.Start(ctx, "ReactionService.Dislike")
  defer span.End()

  userId, err := r.getUserClaims(ctx)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrUnAuthorized(err)
  }

  // Validate item
  if err := r.isExists(ctx, itemType, itemId); err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrExternal(err)
  }

  domain := entity.NewDislikeReaction(userId, itemType, itemId)
  err = r.repo.Delsert(ctx, &domain)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepositoryOverride(err, types.NewPair(status.INTERNAL_SERVER_ERROR, errors.New("something wrong in Delsert")))
  }

  return status.Success()
}

func (r *reactionService) GetItemsReactions(ctx context.Context, itemType entity.ItemType, itemId types.Id, pageDTO sharedDto.PagedElementDTO) (sharedDto.PagedElementResult[dto.ReactionResponseDTO], status.Object) {
  ctx, span := r.tracer.Start(ctx, "ReactionService.GetItemsReactions")
  defer span.End()

  // Validate item
  if err := r.isExists(ctx, itemType, itemId); err != nil {
    spanUtil.RecordError(err, span)
    return sharedDto.NewPagedElementResult2[dto.ReactionResponseDTO](nil, &pageDTO, 0), status.ErrExternal(err)
  }

  result, err := r.repo.FindByItemId(ctx, itemType, itemId, pageDTO.ToQueryParam())
  if err != nil {
    spanUtil.RecordError(err, span)
    return sharedDto.NewPagedElementResult2[dto.ReactionResponseDTO](nil, &pageDTO, result.Total), status.FromRepository(err, status.NullCode)
  }

  resp := sharedUtil.CastSliceP(result.Data, mapper.ToReactionResponseDTO)
  return sharedDto.NewPagedElementResult2(resp, &pageDTO, result.Total), status.Success()
}

func (r *reactionService) GetCounts(ctx context.Context, itemType entity.ItemType, itemIds ...types.Id) ([]dto.CountResponseDTO, status.Object) {
  ctx, span := r.tracer.Start(ctx, "ReactionService.GetCounts")
  defer span.End()

  counts, err := r.repo.GetCounts(ctx, itemType, itemIds...)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, status.FromRepository(err, status.NullCode)
  }

  resp := sharedUtil.CastSliceP(counts, mapper.ToCountResponseDTO)
  return resp, status.Success()
}

func (r *reactionService) Delete(ctx context.Context, itemType entity.ItemType, itemIds ...types.Id) status.Object {
  ctx, span := r.tracer.Start(ctx, "ReactionService.Delete")
  defer span.End()

  err := r.repo.DeleteByItemId(ctx, itemType, itemIds...)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }
  return status.Deleted()
}

func (r *reactionService) ClearUsers(ctx context.Context, userId types.Id) status.Object {
  ctx, span := r.tracer.Start(ctx, "ReactionService.ClearUsers")
  defer span.End()

  // Check if user id the same
  if err := r.checkPermission(ctx, userId, constant.REACTION_PERMISSIONS[constant.REACTION_DELETE_ARB]); err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrUnAuthorized(err)
  }

  err := r.repo.DeleteByUserId(ctx, userId, optional.Null[entity.Item]())
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }
  return status.Deleted()
}
