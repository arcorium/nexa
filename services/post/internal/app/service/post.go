package service

import (
  "context"
  sharedDto "github.com/arcorium/nexa/shared/dto"
  sharedErr "github.com/arcorium/nexa/shared/errors"
  sharedJwt "github.com/arcorium/nexa/shared/jwt"
  "github.com/arcorium/nexa/shared/status"
  "github.com/arcorium/nexa/shared/types"
  sharedUow "github.com/arcorium/nexa/shared/uow"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  authUtil "github.com/arcorium/nexa/shared/util/auth"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "go.opentelemetry.io/otel/trace"
  "nexa/services/post/constant"
  uow "nexa/services/post/internal/app/uow"
  "nexa/services/post/internal/domain/dto"
  "nexa/services/post/internal/domain/entity"
  "nexa/services/post/internal/domain/external"
  "nexa/services/post/internal/domain/mapper"
  "nexa/services/post/internal/domain/service"
  "nexa/services/post/util"
)

func NewPost(unit sharedUow.IUnitOfWork[uow.PostStorage]) service.IPost {
  return &postService{
    unit:   unit,
    tracer: util.GetTracer(),
  }
}

type postService struct {
  unit   sharedUow.IUnitOfWork[uow.PostStorage]
  tracer trace.Tracer

  mediaClient   external.IMediaStore
  userClient    external.IUserClient
  followClient  external.IFollowClient
  likeClient    external.ILikeClient
  commentClient external.ICommentClient
}

func (p *postService) getUserClaims(ctx context.Context) (types.Id, error) {
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

func (p *postService) checkPermission(ctx context.Context, targetId types.Id, permission string) error {
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

func (p *postService) getExternalPostData(ctx context.Context, posts ...entity.Post) error {
  // TODO: The post versions media and tagged user is not get
  // Get file ids
  var fileIds []types.Id
  for _, post := range posts {
    ids := sharedUtil.CastSliceP(post.Medias, func(from *entity.Media) types.Id {
      return from.Id
    })
    fileIds = append(fileIds, ids...)
  }

  urls, err := p.mediaClient.GetUrls(ctx, fileIds...)
  if err != nil {
    return err
  }

  // Get user tags
  var userIds []types.Id
  for _, post := range posts {
    ids := sharedUtil.CastSliceP(post.Tags, func(user *entity.TaggedUser) types.Id {
      return user.Id
    })
    userIds = append(userIds, ids...)
  }

  names, err := p.userClient.GetUserNames(ctx, userIds...)
  if err != nil {
    return err
  }

  // Get post likes
  postIds := sharedUtil.CastSliceP(posts, func(post *entity.Post) types.Id {
    return post.Id
  })
  likeCounts, err := p.likeClient.GetPostCounts(ctx, postIds...)
  if err != nil {
    return err
  }
  // Get post comments
  commentCounts, err := p.commentClient.GetPostCounts(ctx, postIds...)
  if err != nil {
    return err
  }

  // Set into domain
  mediaIdx := 0
  userIdx := 0
  for i := range len(posts) {
    // Set data
    posts[i].Likes = likeCounts[i].TotalLike
    posts[i].Dislikes = likeCounts[i].TotalDislike
    posts[i].Comments = commentCounts[i]

    // Append media urls
    for j := range len(posts[i].Medias) {
      posts[i].Medias[i].Url = urls[mediaIdx+j]
    }
    mediaIdx += len(posts[i].Medias)

    // Append user tags
    for j := range len(posts[i].Tags) {
      posts[i].Tags[i].Name = names[userIdx+j]
    }
    userIdx += len(posts[i].Tags)
  }

  return nil
}

func (p *postService) GetAll(ctx context.Context, pageDTO *sharedDto.PagedElementDTO) (sharedDto.PagedElementResult[dto.PostResponseDTO], status.Object) {
  ctx, span := p.tracer.Start(ctx, "PostService.GetAll")
  defer span.End()

  repos := p.unit.Repositories()
  result, err := repos.Post().Get(ctx, false, entity.VisibilityOnlyMe, pageDTO.ToQueryParam())
  if err != nil {
    spanUtil.RecordError(err, span)
    return sharedDto.PagedElementResult[dto.PostResponseDTO]{}, status.FromRepository(err, status.NullCode)
  }

  err = p.getExternalPostData(ctx, result.Data...)
  if err != nil {
    spanUtil.RecordError(err, span)
    return sharedDto.PagedElementResult[dto.PostResponseDTO]{}, status.ErrExternal(err)
  }

  resp := sharedUtil.CastSliceP(result.Data, mapper.ToPostResponseDTO)
  return sharedDto.NewPagedElementResult2(resp, pageDTO, result.Total), status.Success()
}

func (p *postService) GetEdited(ctx context.Context, postId types.Id) (dto.EditedPostResponseDTO, status.Object) {
  ctx, span := p.tracer.Start(ctx, "PostService.GetEdited")
  defer span.End()

  repos := p.unit.Repositories()
  post, err := repos.Post().FindById(ctx, true, postId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.EditedPostResponseDTO{}, status.FromRepository(err, status.NullCode)
  }

  err = p.getExternalPostData(ctx, post...)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.EditedPostResponseDTO{}, status.ErrExternal(err)
  }

  resp := mapper.ToEditedPostResponseDTO(&post[0])
  return resp, status.Success()
}

func (p *postService) FindById(ctx context.Context, id types.Id) (dto.PostResponseDTO, status.Object) {
  ctx, span := p.tracer.Start(ctx, "PostService.FindById")
  defer span.End()

  repos := p.unit.Repositories()
  post, err := repos.Post().FindById(ctx, false, id)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.PostResponseDTO{}, status.FromRepository(err, status.NullCode)
  }

  err = p.getExternalPostData(ctx, post...)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.PostResponseDTO{}, status.ErrExternal(err)
  }

  resp := mapper.ToPostResponseDTO(&post[0])
  return resp, status.Success()
}

func (p *postService) getUserPostVisibility(ctx context.Context, userId, expectedUserId types.Id) (entity.Visibility, error) {
  if userId == expectedUserId {
    return entity.VisibilityOnlyMe, nil
  }

  res, err := p.followClient.IsFollower(ctx, userId, expectedUserId)
  if err != nil {
    return entity.VisibilityUnknown, err
  }

  return sharedUtil.Ternary(res, entity.VisibilityFollower, entity.VisibilityPublic), nil
}

func (p *postService) FindByUserId(ctx context.Context, userId types.Id, pageDTO *sharedDto.PagedElementDTO) (sharedDto.PagedElementResult[dto.PostResponseDTO], status.Object) {
  ctx, span := p.tracer.Start(ctx, "PostService.FindByUserId")
  defer span.End()

  id, err := p.getUserClaims(ctx)
  if err != nil {
    spanUtil.RecordError(err, span)
    return sharedDto.PagedElementResult[dto.PostResponseDTO]{}, status.ErrUnAuthenticated(err)
  }

  // Get claims user relation to user relation
  visibility, err := p.getUserPostVisibility(ctx, id, userId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return sharedDto.PagedElementResult[dto.PostResponseDTO]{}, status.ErrExternal(err)
  }

  repos := p.unit.Repositories()
  result, err := repos.Post().FindByUserId(ctx, false, userId, visibility, pageDTO.ToQueryParam())
  if err != nil {
    spanUtil.RecordError(err, span)
    return sharedDto.PagedElementResult[dto.PostResponseDTO]{}, status.FromRepository(err, status.NullCode)
  }

  err = p.getExternalPostData(ctx, result.Data...)
  if err != nil {
    spanUtil.RecordError(err, span)
    return sharedDto.PagedElementResult[dto.PostResponseDTO]{}, status.ErrExternal(err)
  }

  resp := sharedUtil.CastSliceP(result.Data, mapper.ToPostResponseDTO)
  return sharedDto.NewPagedElementResult2(resp, pageDTO, result.Total), status.Success()
}

func (p *postService) Create(ctx context.Context, createDTO *dto.CreatePostDTO) (types.Id, status.Object) {
  ctx, span := p.tracer.Start(ctx, "PostService.Create")
  defer span.End()

  userId, err := p.getUserClaims(ctx)
  if err != nil {
    spanUtil.RecordError(err, span)
    return types.NullId(), status.ErrUnAuthenticated(err)
  }

  post := createDTO.ToDomain(userId)
  repos := p.unit.Repositories()
  err = repos.Post().Create(ctx, &post)
  if err != nil {
    spanUtil.RecordError(err, span)
    return types.NullId(), status.FromRepositoryExist(err)
  }

  return post.Id, status.Created()
}

func (p *postService) UpdateVisibility(ctx context.Context, id types.Id, newVisibility entity.Visibility) status.Object {
  ctx, span := p.tracer.Start(ctx, "PostService.UpdateVisibility")
  defer span.End()

  repos := p.unit.Repositories()
  err := repos.Post().UpdateVisibility(ctx, id, newVisibility)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }

  return status.Updated()
}

func (p *postService) Edit(ctx context.Context, editDTO *dto.EditPostDTO) status.Object {
  ctx, span := p.tracer.Start(ctx, "PostService.Edit")
  defer span.End()

  userId, err := p.getUserClaims(ctx)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrUnAuthenticated(err)
  }

  post := editDTO.ToDomain(userId)
  repos := p.unit.Repositories()
  err = repos.Post().Edit(ctx, &post)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }
  return status.Updated()
}

func (p *postService) deleteArbitraryPost(ctx context.Context, id types.Id) status.Object {
  span := trace.SpanFromContext(ctx)

  stat := status.Deleted()
  _ = p.unit.DoTx(ctx, func(ctx context.Context, storage uow.PostStorage) error {
    err := storage.Post().Delete(ctx, id)
    if err != nil {
      spanUtil.RecordError(err, span)
      stat = status.FromRepository(err, status.NullCode)
      return err
    }

    // NOTE: Comment client will delete each comments like
    err = p.commentClient.DeletePostsComments(ctx, id)
    if err != nil {
      spanUtil.RecordError(err, span)
      stat = status.ErrExternal(err)
      return err
    }

    err = p.likeClient.DeletePostsLikes(ctx, id)
    if err != nil {
      spanUtil.RecordError(err, span)
      stat = status.ErrExternal(err)
      return err
    }

    return nil
  })

  return stat
}

func (p *postService) Delete(ctx context.Context, id types.Id) status.Object {
  ctx, span := p.tracer.Start(ctx, "PostService.Delete")
  defer span.End()

  // Check if user wanted to remove other user post
  claims, err := sharedJwt.GetUserClaimsFromCtx(ctx)
  if err != nil {
    return status.ErrUnAuthenticated(sharedErr.ErrUnauthenticated)
  }

  // Allow to remove arbitrary post
  if authUtil.ContainsPermission(claims.Roles, constant.POST_PERMISSIONS[constant.POST_DELETE_ARB]) {
    return p.deleteArbitraryPost(ctx, id)
  }

  userId, err := types.IdFromString(claims.UserId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrBadRequest(err)
  }

  stat := status.Deleted()
  _ = p.unit.DoTx(ctx, func(ctx context.Context, storage uow.PostStorage) error {
    // Delete the post
    _, err = storage.Post().DeleteUsers(ctx, userId, id)
    if err != nil {
      spanUtil.RecordError(err, span)
      stat = status.FromRepository(err, status.NullCode)
      return err
    }

    // NOTE: Comment client will delete each comments like
    err = p.commentClient.DeletePostsComments(ctx, id)
    if err != nil {
      spanUtil.RecordError(err, span)
      stat = status.ErrExternal(err)
      return err
    }

    err = p.likeClient.DeletePostsLikes(ctx, id)
    if err != nil {
      spanUtil.RecordError(err, span)
      stat = status.ErrExternal(err)
      return err
    }

    return nil
  })

  return stat
}

func (p *postService) ToggleBookmark(ctx context.Context, postId types.Id) status.Object {
  ctx, span := p.tracer.Start(ctx, "PostService.ToggleBookmark")
  defer span.End()

  userId, err := p.getUserClaims(ctx)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrUnAuthenticated(err)
  }

  repos := p.unit.Repositories()
  err = repos.Post().DelsertBookmark(ctx, userId, postId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }
  return status.Updated()
}

func (p *postService) GetBookmarked(ctx context.Context, userId types.Id, elementDTO *sharedDto.PagedElementDTO) (sharedDto.PagedElementResult[dto.PostResponseDTO], status.Object) {
  ctx, span := p.tracer.Start(ctx, "PostService.GetBookmarked")
  defer span.End()

  err := p.checkPermission(ctx, userId, constant.POST_PERMISSIONS[constant.POST_GET_ARB])
  if err != nil {
    spanUtil.RecordError(err, span)
    return sharedDto.PagedElementResult[dto.PostResponseDTO]{}, status.ErrUnAuthorized(err)
  }

  repos := p.unit.Repositories()
  result, err := repos.Post().GetBookmarked(ctx, userId, elementDTO.ToQueryParam())
  if err != nil {
    spanUtil.RecordError(err, span)
    return sharedDto.PagedElementResult[dto.PostResponseDTO]{}, status.FromRepository(err, status.NullCode)
  }

  resp := sharedUtil.CastSliceP(result.Data, mapper.ToPostResponseDTO)
  return sharedDto.NewPagedElementResult2(resp, elementDTO, result.Total), status.Success()
}

func (p *postService) ClearUsers(ctx context.Context, userId types.Id) status.Object {
  ctx, span := p.tracer.Start(ctx, "PostService.ClearUsers")
  defer span.End()

  err := p.checkPermission(ctx, userId, constant.POST_PERMISSIONS[constant.POST_DELETE_ARB])
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrUnAuthorized(err)
  }

  // Delete posts
  stat := status.Deleted()
  _ = p.unit.DoTx(ctx, func(ctx context.Context, storage uow.PostStorage) error {
    ctx, span := p.tracer.Start(ctx, "UOW.ClearUsers")
    defer span.End()

    // Delete all user's posts
    ids, err := storage.Post().DeleteUsers(ctx, userId)
    if err != nil {
      spanUtil.RecordError(err, span)
      stat = status.FromRepository(err, status.NullCode)
      return err
    }

    // Delete all post's comments
    // NOTE: Each comments like deletions is performed by comment client
    err = p.commentClient.DeletePostsComments(ctx, ids...)
    if err != nil {
      spanUtil.RecordError(err, span)
      stat = status.ErrExternal(err)
      return err
    }

    // Delete posts likes
    err = p.likeClient.DeletePostsLikes(ctx, ids...)
    if err != nil {
      spanUtil.RecordError(err, span)
      stat = status.ErrExternal(err)
      return err
    }

    return nil
  })

  return stat
}
