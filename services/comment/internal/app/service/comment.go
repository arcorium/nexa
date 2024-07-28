package service

import (
  "context"
  sharedDto "github.com/arcorium/nexa/shared/dto"
  sharedErr "github.com/arcorium/nexa/shared/errors"
  sharedJwt "github.com/arcorium/nexa/shared/jwt"
  "github.com/arcorium/nexa/shared/optional"
  "github.com/arcorium/nexa/shared/status"
  "github.com/arcorium/nexa/shared/types"
  sharedUOW "github.com/arcorium/nexa/shared/uow"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  authUtil "github.com/arcorium/nexa/shared/util/auth"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "go.opentelemetry.io/otel/trace"
  "nexa/services/comment/constant"
  "nexa/services/comment/internal/app/uow"
  "nexa/services/comment/internal/domain/dto"
  "nexa/services/comment/internal/domain/entity"
  "nexa/services/comment/internal/domain/external"
  "nexa/services/comment/internal/domain/mapper"
  "nexa/services/comment/internal/domain/service"
  "nexa/services/comment/util"
  "nexa/services/comment/util/errs"
)

func NewComment(unit sharedUOW.IUnitOfWork[uow.CommentStorage], postClient external.IPostClient, /*reactClient external.IReactionClient*/) service.IComment {
  return &commentService{
    unit:   unit,
    tracer: util.GetTracer(),
    //reactClient: reactClient,
    postClient: postClient,
  }
}

type commentService struct {
  unit   sharedUOW.IUnitOfWork[uow.CommentStorage]
  tracer trace.Tracer

  //reactClient external.IReactionClient
  postClient external.IPostClient
}

func (c *commentService) getUserClaims(ctx context.Context) (types.Id, error) {
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

func (c *commentService) checkPermission(ctx context.Context, targetId types.Id, permission string) error {
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

func (c *commentService) checkAvailability(ctx context.Context, postId types.Id) error {
  exists, err := c.postClient.Validate(ctx, postId)
  if err != nil {
    return err
  }
  return sharedUtil.Ternary(exists, nil, errs.ErrPostNotFound)
}

func (c *commentService) Create(ctx context.Context, commentDTO *dto.CreateCommentDTO) (types.Id, status.Object) {
  ctx, span := c.tracer.Start(ctx, "CommentService.Create")
  defer span.End()

  userId, err := c.getUserClaims(ctx)
  if err != nil {
    spanUtil.RecordError(err, span)
    return types.NullId(), status.ErrUnAuthorized(err)
  }

  // Validate post
  if err := c.checkAvailability(ctx, commentDTO.PostId); err != nil {
    spanUtil.RecordError(err, span)
    return types.NullId(), status.ErrExternal(err)
  }

  comment, err := commentDTO.ToDomain(userId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return types.NullId(), status.ErrInternal(err)
  }

  repos := c.unit.Repositories()
  err = repos.Comment().Create(ctx, &comment)
  if err != nil {
    spanUtil.RecordError(err, span)
    return types.NullId(), status.FromRepository2(err, status.Null, optional.Some(status.ErrNotFound()))
  }

  return comment.Id, status.Created()
}

func (c *commentService) Edit(ctx context.Context, commentDTO *dto.EditCommentDTO) status.Object {
  ctx, span := c.tracer.Start(ctx, "CommentService.Edit")
  defer span.End()

  userId, err := c.getUserClaims(ctx)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }

  repos := c.unit.Repositories()
  err = repos.Comment().UpdateContent(ctx, userId, commentDTO.Id, commentDTO.Content)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, optional.Some(status.INTERNAL_SERVER_ERROR))
  }

  return status.Updated()
}

func (c *commentService) deleteArbitrary(ctx context.Context, commentIds ...types.Id) status.Object {
  span := trace.SpanFromContext(ctx)

  stat := status.Deleted()
  _ = c.unit.DoTx(ctx, func(ctx context.Context, storage uow.CommentStorage) error {
    err := storage.Comment().DeleteByIds(ctx, commentIds...)
    if err != nil {
      spanUtil.RecordError(err, span)
      stat = status.FromRepository(err, optional.Some(status.INTERNAL_SERVER_ERROR))
      return err
    }

    //err = c.reactClient.DeleteComments(ctx, commentIds...)
    //if err != nil {
    //  spanUtil.RecordError(err, span)
    //  stat = status.ErrExternal(err)
    //  return err
    //}

    return nil
  })

  return stat
}

func (c *commentService) Delete(ctx context.Context, commentIds ...types.Id) status.Object {
  ctx, span := c.tracer.Start(ctx, "CommentService.Delete")
  defer span.End()

  claims, err := sharedJwt.GetUserClaimsFromCtx(ctx)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrUnAuthorized(err)
  }

  userId, err := types.IdFromString(claims.UserId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrUnAuthorized(err)
  }

  if authUtil.ContainsPermission(claims.Roles, constant.COMMENT_PERMISSIONS[constant.COMMENT_DELETE_ARB]) {
    return c.deleteArbitrary(ctx, commentIds...)
  }

  stat := status.Deleted()
  _ = c.unit.DoTx(ctx, func(ctx context.Context, storage uow.CommentStorage) error {
    ctx, span := c.tracer.Start(ctx, "UOW.Delete")
    defer span.End()

    // Delete user's comments
    _, err = storage.Comment().DeleteUsers(ctx, userId, commentIds...)
    if err != nil {
      spanUtil.RecordError(err, span)
      stat = status.FromRepository(err, status.NullCode)
      return err
    }

    // Delete each comment's reactions
    //err = c.reactClient.DeleteComments(ctx, commentIds...)
    //if err != nil {
    //  spanUtil.RecordError(err, span)
    //  stat = status.ErrExternal(err)
    //  return err
    //}
    return nil
  })

  return stat
}

func (c *commentService) FindById(ctx context.Context, findDTO *dto.FindCommentByIdDTO) (dto.CommentResponseDTO, status.Object) {
  ctx, span := c.tracer.Start(ctx, "CommentService.FindById")
  defer span.End()

  repos := c.unit.Repositories()
  comment, err := repos.Comment().FindById(ctx, findDTO.ShowReply, findDTO.CommentId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.CommentResponseDTO{}, status.FromRepository(err, status.NullCode)
  }

  // Validate post
  if err := c.checkAvailability(ctx, comment[0].PostId); err != nil { // Index 0 is the comment itself (main parent)
    spanUtil.RecordError(err, span)
    return dto.CommentResponseDTO{}, status.ErrExternal(err)
  }

  resp := mapper.ToCommentsResponse(comment)
  if len(resp) > 1 {
    err = errs.ErrResultWithDifferentLength
    spanUtil.RecordError(err, span)
    return dto.CommentResponseDTO{}, status.ErrInternal(err)
  }
  return resp[0], status.Success()
}

func (c *commentService) GetPosts(ctx context.Context, getDTO *dto.GetPostsCommentsDTO, pageDTO sharedDto.PagedElementDTO) (sharedDto.PagedElementResult[dto.CommentResponseDTO], status.Object) {
  ctx, span := c.tracer.Start(ctx, "CommentService.GetPosts")
  defer span.End()

  // Validate post
  if err := c.checkAvailability(ctx, getDTO.PostId); err != nil {
    spanUtil.RecordError(err, span)
    return sharedDto.PagedElementResult[dto.CommentResponseDTO]{}, status.ErrExternal(err)
  }

  repos := c.unit.Repositories()
  result, err := repos.Comment().FindByPostId(ctx, getDTO.ShowReply, getDTO.PostId, pageDTO.ToQueryParam())
  if err != nil {
    spanUtil.RecordError(err, span)
    return sharedDto.NewPagedElementResult2[dto.CommentResponseDTO](nil, &pageDTO, result.Total), status.FromRepository(err, status.NullCode)
  }

  // Get each comment's reactions
  //commentIds := sharedUtil.CastSliceP(result.Data, func(from *entity.Comment) types.Id {
  //  return from.Id
  //})
  //reactions, err := c.reactClient.GetCommentsCounts(ctx, commentIds...)
  //if err != nil {
  //  spanUtil.RecordError(err, span)
  //  return sharedDto.NewPagedElementResult2[dto.CommentResponseDTO](nil, &pageDTO, result.Total), status.ErrExternal(err)
  //}
  //
  // Error different length
  //if len(reactions) != len(commentIds) {
  //  err = errs.ErrResultWithDifferentLength
  //  spanUtil.RecordError(err, span)
  //  return sharedDto.NewPagedElementResult2[dto.CommentResponseDTO](nil, &pageDTO, result.Total), status.ErrInternal(err)
  //}

  resp := mapper.ToCommentsResponse(result.Data)
  return sharedDto.NewPagedElementResult2(resp, &pageDTO, result.Total), status.Success()
}

func (c *commentService) GetReplies(ctx context.Context, repliesDTO *dto.GetCommentsRepliesDTO, pageDTO sharedDto.PagedElementDTO) (sharedDto.PagedElementResult[dto.CommentResponseDTO], status.Object) {
  ctx, span := c.tracer.Start(ctx, "CommentService.GetReplies")
  defer span.End()

  repos := c.unit.Repositories()
  comment, err := repos.Comment().FindById(ctx, false, repliesDTO.CommentId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return sharedDto.PagedElementResult[dto.CommentResponseDTO]{}, status.FromRepository(err, status.NullCode)
  }

  if err := c.checkAvailability(ctx, comment[0].PostId); err != nil {
    spanUtil.RecordError(err, span)
    return sharedDto.PagedElementResult[dto.CommentResponseDTO]{}, status.ErrExternal(err)
  }

  result, err := repos.Comment().GetReplies(ctx, repliesDTO.ShowReply, repliesDTO.CommentId, pageDTO.ToQueryParam())
  if err != nil {
    spanUtil.RecordError(err, span)
    return sharedDto.NewPagedElementResult2[dto.CommentResponseDTO](nil, &pageDTO, result.Total), status.FromRepository(err, status.NullCode)
  }

  //// Get each comment's reactions
  //commentIds := sharedUtil.CastSliceP(result.Data, func(from *entity.Comment) types.Id {
  //  return from.Id
  //})
  //reactions, err := c.reactClient.GetCommentsCounts(ctx, commentIds...)
  //if err != nil {
  //  spanUtil.RecordError(err, span)
  //  return sharedDto.NewPagedElementResult2[dto.CommentResponseDTO](nil, &pageDTO, result.Total), status.ErrExternal(err)
  //}
  //
  //// Error different length
  //if len(reactions) != len(commentIds) {
  //  spanUtil.RecordError(err, span)
  //  return sharedDto.NewPagedElementResult2[dto.CommentResponseDTO](nil, &pageDTO, result.Total), status.ErrInternal(err)
  //}

  resp := mapper.ToCommentsResponse(result.Data)
  return sharedDto.NewPagedElementResult2(resp, &pageDTO, result.Total), status.Success()
}

func (c *commentService) GetCounts(ctx context.Context, itemType entity.ItemType, itemIds ...types.Id) ([]uint64, status.Object) {
  ctx, span := c.tracer.Start(ctx, "CommentService.GetReplies")
  defer span.End()

  var counts []entity.Count
  var err error
  repos := c.unit.Repositories()

  switch itemType {
  case entity.ItemPostComment:
    counts, err = repos.Comment().GetPostCounts(ctx, itemIds...)
  case entity.ItemCommentReply:
    counts, err = repos.Comment().GetReplyCounts(ctx, itemIds...)
  default:
    return nil, status.ErrInternal(sharedErr.ErrEnumOutOfBounds)
  }

  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, status.FromRepository(err, status.NullCode)
  }

  resp := sharedUtil.CastSliceP(counts, func(from *entity.Count) uint64 {
    return from.TotalComments
  })
  return resp, status.Success()
}

func (c *commentService) ClearPosts(ctx context.Context, postIds ...types.Id) status.Object {
  ctx, span := c.tracer.Start(ctx, "CommentService.ClearPosts")
  defer span.End()

  for _, id := range postIds {
    if err := c.checkAvailability(ctx, id); err != nil {
      spanUtil.RecordError(err, span)
      return status.ErrExternal(err)
    }
  }

  stat := status.Deleted()
  _ = c.unit.DoTx(ctx, func(ctx context.Context, storage uow.CommentStorage) error {
    // Clear post's comments
    _, err := storage.Comment().DeleteByPostIds(ctx, postIds...)
    if err != nil {
      spanUtil.RecordError(err, span)
      stat = status.FromRepository(err, status.NullCode)
      return err
    }

    // Delete each comment's reactions
    //err = c.reactClient.DeleteComments(ctx, commentIds...)
    //if err != nil {
    //  spanUtil.RecordError(err, span)
    //  stat = status.ErrExternal(err)
    //  return err
    //}

    return nil
  })

  return stat
}

func (c *commentService) ClearUsers(ctx context.Context, userId types.Id) status.Object {
  ctx, span := c.tracer.Start(ctx, "CommentService.ClearUsers")
  defer span.End()

  if err := c.checkPermission(ctx, userId, constant.COMMENT_PERMISSIONS[constant.COMMENT_DELETE_ARB]); err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrUnAuthorized(err)
  }

  stat := status.Deleted()
  _ = c.unit.DoTx(ctx, func(ctx context.Context, storage uow.CommentStorage) error {
    // Delete all comment's created by the user
    _, err := storage.Comment().DeleteUsers(ctx, userId)
    if err != nil {
      spanUtil.RecordError(err, span)
      stat = status.FromRepository(err, status.NullCode)
      return err
    }

    // Delete each comment's reactions
    //err = c.reactClient.DeleteComments(ctx, commentIds...)
    //if err != nil {
    //  spanUtil.RecordError(err, span)
    //  stat = status.ErrExternal(err)
    //  return err
    //}

    return nil
  })

  return stat
}

func (c *commentService) IsExists(ctx context.Context, commentIds ...types.Id) (bool, status.Object) {
  ctx, span := c.tracer.Start(ctx, "CommentService.IsExists")
  defer span.End()

  repos := c.unit.Repositories()
  count, err := repos.Comment().Count(ctx, commentIds...)
  if err != nil {
    spanUtil.RecordError(err, span)
    return false, status.FromRepository(err, status.NullCode)
  }

  return count == uint64(len(commentIds)), status.Success()
}
