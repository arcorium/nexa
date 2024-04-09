package handler

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"nexa/services/authentication/internal/api/grpc/mapper"
	"nexa/services/authentication/internal/domain/service"
	"nexa/services/authentication/shared/proto"
	sharedProto "nexa/shared/proto"
	"nexa/shared/status"
	"nexa/shared/types"
	"nexa/shared/util"
)

func NewToken(token service.IToken) TokenHandler {
	return TokenHandler{svc: token}
}

type TokenHandler struct {
	proto.UnimplementedTokenServer
	svc service.IToken
}

func (t *TokenHandler) RegisterHandler(server *grpc.Server) {
	proto.RegisterTokenServer(server, t)
}

func (t *TokenHandler) Request(ctx context.Context, input *proto.RequestInput) (*proto.RequestOutput, error) {
	dto := mapper.ToTokenRequestDTO(input)
	if stats := util.ValidateStruct(ctx, &dto); stats.IsError() {
		return nil, stats.ToGRPCError()
	}

	result, stats := t.svc.Request(ctx, &dto)
	if stats.IsError() {
		return nil, stats.ToGRPCError()
	}
	return &proto.RequestOutput{Token: result.Token}, nil
}

func (t *TokenHandler) Verify(ctx context.Context, input *proto.VerifyInput) (*emptypb.Empty, error) {
	dto := mapper.ToTokenVerifyDTO(input)
	if stats := util.ValidateStruct(ctx, &dto); stats.IsError() {
		return nil, stats.ToGRPCError()
	}

	stats := t.svc.Verify(ctx, &dto)
	if stats.IsError() {
		return nil, stats.ToGRPCError()
	}
	return &emptypb.Empty{}, nil
}

func (t *TokenHandler) AddUsage(ctx context.Context, input *proto.AddUsageInput) (*proto.AddUsageOutput, error) {
	dto := mapper.ToTokenAddUsageDTO(input)
	if stats := util.ValidateStruct(ctx, &dto); stats.IsError() {
		return nil, stats.ToGRPCError()
	}

	id, stats := t.svc.AddUsage(ctx, &dto)
	if stats.IsError() {
		return nil, stats.ToGRPCError()
	}
	return &proto.AddUsageOutput{Id: id.Underlying().String()}, nil
}

func (t *TokenHandler) RemoveUsage(ctx context.Context, input *proto.RemoveUsageInput) (*emptypb.Empty, error) {
	id := types.IdFromString(input.Id)
	if err := id.Validate(); err != nil {
		stats := status.ErrFieldValidation(err)
		return nil, stats.ToGRPCError()
	}

	stats := t.svc.RemoveUsage(ctx, id)
	if stats.IsError() {
		return nil, stats.ToGRPCError()
	}
	return &emptypb.Empty{}, nil
}

func (t *TokenHandler) UpdateUsage(ctx context.Context, input *proto.UpdateUsageInput) (*emptypb.Empty, error) {
	dto := mapper.ToTokenUpdateUsageDTO(input)
	if stats := util.ValidateStruct(ctx, &dto); stats.IsError() {
		return nil, stats.ToGRPCError()
	}

	stats := t.svc.UpdateUsage(ctx, &dto)
	if stats.IsError() {
		return nil, stats.ToGRPCError()
	}
	return &emptypb.Empty{}, nil
}

func (t *TokenHandler) FindUsage(ctx context.Context, input *proto.FindUsageInput) (*proto.TokenUsageResponse, error) {
	id := types.IdFromString(input.Id)
	if err := id.Validate(); err != nil {
		stats := status.ErrFieldValidation(err)
		return nil, stats.ToGRPCError()
	}

	usage, stats := t.svc.FindUsage(ctx, id)
	if stats.IsError() {
		return nil, stats.ToGRPCError()
	}
	return mapper.ToTokenUsageResponse(&usage), nil
}

func (t *TokenHandler) FindAllUsages(ctx context.Context, input *sharedProto.PagedElementInput) (*proto.FindAllUsagesOutput, error) {
	dto := input.ToDTO()

	result, stats := t.svc.FindAllUsages(ctx, &dto)
	if stats.IsError() {
		return nil, stats.ToGRPCError()
	}
	return mapper.ToTokenUsageResponses(&result), nil
}
