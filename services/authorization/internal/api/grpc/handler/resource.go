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

func NewResource(resource service.IResource) ResourceHandler {
	return ResourceHandler{svc: resource}
}

type ResourceHandler struct {
	proto.UnimplementedResourceServiceServer

	svc service.IResource
}

func (r *ResourceHandler) Register(server *grpc.Server) {
	proto.RegisterResourceServiceServer(server, r)
}

func (r *ResourceHandler) Create(ctx context.Context, input *proto.ResourceCreateInput) (*proto.ResourceCreateOutput, error) {
	createDTO := mapper.ToResourceCreateDTO(input)
	if stats := util.ValidateStruct(ctx, &createDTO); stats.IsError() {
		return nil, stats.ToGRPCError()
	}

	id, stats := r.svc.Create(ctx, &createDTO)
	if stats.IsError() {
		return nil, stats.ToGRPCError()
	}
	return &proto.ResourceCreateOutput{Id: id.Underlying().String()}, stats.ToGRPCError()
}

func (r *ResourceHandler) Update(ctx context.Context, input *proto.ResourceUpdateInput) (*emptypb.Empty, error) {
	updateDTO := mapper.ToResourceUpdateDTO(input)
	if stats := util.ValidateStruct(ctx, &updateDTO); stats.IsError() {
		return nil, stats.ToGRPCError()
	}

	stats := r.svc.Update(ctx, &updateDTO)
	if stats.IsError() {
		return nil, stats.ToGRPCError()
	}
	return &emptypb.Empty{}, stats.ToGRPCError()
}

func (r *ResourceHandler) Delete(ctx context.Context, input *proto.ResourceDeleteInput) (*emptypb.Empty, error) {
	id := types.IdFromString(input.Id)
	if err := id.Validate(); err != nil {
		stats := status.ErrBadRequest(err)
		return nil, stats.ToGRPCError()
	}

	stats := r.svc.Delete(ctx, id)
	if stats.IsError() {
		return nil, stats.ToGRPCError()
	}
	return &emptypb.Empty{}, stats.ToGRPCError()
}

func (r *ResourceHandler) Find(ctx context.Context, input *proto.ResourceFindInput) (*proto.ResourceResponse, error) {
	id := types.IdFromString(input.Id)
	if err := id.Validate(); err != nil {
		stats := status.ErrBadRequest(err)
		return nil, stats.ToGRPCError()
	}

	response, stats := r.svc.Find(ctx, id)
	if stats.IsError() {
		return nil, stats.ToGRPCError()
	}
	return mapper.ToResourceResponse(&response), stats.ToGRPCError()
}

func (r *ResourceHandler) FindAll(ctx context.Context, input *sharedProto.PagedElementInput) (*proto.ResourceFindAllResponse, error) {
	pagedElementDTO := input.ToDTO()

	result, stats := r.svc.FindAll(ctx, &pagedElementDTO)
	if stats.IsError() {
		return nil, stats.ToGRPCError()
	}
	return mapper.ToResourceResponses(&result), stats.ToGRPCError()
}
