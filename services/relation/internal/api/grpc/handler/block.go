package handler

import (
  "context"
  common "github.com/arcorium/nexa/proto/gen/go/common"
  relationv1 "github.com/arcorium/nexa/proto/gen/go/relation/v1"
  sharedDto "github.com/arcorium/nexa/shared/dto"
  sharedErr "github.com/arcorium/nexa/shared/errors"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  emptypb "google.golang.org/protobuf/types/known/emptypb"
  "nexa/services/relation/internal/api/grpc/mapper"
  "nexa/services/relation/internal/domain/dto"
  "nexa/services/relation/internal/domain/service"
  "nexa/services/relation/util"
)

func NewBlock(block service.IBlock) BlockHandler {
  return BlockHandler{
    svc:    block,
    tracer: util.GetTracer(),
  }
}

type BlockHandler struct {
  relationv1.UnimplementedBlockServiceServer

  svc    service.IBlock
  tracer trace.Tracer
}

func (b *BlockHandler) Register(server *grpc.Server) {
  relationv1.RegisterBlockServiceServer(server, b)
}

func (b *BlockHandler) Block(ctx context.Context, request *relationv1.BlockUserRequest) (*emptypb.Empty, error) {
  ctx, span := b.tracer.Start(ctx, "BlockHandler.Block")
  defer span.End()

  userId, err := types.IdFromString(request.TargetUserId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, sharedErr.NewFieldError("target_user_id", err).ToGrpcError()
  }

  stat := b.svc.Block(ctx, userId)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (b *BlockHandler) Unblock(ctx context.Context, request *relationv1.UnblockUserRequest) (*emptypb.Empty, error) {
  ctx, span := b.tracer.Start(ctx, "BlockHandler.Unblock")
  defer span.End()

  userId, err := types.IdFromString(request.TargetUserId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, sharedErr.NewFieldError("target_user_id", err).ToGrpcError()
  }

  stat := b.svc.Unblock(ctx, userId)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (b *BlockHandler) GetBlocked(ctx context.Context, input *common.PagedElementInput) (*relationv1.GetBlockedResponse, error) {
  ctx, span := b.tracer.Start(ctx, "BlockHandler.Unblock")
  defer span.End()

  pageDTO := sharedDto.PagedElementDTO{
    Element: input.Element,
    Page:    input.Page,
  }
  result, stat := b.svc.GetUsers(ctx, pageDTO)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  resp := &relationv1.GetBlockedResponse{
    UserIds: sharedUtil.CastSliceP(result.Data, func(respDTO *dto.BlockResponseDTO) string {
      return respDTO.UserId.String()
    }),
    Details: mapper.ToProtoPagedElementOutput(&result),
  }

  return resp, nil
}

func (b *BlockHandler) GetUsersCount(ctx context.Context, request *relationv1.GetUsersBlockCountRequest) (*relationv1.GetUsersBlockCountResponse, error) {
  ctx, span := b.tracer.Start(ctx, "BlockHandler.ClearUsers")
  defer span.End()

  userId, err := types.IdFromString(request.UserId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, sharedErr.NewFieldError("user_id", err).ToGrpcError()
  }

  count, stat := b.svc.GetUsersCount(ctx, userId)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  resp := &relationv1.GetUsersBlockCountResponse{
    Counts: mapper.ToProtoBlockCount(&count),
  }
  return resp, nil
}

func (b *BlockHandler) ClearUsers(ctx context.Context, request *relationv1.ClearUsersBlockRequest) (*emptypb.Empty, error) {
  ctx, span := b.tracer.Start(ctx, "BlockHandler.ClearUsers")
  defer span.End()

  userId, err := types.IdFromString(request.UserId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, sharedErr.NewFieldError("user_id", err).ToGrpcError()
  }

  stat := b.svc.ClearUsers(ctx, userId)
  return nil, stat.ToGRPCErrorWithSpan(span)
}
