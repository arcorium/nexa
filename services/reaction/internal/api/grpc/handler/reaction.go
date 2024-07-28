package handler

import (
  "context"
  reactionv1 "github.com/arcorium/nexa/proto/gen/go/reaction/v1"
  sharedErr "github.com/arcorium/nexa/shared/errors"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  "google.golang.org/protobuf/types/known/emptypb"
  "nexa/services/reaction/internal/api/grpc/mapper"
  "nexa/services/reaction/internal/domain/service"
  "nexa/services/reaction/util"
)

func NewReaction(reaction service.IReaction) ReactionHandler {
  return ReactionHandler{
    svc:    reaction,
    tracer: util.GetTracer(),
  }
}

type ReactionHandler struct {
  reactionv1.UnimplementedReactionServiceServer
  svc    service.IReaction
  tracer trace.Tracer
}

func (r *ReactionHandler) Register(server *grpc.Server) {
  reactionv1.RegisterReactionServiceServer(server, r)
}

func (r *ReactionHandler) Like(ctx context.Context, request *reactionv1.LikeRequest) (*emptypb.Empty, error) {
  ctx, span := r.tracer.Start(ctx, "ReactionHandler.Like")
  defer span.End()

  itemType, err := mapper.ToEntityItemType(request.ItemType)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, sharedErr.NewFieldError("item_type", err).ToGrpcError()
  }

  itemId, err := types.IdFromString(request.ItemId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, sharedErr.NewFieldError("id", err).ToGrpcError()
  }

  stat := r.svc.Like(ctx, itemType, itemId)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (r *ReactionHandler) Dislike(ctx context.Context, request *reactionv1.DislikeRequest) (*emptypb.Empty, error) {
  ctx, span := r.tracer.Start(ctx, "ReactionHandler.Dislike")
  defer span.End()

  itemType, err := mapper.ToEntityItemType(request.ItemType)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, sharedErr.NewFieldError("item_type", err).ToGrpcError()
  }

  itemId, err := types.IdFromString(request.ItemId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, sharedErr.NewFieldError("id", err).ToGrpcError()
  }

  stat := r.svc.Dislike(ctx, itemType, itemId)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (r *ReactionHandler) GetItems(ctx context.Context, request *reactionv1.GetItemReactionsRequest) (*reactionv1.GetItemReactionsResponse, error) {
  ctx, span := r.tracer.Start(ctx, "ReactionHandler.GetItems")
  defer span.End()

  itemType, err := mapper.ToEntityItemType(request.ItemType)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, sharedErr.NewFieldError("item_type", err).ToGrpcError()
  }

  itemId, err := types.IdFromString(request.ItemId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, sharedErr.NewFieldError("id", err).ToGrpcError()
  }

  pageDTO := mapper.ToPagedElementDTO(request.Details)
  result, stat := r.svc.GetItemsReactions(ctx, itemType, itemId, pageDTO)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCErrorWithSpan(span)
  }

  resp := sharedUtil.CastSliceP(result.Data, mapper.ToProtoReaction)
  return &reactionv1.GetItemReactionsResponse{
    ItemType:  request.ItemType,
    ItemId:    request.ItemId,
    Reactions: resp,
    Details:   mapper.ToCommonPagedOutput(&result),
  }, nil
}

func (r *ReactionHandler) GetCount(ctx context.Context, request *reactionv1.GetCountRequest) (*reactionv1.GetCountResponse, error) {
  ctx, span := r.tracer.Start(ctx, "ReactionHandler.GetCount")
  defer span.End()

  itemType, err := mapper.ToEntityItemType(request.ItemType)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, sharedErr.NewFieldError("item_type", err).ToGrpcError()
  }

  itemIds, ierr := sharedUtil.CastSliceErrs(request.ItemIds, types.IdFromString)
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return nil, ierr.ToGRPCError("item_ids")
  }

  counts, stat := r.svc.GetCounts(ctx, itemType, itemIds...)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  resp := sharedUtil.CastSliceP(counts, mapper.ToProtoCount)
  return &reactionv1.GetCountResponse{Counts: resp}, nil
}

func (r *ReactionHandler) DeleteItems(ctx context.Context, request *reactionv1.DeleteItemsReactionsRequest) (*emptypb.Empty, error) {
  ctx, span := r.tracer.Start(ctx, "ReactionHandler.Delete")
  defer span.End()

  itemType, err := mapper.ToEntityItemType(request.ItemType)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, sharedErr.NewFieldError("item_type", err).ToGrpcError()
  }

  itemIds, ierr := sharedUtil.CastSliceErrs(request.ItemIds, types.IdFromString)
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return nil, ierr.ToGRPCError("item_ids")
  }

  stat := r.svc.Delete(ctx, itemType, itemIds...)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (r *ReactionHandler) ClearUsers(ctx context.Context, request *reactionv1.ClearUsersReactionsRequest) (*emptypb.Empty, error) {
  ctx, span := r.tracer.Start(ctx, "ReactionHandler.ClearUsers")
  defer span.End()

  userId, err := types.IdFromString(request.UserId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, sharedErr.NewFieldError("id", err).ToGrpcError()
  }

  stat := r.svc.ClearUsers(ctx, userId)
  return nil, stat.ToGRPCErrorWithSpan(span)
}
