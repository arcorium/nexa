package service

import (
  "context"
  sharedDto "github.com/arcorium/nexa/shared/dto"
  sharedErr "github.com/arcorium/nexa/shared/errors"
  sharedJwt "github.com/arcorium/nexa/shared/jwt"
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

func NewBlock(block repository.IBlock, client external.IUserClient) service.IBlock {
  return &blockService{
    repo:       block,
    userClient: client,
    tracer:     util.GetTracer(),
  }
}

type blockService struct {
  repo       repository.IBlock
  userClient external.IUserClient
  tracer     trace.Tracer
}

func (b *blockService) getUserClaims(ctx context.Context) (types.Id, error) {
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

func (b *blockService) checkPermission(ctx context.Context, targetId types.Id, permission string) error {
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

func (b *blockService) isUserExists(ctx context.Context, userId types.Id) error {
  // Validate the user id
  exist, err := b.userClient.ValidateUsers(ctx, userId)
  if err != nil {
    return err
  }

  return sharedUtil.Ternary(exist, nil, errs.ErrUserNotFound)
}

func (b *blockService) Block(ctx context.Context, targetUserId types.Id) status.Object {
  ctx, span := b.tracer.Start(ctx, "BlockService.Block")
  defer span.End()

  // Check trying to block itself
  userId, err := b.getUserClaims(ctx)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrUnAuthorized(err)
  }

  if userId.Eq(targetUserId) {
    spanUtil.RecordError(errs.ErrBlockItself, span)
    return status.ErrBadRequest(errs.ErrBlockItself)
  }

  // Validate user id
  if err := b.isUserExists(ctx, targetUserId); err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrExternal(err)
  }

  block := entity.NewBlock(userId, targetUserId)
  err = b.repo.Create(ctx, &block)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository2(err, status.Null, status.SomeSuccess())
  }

  return status.Created()
}

func (b *blockService) Unblock(ctx context.Context, targetUserId types.Id) status.Object {
  ctx, span := b.tracer.Start(ctx, "BlockService.Unblock")
  defer span.End()

  userId, err := b.getUserClaims(ctx)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrUnAuthorized(err)
  }

  block := entity.NewBlock(userId, targetUserId)
  err = b.repo.Delete(ctx, &block)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepositoryOverrideObject(err, status.Deleted())

  }

  return status.Deleted()
}

func (b *blockService) IsBlocked(ctx context.Context, userId types.Id, targetUserId types.Id) (bool, status.Object) {
  ctx, span := b.tracer.Start(ctx, "BlockService.IsBlocked")
  defer span.End()

  if err := b.isUserExists(ctx, targetUserId); err != nil {
    spanUtil.RecordError(err, span)
    return false, status.ErrBadRequest(err)
  }

  blocked, err := b.repo.IsBlocked(ctx, userId, targetUserId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return false, status.FromRepository(err, status.NullCode)
  }
  return blocked, status.Success()
}

func (b *blockService) GetUsers(ctx context.Context, userId types.Id, pageDTO sharedDto.PagedElementDTO) (sharedDto.PagedElementResult[dto.BlockResponseDTO], status.Object) {
  ctx, span := b.tracer.Start(ctx, "BlockService.GetUsers")
  defer span.End()

  result, err := b.repo.GetBlocked(ctx, userId, pageDTO.ToQueryParam())
  if err != nil {
    spanUtil.RecordError(err, span)
    return sharedDto.NewPagedElementResult2[dto.BlockResponseDTO](nil, &pageDTO, result.Total), status.FromRepositoryOverrideObject(err, status.Success())
  }

  resp := sharedUtil.CastSliceP(result.Data, mapper.ToBlockResponseDTO)
  return sharedDto.NewPagedElementResult2(resp, &pageDTO, result.Total), status.Success()
}

func (b *blockService) GetUsersCount(ctx context.Context, userId types.Id) (dto.BlockCountResponseDTO, status.Object) {
  ctx, span := b.tracer.Start(ctx, "BlockService.GetUsersCount")
  defer span.End()
  var err error

  counts, err := b.repo.GetCounts(ctx, userId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.BlockCountResponseDTO{}, status.FromRepository(err, status.NullCode)
  }

  resp := sharedUtil.CastSliceP(counts, mapper.ToBlockCountResponseDTO)
  return resp[0], status.Success()
}

func (b *blockService) ClearUsers(ctx context.Context, userId types.Id) status.Object {
  ctx, span := b.tracer.Start(ctx, "BlockService.ClearUsers")
  defer span.End()

  if err := b.checkPermission(ctx, userId, constant.RELATION_PERMISSIONS[constant.RELATION_DELETE_BLOCK_ARB]); err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrUnAuthorized(err)
  }

  err := b.repo.DeleteByUserId(ctx, true, userId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepositoryOverrideObject(err, status.Deleted())
  }

  return status.Deleted()
}
