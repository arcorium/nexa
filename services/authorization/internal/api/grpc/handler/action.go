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

func NewAction(action service.IAction) ActionHandler {
	return ActionHandler{svc: action}
}

type ActionHandler struct {
	proto.UnimplementedActionServerServer

	svc service.IAction
}

func (a *ActionHandler) Register(server *grpc.Server) {
	proto.RegisterActionServerServer(server, a)
}

func (a *ActionHandler) Create(ctx context.Context, input *proto.ActionCreateInput) (*proto.ActionCreateOutput, error) {

	createDTO := mapper.ToActionCreateDTO(input)
	stats := util.ValidateStruct(ctx, &createDTO)
	if stats.IsError() {
		return nil, stats.ToGRPCError()
	}

	id, stats := a.svc.Create(ctx, &createDTO)
	if stats.IsError() {
		return nil, stats.ToGRPCError()
	}
	return &proto.ActionCreateOutput{Id: id.Underlying().String()}, stats.ToGRPCError()
}

func (a *ActionHandler) Update(ctx context.Context, input *proto.ActionUpdateInput) (*emptypb.Empty, error) {
	updateDTO := mapper.ToActionUpdateDTO(input)
	stats := util.ValidateStruct(ctx, &updateDTO)
	if stats.IsError() {
		return nil, stats.ToGRPCError()
	}

	stats = a.svc.Update(ctx, &updateDTO)
	if stats.IsError() {
		return nil, stats.ToGRPCError()
	}
	return &emptypb.Empty{}, stats.ToGRPCError()
}

func (a *ActionHandler) Delete(ctx context.Context, input *proto.ActionDeleteInput) (*emptypb.Empty, error) {
	id := types.IdFromString(input.Id)
	if err := id.Validate(); err != nil {
		stats := status.ErrBadRequest(err)
		return nil, stats.ToGRPCError()
	}

	stats := a.svc.Delete(ctx, id)
	if stats.IsError() {
		return nil, stats.ToGRPCError()
	}
	return &emptypb.Empty{}, stats.ToGRPCError()
}

func (a *ActionHandler) Find(ctx context.Context, input *proto.ActionFindInput) (*proto.ActionResponse, error) {
	id := types.IdFromString(input.Id)
	if err := id.Validate(); err != nil {
		stats := status.ErrBadRequest(err)
		return nil, stats.ToGRPCError()
	}

	responseDTO, stats := a.svc.Find(ctx, id)
	if stats.IsError() {
		return nil, stats.ToGRPCError()
	}
	return mapper.ToActionResponse(&responseDTO), stats.ToGRPCError()
}

func (a *ActionHandler) FindAll(ctx context.Context, input *sharedProto.PagedElementInput) (*proto.ActionFindAllResponse, error) {
	pagedElementDTO := input.ToDTO()

	responseDTOs, stats := a.svc.FindAll(ctx, &pagedElementDTO)
	if stats.IsError() {
		return nil, stats.ToGRPCError()
	}
	return mapper.ToActionResponses(&responseDTOs), stats.ToGRPCError()
}
