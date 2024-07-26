package service

import (
  "context"
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
  "nexa/services/relation/constant"
  "nexa/services/relation/internal/domain/dto"
  "nexa/services/relation/internal/domain/entity"
  "nexa/services/relation/internal/domain/external"
  "nexa/services/relation/internal/domain/mapper"
  "nexa/services/relation/internal/domain/repository"
  "nexa/services/relation/internal/domain/service"
  "nexa/services/relation/util"
  "nexa/services/relation/util/errs"
)

func NewFollow(follow repository.IFollow, client external.IUserClient) service.IFollow {
  return &followService{
    repo:       follow,
    userClient: client,
    tracer:     util.GetTracer(),
  }
}

type followService struct {
  repo       repository.IFollow
  userClient external.IUserClient
  tracer     trace.Tracer
}

func (f *followService) getUserClaims(ctx context.Context) (types.Id, error) {
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

func (f *followService) checkPermission(ctx context.Context, targetId types.Id, permission string) error {
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

func (f *followService) isUserExists(ctx context.Context, userIds ...types.Id) error {
  // Validate the user id
  exist, err := f.userClient.Validate(ctx, userIds...)
  if err != nil {
    return err
  }

  return sharedUtil.Ternary(exist, nil, errs.ErrUserNotFound)
}

func (f *followService) Follow(ctx context.Context, targetUserId ...types.Id) status.Object {
  ctx, span := f.tracer.Start(ctx, "FollowService.Follow")
  defer span.End()

  userId, err := f.getUserClaims(ctx)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrUnAuthorized(err)
  }

  // Prevent self-follow
  for _, id := range targetUserId {
    if userId.Eq(id) {
      err = errs.ErrFollowItself
      spanUtil.RecordError(err, span)
      return status.ErrBadRequest(err)
    }
  }

  // Validate user id
  if err := f.isUserExists(ctx, targetUserId...); err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrBadRequest(err)
  }

  follows := entity.NewFollows(userId, targetUserId...)
  err = f.repo.Creates(ctx, follows)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository2(err, status.Null, optional.Some(status.Success()))
  }

  return status.Success()
}

func (f *followService) Unfollow(ctx context.Context, targetUserId ...types.Id) status.Object {
  ctx, span := f.tracer.Start(ctx, "FollowService.Unfollow")
  defer span.End()

  userId, err := f.getUserClaims(ctx)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrUnAuthorized(err)
  }

  follows := entity.NewFollows(userId, targetUserId...)
  err = f.repo.Delete(ctx, follows...)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepositoryOverride(err, types.NewPair[status.Code, error](status.DELETED, nil))
  }

  return status.Deleted()
}

func (f *followService) GetFollowers(ctx context.Context, userId types.Id, pageDTO sharedDto.PagedElementDTO) (sharedDto.PagedElementResult[dto.FollowResponseDTO], status.Object) {
  ctx, span := f.tracer.Start(ctx, "FollowService.GetFollowers")
  defer span.End()

  claims := types.Must(sharedJwt.GetUserClaimsFromCtx(ctx))
  if !userId.EqWithString(claims.UserId) {
    // Validate user id
    if err := f.isUserExists(ctx, userId); err != nil {
      spanUtil.RecordError(err, span)
      return sharedDto.NewPagedElementResult2[dto.FollowResponseDTO](nil, &pageDTO, 0), status.ErrExternal(err)
    }
  }

  followers, err := f.repo.GetFollowers(ctx, userId, pageDTO.ToQueryParam())
  if err != nil {
    spanUtil.RecordError(err, span)
    return sharedDto.NewPagedElementResult2[dto.FollowResponseDTO](nil, &pageDTO, followers.Total), status.FromRepositoryOverrideObject(err, status.Success())
  }

  resp := sharedUtil.CastSliceP(followers.Data, mapper.ToFollowerResponseDTO)
  return sharedDto.NewPagedElementResult2(resp, &pageDTO, followers.Total), status.Success()
}

func (f *followService) GetFollowings(ctx context.Context, userId types.Id, pageDTO sharedDto.PagedElementDTO) (sharedDto.PagedElementResult[dto.FollowResponseDTO], status.Object) {
  ctx, span := f.tracer.Start(ctx, "FollowService.GetFollowings")
  defer span.End()

  claims := types.Must(sharedJwt.GetUserClaimsFromCtx(ctx))
  if !userId.EqWithString(claims.UserId) {
    // Validate user id
    if err := f.isUserExists(ctx, userId); err != nil {
      spanUtil.RecordError(err, span)
      return sharedDto.NewPagedElementResult2[dto.FollowResponseDTO](nil, &pageDTO, 0), status.ErrExternal(err)
    }
  }

  followers, err := f.repo.GetFollowings(ctx, userId, pageDTO.ToQueryParam())
  if err != nil {
    spanUtil.RecordError(err, span)
    return sharedDto.NewPagedElementResult2[dto.FollowResponseDTO](nil, &pageDTO, followers.Total), status.FromRepositoryOverrideObject(err, status.Success())
  }

  resp := sharedUtil.CastSliceP(followers.Data, mapper.ToFolloweeResponseDTO)
  return sharedDto.NewPagedElementResult2(resp, &pageDTO, followers.Total), status.Success()
}

func (f *followService) GetStatus(ctx context.Context, userId types.Id, targetUserIds ...types.Id) ([]entity.FollowStatus, status.Object) {
  // NOTE: Currently only return either the user is following the target or not
  ctx, span := f.tracer.Start(ctx, "FollowService.GetFollowings")
  defer span.End()

  // Validate user id
  claims := types.Must(sharedJwt.GetUserClaimsFromCtx(ctx))
  if !userId.EqWithString(claims.UserId) {
    if err := f.isUserExists(ctx, userId); err != nil {
      spanUtil.RecordError(err, span)
      return nil, status.ErrExternal(err)
    }
  }

  // Validate target user ids
  //if err := f.isUserExists(ctx, targetUserIds...); err != nil {
  //  spanUtil.RecordError(err, span)
  //  return nil, status.ErrExternal(err)
  //}

  result, err := f.repo.IsFollowing(ctx, userId, targetUserIds...)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, status.Object{}
  }

  resp := sharedUtil.CastSlice(result, mapper.ToFollowStatus)
  return resp, status.Success()
}

func (f *followService) GetUsersCount(ctx context.Context, userIds ...types.Id) ([]dto.FollowCountResponseDTO, status.Object) {
  ctx, span := f.tracer.Start(ctx, "FollowService.GetUsersCount")
  defer span.End()

  // Validate
  //if err := f.isUserExists(ctx, userIds...); err != nil {
  //  spanUtil.RecordError(err, span)
  //  return nil, status.ErrExternal(err)
  //}

  counts, err := f.repo.GetCounts(ctx, userIds...)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, status.FromRepository(err, status.NullCode)
  }

  resp := sharedUtil.CastSliceP(counts, mapper.ToFollowCountResponseDTO)
  return resp, status.Success()
}

func (f *followService) ClearUsers(ctx context.Context, userId types.Id) status.Object {
  ctx, span := f.tracer.Start(ctx, "FollowService.ClearUsers")
  defer span.End()

  if err := f.checkPermission(ctx, userId, constant.RELATION_PERMISSIONS[constant.RELATION_DELETE_FOLLOW_ARB]); err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrUnAuthorized(err)
  }

  err := f.repo.DeleteByUserId(ctx, true, userId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }

  return status.Deleted()
}
