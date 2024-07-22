package handler

import (
  "context"
  common "github.com/arcorium/nexa/proto/gen/go/common"
  postv1 "github.com/arcorium/nexa/proto/gen/go/post/v1"
  sharedDto "github.com/arcorium/nexa/shared/dto"
  sharedErr "github.com/arcorium/nexa/shared/errors"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  emptypb "google.golang.org/protobuf/types/known/emptypb"
  "nexa/services/post/internal/api/grpc/mapper"
  "nexa/services/post/internal/domain/service"
  "nexa/services/post/util"
)

func NewPost(post service.IPost) PostHandler {
  return PostHandler{
    svc:    post,
    tracer: util.GetTracer(),
  }
}

type PostHandler struct {
  postv1.UnimplementedPostServiceServer
  svc    service.IPost
  tracer trace.Tracer
}

func (p *PostHandler) Register(server *grpc.Server) {
  postv1.RegisterPostServiceServer(server, p)
}

func (p *PostHandler) Find(ctx context.Context, input *common.PagedElementInput) (*postv1.FindPostResponse, error) {
  ctx, span := p.tracer.Start(ctx, "PostHandler.Find")
  defer span.End()

  pageDto := sharedDto.PagedElementDTO{
    Element: input.Element,
    Page:    input.Page,
  }

  result, stat := p.svc.GetAll(ctx, &pageDto)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  posts := sharedUtil.CastSliceP(result.Data, mapper.ToProtoPost)
  return &postv1.FindPostResponse{Posts: posts, Details: mapper.ToProtoPagedOutput(&result)}, nil
}

func (p *PostHandler) FindEdited(ctx context.Context, request *postv1.FindEditedPostRequest) (*postv1.FindEditedPostResponse, error) {
  ctx, span := p.tracer.Start(ctx, "PostHandler.FindEdited")
  defer span.End()

  postId, err := types.IdFromString(request.PostId)
  if err != nil {
    cerr := sharedErr.NewFieldError("post_id", err)
    spanUtil.RecordError(cerr, span)
    return nil, cerr.ToGrpcError()
  }

  result, stat := p.svc.GetEdited(ctx, postId)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  posts := mapper.ToProtoEditedPost(&result)
  return &postv1.FindEditedPostResponse{Post: posts}, nil
}

func (p *PostHandler) FindById(ctx context.Context, request *postv1.FindPostByIdRequest) (*postv1.FindPostByIdResponse, error) {
  ctx, span := p.tracer.Start(ctx, "PostHandler.FindById")
  defer span.End()

  postId, err := types.IdFromString(request.PostId)
  if err != nil {
    cerr := sharedErr.NewFieldError("post_id", err)
    spanUtil.RecordError(cerr, span)
    return nil, cerr.ToGrpcError()
  }

  post, stat := p.svc.FindById(ctx, postId)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  resp := mapper.ToProtoPost(&post)
  return &postv1.FindPostByIdResponse{Post: resp}, nil
}

func (p *PostHandler) FindUsers(ctx context.Context, request *postv1.FindUserPostRequest) (*postv1.FindUserPostResponse, error) {
  ctx, span := p.tracer.Start(ctx, "PostHandler.FindUsers")
  defer span.End()

  pageDto := sharedDto.PagedElementDTO{
    Element: request.Details.Element,
    Page:    request.Details.Page,
  }

  userId, err := types.IdFromString(request.UserId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, sharedErr.NewFieldError("user_id", err).ToGrpcError()
  }

  result, stat := p.svc.FindByUserId(ctx, userId, &pageDto)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  resp := sharedUtil.CastSliceP(result.Data, mapper.ToProtoPost)
  return &postv1.FindUserPostResponse{Posts: resp, Details: mapper.ToProtoPagedOutput(&result)}, nil
}

func (p *PostHandler) Create(ctx context.Context, request *postv1.CreatePostRequest) (*postv1.CreatePostResponse, error) {
  ctx, span := p.tracer.Start(ctx, "PostHandler.Create")
  defer span.End()

  postDTO, err := mapper.ToCreatePostDTO(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  id, stat := p.svc.Create(ctx, &postDTO)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  return &postv1.CreatePostResponse{PostId: id.String()}, nil
}

func (p *PostHandler) UpdateVisibility(ctx context.Context, request *postv1.UpdatePostVisibilityRequest) (*emptypb.Empty, error) {
  ctx, span := p.tracer.Start(ctx, "PostHandler.UpdateVisibility")
  defer span.End()

  postId, err := types.IdFromString(request.PostId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, sharedErr.NewFieldError("post_id", err).ToGrpcError()
  }

  visibility, err := mapper.ToEntityVisibility(request.NewVisibility)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, sharedErr.NewFieldError("new_visibility", err).ToGrpcError()
  }

  stat := p.svc.UpdateVisibility(ctx, postId, visibility)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (p *PostHandler) Edit(ctx context.Context, request *postv1.EditPostRequest) (*emptypb.Empty, error) {
  ctx, span := p.tracer.Start(ctx, "PostHandler.Edit")
  defer span.End()

  postDTO, err := mapper.ToEditPostDTO(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  stat := p.svc.Edit(ctx, &postDTO)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (p *PostHandler) Bookmark(ctx context.Context, request *postv1.BookmarkPostRequest) (*emptypb.Empty, error) {
  ctx, span := p.tracer.Start(ctx, "PostHandler.Bookmark")
  defer span.End()

  postId, err := types.IdFromString(request.PostId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, sharedErr.NewFieldError("post_id", err).ToGrpcError()
  }

  stat := p.svc.ToggleBookmark(ctx, postId)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (p *PostHandler) GetBookmarked(ctx context.Context, request *postv1.GetBookmarkedPostRequest) (*postv1.GetBookmarkedPostResponse, error) {
  ctx, span := p.tracer.Start(ctx, "PostHandler.GetBookmarked")
  defer span.End()

  userId, err := types.IdFromString(request.UserId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, sharedErr.NewFieldError("user_id", err).ToGrpcError()
  }

  pageDTO := sharedDto.PagedElementDTO{
    Element: request.Details.Element,
    Page:    request.Details.Page,
  }

  result, stat := p.svc.GetBookmarked(ctx, userId, &pageDTO)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  resp := &postv1.GetBookmarkedPostResponse{
    Posts:   sharedUtil.CastSliceP(result.Data, mapper.ToProtoPost),
    Details: mapper.ToProtoPagedOutput(&result),
  }
  return resp, nil
}

func (p *PostHandler) Delete(ctx context.Context, request *postv1.DeletePostRequest) (*emptypb.Empty, error) {
  ctx, span := p.tracer.Start(ctx, "PostHandler.Delete")
  defer span.End()

  postId, err := types.IdFromString(request.PostId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, sharedErr.NewFieldError("post_id", err).ToGrpcError()
  }

  stat := p.svc.Delete(ctx, postId)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (p *PostHandler) ClearUsers(ctx context.Context, request *postv1.ClearUserPostsRequest) (*emptypb.Empty, error) {
  ctx, span := p.tracer.Start(ctx, "PostHandler.ClearUsers")
  defer span.End()

  userId, err := types.IdFromString(request.UserId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, sharedErr.NewFieldError("user_id", err).ToGrpcError()
  }

  stat := p.svc.ClearUsers(ctx, userId)
  return nil, stat.ToGRPCErrorWithSpan(span)
}
