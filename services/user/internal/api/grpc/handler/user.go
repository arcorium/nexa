package handler

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"nexa/services/user/internal/api/grpc/mapper"
	"nexa/services/user/internal/domain/service"
	"nexa/services/user/shared/proto"
	"nexa/shared/types"
	"nexa/shared/util"
)

func NewUserHandler(user service.IUser) UserHandler {
	return UserHandler{userService: user}
}

type UserHandler struct {
	proto.UnimplementedUserServiceServer

	userService service.IUser
}

func (u *UserHandler) Register(server *grpc.Server) {
	proto.RegisterUserServiceServer(server, u)
}

func (u *UserHandler) Create(ctx context.Context, request *proto.CreateUserRequest) (*emptypb.Empty, error) {
	dtoInput := mapper.ToDTOCreateInput(request)

	// DTO Validation
	stats := util.ValidateStruct(ctx, &dtoInput)
	if stats.IsError() {
		return nil, stats.ToGRPCError()
	}

	stats = u.userService.Create(ctx, &dtoInput)
	return &emptypb.Empty{}, stats.ToGRPCError()
}

func (u *UserHandler) Update(ctx context.Context, request *proto.UpdateUserRequest) (*emptypb.Empty, error) {
	dtoInput := mapper.ToDTOUserUpdateInput(request)
	// TODO: Get user id from access token claims from ctx

	err := util.GetValidator().StructCtx(ctx, &dtoInput)
	if err != nil {
		return nil, err
	}

	stats := u.userService.Update(ctx, &dtoInput)
	return &emptypb.Empty{}, stats.Error
}

func (u *UserHandler) UpdateVerified(ctx context.Context, request *proto.UpdateUserVerifiedRequest) (*emptypb.Empty, error) {
	id := types.IdFromString(request.Id)
	// TODO: Get user id from access token claims from ctx

	err := id.Validate()
	if err != nil {
		return nil, err
	}

	stats := u.userService.UpdateVerified(ctx, id)
	return &emptypb.Empty{}, stats.Error
}

func (u *UserHandler) UpdatePassword(ctx context.Context, request *proto.UpdateUserPasswordRequest) (*emptypb.Empty, error) {
	dtoInput := mapper.ToDTOUserUpdatePasswordInput(request)
	// TODO: Get user id from access token claims from ctx

	err := util.GetValidator().Struct(&dtoInput)
	if err != nil {
		return nil, err
	}

	stats := u.userService.UpdatePassword(ctx, &dtoInput)
	return &emptypb.Empty{}, stats.Error
}

func (u *UserHandler) ResetPassword(ctx context.Context, request *proto.ResetUserPasswordRequest) (*emptypb.Empty, error) {
	dtoInput := mapper.ToDTOUserResetPasswordInput(request)

	err := util.GetValidator().StructCtx(ctx, &dtoInput)
	if err != nil {
		return nil, err
	}

	stats := u.userService.ResetPassword(ctx, &dtoInput)
	return &emptypb.Empty{}, stats.Error
}

func (u *UserHandler) FindUserByIds(request *proto.FindUsersByIdsRequest, server proto.UserService_FindUserByIdsServer) error {
	ids, err := util.CastSliceErr(request.Ids, func(from *string) (types.Id, error) {
		id := types.IdFromString(*from)
		return id, id.Validate()
	})
	if err != nil {
		return err
	}

	users, stats := u.userService.FindByIds(server.Context(), ids)
	if stats.IsError() {
		return stats.Error
	}

	for _, val := range users {
		response := mapper.ToProtoUserResponse(&val)
		if err := server.Send(&response); err != nil {
			return err
		}
	}
	return nil
}

func (u *UserHandler) FindUserByEmail(request *proto.FindUsersByEmailRequest, server proto.UserService_FindUserByEmailServer) error {
	emails, err := util.CastSliceErr(request.Emails, func(from *string) (types.Email, error) {
		email := types.EmailFromString(*from)
		return email, email.Validate()
	})
	if err != nil {
		return err
	}

	users, stats := u.userService.FindByEmails(server.Context(), emails)
	if stats.IsError() {
		return stats.Error
	}

	for _, user := range users {
		response := mapper.ToProtoUserResponse(&user)
		if err := server.Send(&response); err != nil {
			return err
		}
	}

	return nil
}

func (u *UserHandler) BannedUser(ctx context.Context, request *proto.BannedUserRequest) (*emptypb.Empty, error) {
	// TODO: Get user id from access token claims from ctx
	dtoInput := mapper.ToDTOUserBannedInput(request)

	err := util.GetValidator().StructCtx(ctx, &dtoInput)
	if err != nil {
		return nil, err
	}

	stats := u.userService.BannedUser(ctx, &dtoInput)
	if stats.IsError() {
		return nil, stats.Error
	}

	return &emptypb.Empty{}, nil
}

func (u *UserHandler) DeleteUser(ctx context.Context, request *proto.DeleteUserRequest) (*emptypb.Empty, error) {
	// TODO: Get user id from access token claims from ctx
	id := types.IdFromString(request.Id)
	if err := id.Validate(); err != nil {
		return nil, err
	}

	stats := u.userService.DeleteById(ctx, id)
	if stats.IsError() {
		return nil, stats.Error
	}

	return &emptypb.Empty{}, nil
}
