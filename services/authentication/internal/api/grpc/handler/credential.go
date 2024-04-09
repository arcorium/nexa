package handler

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"nexa/services/authentication/internal/api/grpc/mapper"
	"nexa/services/authentication/internal/domain/service"
	"nexa/services/authentication/shared/proto"
	"nexa/shared/status"
	"nexa/shared/types"
	"nexa/shared/util"
)

func NewCredential(svc service.ICredential) CredentialHandler {
	return CredentialHandler{
		svc: svc,
	}
}

type CredentialHandler struct {
	proto.UnimplementedCredentialServer
	svc service.ICredential
}

func (c *CredentialHandler) RegisterHandler(server *grpc.Server) {
	proto.RegisterCredentialServer(server, c)
}

func (c *CredentialHandler) Login(ctx context.Context, input *proto.LoginInput) (*proto.LoginOutput, error) {
	dto := mapper.ToLoginDTO(input)
	if stats := util.ValidateStruct(ctx, &dto); stats.IsError() {
		return nil, stats.ToGRPCError()
	}

	accessToken, stats := c.svc.Login(ctx, &dto)
	if stats.IsError() {
		return nil, stats.ToGRPCError()
	}
	return &proto.LoginOutput{AccessToken: accessToken}, nil
}

func (c *CredentialHandler) Register(ctx context.Context, input *proto.RegisterInput) (*emptypb.Empty, error) {
	dto := mapper.ToRegisterDTO(input)
	if stats := util.ValidateStruct(ctx, &dto); stats.IsError() {
		return nil, stats.ToGRPCError()
	}

	stats := c.svc.Register(ctx, &dto)
	if stats.IsError() {
		return nil, stats.ToGRPCError()
	}
	return &emptypb.Empty{}, nil
}

func (c *CredentialHandler) RefreshToken(ctx context.Context, input *proto.RefreshTokenInput) (*proto.RefreshTokenOutput, error) {
	dto := mapper.ToRefreshTokenDTO(input)
	if stats := util.ValidateStruct(ctx, &dto); stats.IsError() {
		return nil, stats.ToGRPCError()
	}

	accessToken, stats := c.svc.RefreshToken(ctx, &dto)
	if stats.IsError() {
		return nil, stats.ToGRPCError()
	}
	return &proto.RefreshTokenOutput{AccessToken: accessToken}, nil
}

func (c *CredentialHandler) GetCredentials(ctx context.Context, _ *emptypb.Empty) (*proto.GetCredentialsOutput, error) {
	// NOTE: User id is placed on context by interceptor, doesn't need to check here
	credentials, stats := c.svc.GetCurrentCredentials(ctx)
	if stats.IsError() {
		return nil, stats.ToGRPCError()
	}
	return &proto.GetCredentialsOutput{Creds: util.CastSlice(credentials, mapper.ToCredentialResponse)}, nil
}

func (c *CredentialHandler) Logout(ctx context.Context, input *proto.LogoutInput) (*emptypb.Empty, error) {
	ids := util.CastSlice2(input.CredIds, types.IdFromString)
	// Validate
	for _, v := range ids {
		if err := v.Validate(); err != nil {
			stats := status.ErrFieldValidation(err)
			return nil, stats.ToGRPCError()
		}
	}

	stats := c.svc.Logout(ctx, ids...)
	if stats.IsError() {
		return nil, stats.ToGRPCError()
	}
	return &emptypb.Empty{}, nil
}

func (c *CredentialHandler) LogoutAll(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	// NOTE: User id is placed on context by interceptor, doesn't need to check here
	stats := c.svc.LogoutAll(ctx)
	if stats.IsError() {
		return nil, stats.ToGRPCError()
	}
	return &emptypb.Empty{}, nil
}
