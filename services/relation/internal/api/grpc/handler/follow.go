package handler

import (
  "context"
  relationv1 "github.com/arcorium/nexa/proto/gen/go/relation/v1"
  sharedErr "github.com/arcorium/nexa/shared/errors"
  "github.com/arcorium/nexa/shared/jwt"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  "google.golang.org/protobuf/types/known/emptypb"
  "nexa/services/relation/internal/api/grpc/mapper"
  "nexa/services/relation/internal/domain/dto"
  "nexa/services/relation/internal/domain/service"
  "nexa/services/relation/util"
)

func NewFollow(follow service.IFollow) FollowHandler {
  return FollowHandler{
    svc:    follow,
    tracer: util.GetTracer(),
  }
}

type FollowHandler struct {
  relationv1.UnimplementedFollowServiceServer
  svc    service.IFollow
  tracer trace.Tracer
}

func (f *FollowHandler) Register(server *grpc.Server) {
  relationv1.RegisterFollowServiceServer(server, f)
}

func (f *FollowHandler) Follow(ctx context.Context, request *relationv1.FollowRequest) (*emptypb.Empty, error) {
  ctx, span := f.tracer.Start(ctx, "FollowHandler.Follow")
  defer span.End()

  userIds, ierr := sharedUtil.CastSliceErrs(request.UserIds, types.IdFromString)
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return nil, ierr.ToGRPCError("user_ids")
  }

  stat := f.svc.Follow(ctx, userIds...)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (f *FollowHandler) Unfollow(ctx context.Context, request *relationv1.UnfollowRequest) (*emptypb.Empty, error) {
  ctx, span := f.tracer.Start(ctx, "FollowHandler.Unfollow")
  defer span.End()

  userIds, ierr := sharedUtil.CastSliceErrs(request.UserIds, types.IdFromString)
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return nil, ierr.ToGRPCError("user_ids")
  }

  stat := f.svc.Unfollow(ctx, userIds...)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (f *FollowHandler) GetFollowers(ctx context.Context, request *relationv1.GetUserFollowersRequest) (*relationv1.GetUserFollowersResponse, error) {
  ctx, span := f.tracer.Start(ctx, "FollowHandler.GetFollowers")
  defer span.End()

  claims := types.Must(jwt.GetUserClaimsFromCtx(ctx))
  id := types.NewNullable(request.UserId)
  userId, err := types.IdFromString(id.ValueOr(claims.UserId))
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, sharedErr.NewFieldError("user_id", err).ToGrpcError()
  }

  pageDTO := mapper.ToPagedElementDTO(request.Details)

  result, stat := f.svc.GetFollowers(ctx, userId, pageDTO)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  resp := &relationv1.GetUserFollowersResponse{
    UserIds: sharedUtil.CastSliceP(result.Data, func(resp *dto.FollowResponseDTO) string {
      return resp.UserId.String()
    }),
    Details: mapper.ToProtoPagedElementOutput(&result),
  }
  return resp, nil
}

func (f *FollowHandler) GetFollowees(ctx context.Context, request *relationv1.GetUserFolloweesRequest) (*relationv1.GetUserFolloweesResponse, error) {
  ctx, span := f.tracer.Start(ctx, "FollowHandler.GetFollowees")
  defer span.End()

  claims := types.Must(jwt.GetUserClaimsFromCtx(ctx))
  id := types.NewNullable(request.UserId)
  userId, err := types.IdFromString(id.ValueOr(claims.UserId))
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, sharedErr.NewFieldError("user_id", err).ToGrpcError()
  }

  pageDTO := mapper.ToPagedElementDTO(request.Details)

  result, stat := f.svc.GetFollowings(ctx, userId, pageDTO)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  resp := &relationv1.GetUserFolloweesResponse{
    UserIds: sharedUtil.CastSliceP(result.Data, func(resp *dto.FollowResponseDTO) string {
      return resp.UserId.String()
    }),
    Details: mapper.ToProtoPagedElementOutput(&result),
  }
  return resp, nil
}

func (f *FollowHandler) GetRelation(ctx context.Context, request *relationv1.GetRelationRequest) (*relationv1.GetRelationResponse, error) {
  ctx, span := f.tracer.Start(ctx, "FollowHandler.GetFollowStatus")
  defer span.End()

  claims := types.Must(jwt.GetUserClaimsFromCtx(ctx))
  id := types.NewNullable(request.UserId)
  userId, err := types.IdFromString(id.ValueOr(claims.UserId))
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, sharedErr.NewFieldError("user_id", err).ToGrpcError()
  }

  opponentIds, ierr := sharedUtil.CastSliceErrs(request.OpponentUserIds, types.IdFromString)
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return nil, ierr.ToGRPCError("opponent_user_ids")
  }

  statuses, stat := f.svc.GetStatus(ctx, userId, opponentIds...)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  resp := &relationv1.GetRelationResponse{
    Status: sharedUtil.CastSlice(statuses, mapper.ToProtoFollowStatus),
  }
  return resp, nil
}

func (f *FollowHandler) GetUsersCount(ctx context.Context, request *relationv1.GetUsersFollowCountRequest) (*relationv1.GetUsersFollowCountResponse, error) {
  ctx, span := f.tracer.Start(ctx, "FollowHandler.GetUsersCount")
  defer span.End()

  // Nil means to get the user itself
  if len(request.UserIds) == 0 {
    // Get id from claims
    claims := types.Must(jwt.GetUserClaimsFromCtx(ctx))
    request.UserIds = append(request.UserIds, claims.UserId)
  }
  userIds, ierr := sharedUtil.CastSliceErrs(request.UserIds, types.IdFromString)
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return nil, ierr.ToGRPCError("user_ids")
  }

  counts, stat := f.svc.GetUsersCount(ctx, userIds...)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }
  resp := &relationv1.GetUsersFollowCountResponse{
    Counts: sharedUtil.CastSliceP(counts, mapper.ToProtoFollowCount),
  }
  return resp, nil
}

func (f *FollowHandler) ClearUsers(ctx context.Context, request *relationv1.ClearUsersRequest) (*emptypb.Empty, error) {
  ctx, span := f.tracer.Start(ctx, "FollowHandler.ClearUsers")
  defer span.End()

  claims := types.Must(jwt.GetUserClaimsFromCtx(ctx))
  id := types.NewNullable(request.UserId)
  userId, err := types.IdFromString(id.ValueOr(claims.UserId))
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, sharedErr.NewFieldError("user_id", err).ToGrpcError()
  }

  stat := f.svc.ClearUsers(ctx, userId)
  return nil, stat.ToGRPCErrorWithSpan(span)
}
