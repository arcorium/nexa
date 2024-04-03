package handler

import (
	"context"
	"google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	"nexa/services/authorization/internal/api/grpc/mapper"
	"nexa/services/authorization/internal/domain/service"
	"nexa/services/authorization/shared/proto"
	sharedProto "nexa/shared/proto"
	"nexa/shared/status"
	"nexa/shared/types"
	"nexa/shared/util"
)

func NewRole(role service.IRole) RoleHandler {
	return RoleHandler{svc: role}
}

type RoleHandler struct {
	proto.UnimplementedRoleServiceServer
	svc service.IRole
}

func (r *RoleHandler) Register(server *grpc.Server) {
	proto.RegisterRoleServiceServer(server, r)
}

func (r *RoleHandler) Create(ctx context.Context, input *proto.RoleCreateInput) (*proto.RoleCreateOutput, error) {
	createDTO := mapper.ToRoleCreateDTO(input)
	if stats := util.ValidateStruct(ctx, &createDTO); stats.IsError() {
		return nil, stats.ToGRPCError()
	}

	id, stats := r.svc.Create(ctx, &createDTO)
	if stats.IsError() {
		return nil, stats.ToGRPCError()
	}
	return &proto.RoleCreateOutput{Id: id.Underlying().String()}, stats.ToGRPCError()
}

func (r *RoleHandler) Update(ctx context.Context, input *proto.RoleUpdateInput) (*emptypb.Empty, error) {
	updateDTO := mapper.ToRoleUpdateDTO(input)
	if stats := util.ValidateStruct(ctx, &updateDTO); stats.IsError() {
		return nil, stats.ToGRPCError()
	}

	stats := r.svc.Update(ctx, &updateDTO)
	if stats.IsError() {
		return nil, stats.ToGRPCError()
	}
	return &emptypb.Empty{}, stats.ToGRPCError()
}

func (r *RoleHandler) Delete(ctx context.Context, input *proto.RoleDeleteInput) (*emptypb.Empty, error) {
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

func (r *RoleHandler) Find(ctx context.Context, input *proto.RoleFindInput) (*proto.RoleFindOutput, error) {
	ids := util.CastSlice2(input.Ids, types.IdFromString)
	for _, id := range ids {
		if err := id.Validate(); err != nil {
			stats := status.ErrBadRequest(err)
			return nil, stats.ToGRPCError()
		}
	}

	responseDTOS, stats := r.svc.Find(ctx, ids...)
	if stats.IsError() {
		return nil, stats.ToGRPCError()
	}

	return &proto.RoleFindOutput{Roles: util.CastSlice(responseDTOS, mapper.ToRoleResponse)}, stats.ToGRPCError()
}

func (r *RoleHandler) FindAll(ctx context.Context, input *sharedProto.PagedElementInput) (*proto.RoleFindAllOutput, error) {
	pagedElementDTO := input.ToDTO()
	result, stats := r.svc.FindAll(ctx, &pagedElementDTO)
	if stats.IsError() {
		return nil, stats.ToGRPCError()
	}

	return mapper.ToRoleResponses(&result), stats.ToGRPCError()
}

func (r *RoleHandler) AddPermissions(ctx context.Context, input *proto.RoleAddPermissionsInput) (*emptypb.Empty, error) {
	dto := mapper.ToAddPermissionsDTO(input)
	if stats := util.ValidateStruct(ctx, &dto); stats.IsError() {
		return nil, stats.ToGRPCError()
	}

	stats := r.svc.AddPermissions(ctx, &dto)
	if stats.IsError() {
		return nil, stats.ToGRPCError()
	}
	return &emptypb.Empty{}, stats.ToGRPCError()
}

func (r *RoleHandler) RemovePermissions(ctx context.Context, input *proto.RoleRemovePermissionsInput) (*emptypb.Empty, error) {
	dto := mapper.ToRemovePermissionsDTO(input)
	if stats := util.ValidateStruct(ctx, &dto); stats.IsError() {
		return nil, stats.ToGRPCError()
	}

	stats := r.svc.RemovePermissions(ctx, &dto)
	if stats.IsError() {
		return nil, stats.ToGRPCError()
	}
	return &emptypb.Empty{}, stats.ToGRPCError()
}

func (r *RoleHandler) AddUsers(ctx context.Context, input *proto.RoleAddUserInput) (*emptypb.Empty, error) {
	dto := mapper.ToAddUsersDTO(input)
	if stats := util.ValidateStruct(ctx, &dto); stats.IsError() {
		return nil, stats.ToGRPCError()
	}

	stats := r.svc.AddUsers(ctx, &dto)
	if stats.IsError() {
		return nil, stats.ToGRPCError()
	}
	return &emptypb.Empty{}, stats.ToGRPCError()
}

func (r *RoleHandler) RemoveUsers(ctx context.Context, input *proto.RoleRemoveUserInput) (*emptypb.Empty, error) {
	dto := mapper.ToRemoveUsersDTO(input)
	if stats := util.ValidateStruct(ctx, &dto); stats.IsError() {
		return nil, stats.ToGRPCError()
	}

	stats := r.svc.RemoveUsers(ctx, &dto)
	if stats.IsError() {
		return nil, stats.ToGRPCError()
	}
	return &emptypb.Empty{}, stats.ToGRPCError()
}
