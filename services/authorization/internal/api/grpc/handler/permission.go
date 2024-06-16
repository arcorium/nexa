package handler

import (
  "context"
  "google.golang.org/grpc"
  "google.golang.org/protobuf/types/known/emptypb"
  "nexa/services/authorization/internal/api/grpc/mapper"
  "nexa/services/authorization/internal/domain/service"
  "nexa/services/authorization/shared/proto"
  sharedProto "nexa/shared/proto"
  "nexa/shared/status"
  "nexa/shared/types"
  "nexa/shared/util"
)

func NewPermission(permission service.IPermission) PermissionHandler {
  return PermissionHandler{svc: permission}
}

type PermissionHandler struct {
  proto.UnimplementedPermissionServiceServer

  svc service.IPermission
}

func (p *PermissionHandler) Register(server *grpc.Server) {
  proto.RegisterPermissionServiceServer(server, p)
}

func (p *PermissionHandler) Create(ctx context.Context, input *proto.PermissionCreateInput) (*proto.PermissionCreateOutput, error) {
  dto := mapper.ToPermissionCreateDTO(input)
  if stats := util.ValidateStruct(ctx, &dto); stats.IsError() {
    return nil, stats.ToGRPCError()
  }

  id, stats := p.svc.Create(ctx, &dto)
  if stats.IsError() {
    return nil, stats.ToGRPCError()
  }
  return &proto.PermissionCreateOutput{Id: id.Underlying().String()}, stats.ToGRPCError()
}

func (p *PermissionHandler) Find(ctx context.Context, input *proto.PermissionFindInput) (*proto.PermissionResponse, error) {
  id := types.IdFromString(input.Id)
  if err := id.Validate(); err != nil {
    stats := status.ErrBadRequest(err)
    return nil, stats.ToGRPCError()
  }

  response, stats := p.svc.Find(ctx, id)
  if stats.IsError() {
    return nil, stats.ToGRPCError()
  }
  return mapper.ToPermissionResponse(&response), stats.ToGRPCError()
}

func (p *PermissionHandler) FindAll(ctx context.Context, input *sharedProto.PagedElementInput) (*proto.PermissionFindAllOutput, error) {
  dto := input.ToDTO()
  result, stats := p.svc.FindAll(ctx, &dto)
  if stats.IsError() {
    return nil, stats.ToGRPCError()
  }
  return mapper.ToPermissionResponses(&result), stats.ToGRPCError()
}

func (p *PermissionHandler) Delete(ctx context.Context, input *proto.PermissionDeleteInput) (*emptypb.Empty, error) {
  id := types.IdFromString(input.Id)
  if err := id.Validate(); err != nil {
    stats := status.ErrBadRequest(err)
    return nil, stats.ToGRPCError()
  }

  stats := p.svc.Delete(ctx, id)
  if stats.IsError() {
    return nil, stats.ToGRPCError()
  }
  return &emptypb.Empty{}, stats.ToGRPCError()
}

func (p *PermissionHandler) CheckUser(ctx context.Context, input *proto.CheckUserInput) (*emptypb.Empty, error) {
  dto := mapper.ToCheckUserPermissionDTO(input)
  if stats := util.ValidateStruct(ctx, &dto); stats.IsError() {
    return nil, stats.ToGRPCError()
  }

  stats := p.svc.CheckUserPermission(ctx, &dto)
  if stats.IsError() {
    return nil, stats.ToGRPCError()
  }
  return &emptypb.Empty{}, stats.ToGRPCError()
}
