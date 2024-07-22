package handler

import (
  "context"
  commentv1 "github.com/arcorium/nexa/proto/gen/go/comment/v1"
  sharedErr "github.com/arcorium/nexa/shared/errors"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  "google.golang.org/protobuf/types/known/emptypb"
  "nexa/services/comment/internal/api/grpc/mapper"
  "nexa/services/comment/internal/domain/service"
  "nexa/services/comment/util"
)

func NewComment(comment service.IComment) CommentHandler {
  return CommentHandler{
    svc:    comment,
    tracer: util.GetTracer(),
  }
}

type CommentHandler struct {
  commentv1.UnimplementedCommentServiceServer
  svc    service.IComment
  tracer trace.Tracer
}

func (c *CommentHandler) Register(server *grpc.Server) {
  commentv1.RegisterCommentServiceServer(server, c)
}

func (c *CommentHandler) Create(ctx context.Context, request *commentv1.CreateCommentRequest) (*commentv1.CreateCommentResponse, error) {
  ctx, span := c.tracer.Start(ctx, "CommentHandler.Create")
  defer span.End()

  createDTO, err := mapper.ToCreateCommentDTO(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  id, stat := c.svc.Create(ctx, &createDTO)
  if stat.IsError() {
    spanUtil.RecordError(err, span)
    return nil, stat.ToGRPCError()
  }

  resp := &commentv1.CreateCommentResponse{
    PostId:    createDTO.PostId.String(),
    CommentId: id.String(),
  }
  return resp, nil
}

func (c *CommentHandler) Edit(ctx context.Context, request *commentv1.EditCommentRequest) (*emptypb.Empty, error) {
  ctx, span := c.tracer.Start(ctx, "CommentHandler.Edit")
  defer span.End()

  editDTO, err := mapper.ToEditCommentDTO(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  stat := c.svc.Edit(ctx, &editDTO)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (c *CommentHandler) Delete(ctx context.Context, request *commentv1.DeleteCommentRequest) (*emptypb.Empty, error) {
  ctx, span := c.tracer.Start(ctx, "CommentHandler.Delete")
  defer span.End()

  id, err := types.IdFromString(request.CommentId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, sharedErr.NewFieldError("comment_id", err).ToGrpcError()
  }

  stat := c.svc.Delete(ctx, id)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (c *CommentHandler) GetPosts(ctx context.Context, request *commentv1.GetPostCommentsRequest) (*commentv1.GetPostCommentsResponse, error) {
  ctx, span := c.tracer.Start(ctx, "CommentHandler.GetPosts")
  defer span.End()

  getDTO, pageDTO, err := mapper.ToGetPostsCommentsDTO(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  result, stat := c.svc.GetPosts(ctx, &getDTO, pageDTO)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  resp := &commentv1.GetPostCommentsResponse{
    Comments: sharedUtil.CastSliceP(result.Data, mapper.ToProtoComments),
    Details:  mapper.ToProtoPagedElementOutput(&result),
  }
  return resp, nil
}

func (c *CommentHandler) GetReplies(ctx context.Context, request *commentv1.GetCommentRepliesRequest) (*commentv1.GetCommentRepliesResponse, error) {
  ctx, span := c.tracer.Start(ctx, "CommentHandler.GetReplies")
  defer span.End()

  getDTO, pageDTO, err := mapper.ToGetCommentRepliesDTO(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  result, stat := c.svc.GetReplies(ctx, &getDTO, pageDTO)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  resp := &commentv1.GetCommentRepliesResponse{
    Replies: sharedUtil.CastSliceP(result.Data, mapper.ToProtoComments),
    Details: mapper.ToProtoPagedElementOutput(&result),
  }
  return resp, nil
}

func (c *CommentHandler) GetCounts(ctx context.Context, request *commentv1.GetCountsRequest) (*commentv1.GetCountsResponse, error) {
  ctx, span := c.tracer.Start(ctx, "CommentHandler.GetCounts")
  defer span.End()

  itemType, err := mapper.ToEntityItemType(request.ItemType)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, sharedErr.NewFieldError("item_type", err).ToGrpcError()
  }

  ids, ierr := sharedUtil.CastSliceErrs(request.ItemIds, types.IdFromString)
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return nil, ierr.ToGRPCError("item_ids")
  }

  counts, stat := c.svc.GetCounts(ctx, itemType, ids...)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  return &commentv1.GetCountsResponse{Total: counts}, nil
}

func (c *CommentHandler) ClearPosts(ctx context.Context, request *commentv1.ClearPostsCommentsRequest) (*emptypb.Empty, error) {
  ctx, span := c.tracer.Start(ctx, "CommentHandler.ClearPosts")
  defer span.End()

  postIds, ierr := sharedUtil.CastSliceErrs(request.PostIds, types.IdFromString)
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return nil, ierr.ToGRPCError("post_ids")
  }

  stat := c.svc.ClearPosts(ctx, postIds...)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (c *CommentHandler) ClearUsers(ctx context.Context, request *commentv1.ClearUserCommentsRequest) (*emptypb.Empty, error) {
  ctx, span := c.tracer.Start(ctx, "CommentHandler.ClearUsers")
  defer span.End()

  userId, err := types.IdFromString(request.UserId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, sharedErr.NewFieldError("comment_id", err).ToGrpcError()
  }

  stat := c.svc.ClearUsers(ctx, userId)
  return nil, stat.ToGRPCErrorWithSpan(span)
}
