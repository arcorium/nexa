package handler

import (
  "context"
  "errors"
  "google.golang.org/genproto/googleapis/rpc/errdetails"
  "google.golang.org/grpc"
  "google.golang.org/protobuf/types/known/emptypb"
  proto "nexa/proto/generated/golang/user/v1"
  "nexa/services/user/internal/api/grpc/mapper"
  "nexa/services/user/internal/domain/service"
  sharedErr "nexa/shared/errors"
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

  err := util.ValidateStruct(ctx, &dtoInput)
  if err != nil {
    return nil, err
  }

  stats := u.userService.Create(ctx, &dtoInput)
  if stats.IsError() {
    return nil, stats.ToGRPCError()
  }
  return &emptypb.Empty{}, nil
}

func (u *UserHandler) Update(ctx context.Context, request *proto.UpdateUserRequest) (*emptypb.Empty, error) {
  dtoInput := mapper.ToDTOUserUpdateInput(request)
  // TODO: Get user id from access token claims from ctx

  err := util.ValidateStruct(ctx, &dtoInput)
  if err != nil {
    return nil, err
  }

  stats := u.userService.Update(ctx, &dtoInput)
  if stats.IsError() {
    return nil, stats.ToGRPCError()
  }
  return &emptypb.Empty{}, nil
}

func (u *UserHandler) UpdateVerified(ctx context.Context, request *proto.UpdateUserVerifiedRequest) (*emptypb.Empty, error) {
  id, err := types.IdFromString(request.Id)
  // TODO: Get user id from access token claims from ctx
  if errors.Is(err, types.ErrMalformedUUID) {
    return nil, sharedErr.GrpcFieldErrors(&errdetails.BadRequest_FieldViolation{
      Field:       "id",
      Description: err.Error(),
    })
  }

  stats := u.userService.UpdateVerified(ctx, id)
  if stats.IsError() {
    return nil, stats.ToGRPCError()
  }
  return &emptypb.Empty{}, nil
}

func (u *UserHandler) UpdatePassword(ctx context.Context, request *proto.UpdateUserPasswordRequest) (*emptypb.Empty, error) {
  dtoInput := mapper.ToDTOUserUpdatePasswordInput(request)
  // TODO: Get user id from access token claims from ctx

  err := util.ValidateStruct(ctx, &dtoInput)
  if err != nil {
    return nil, err
  }

  stats := u.userService.UpdatePassword(ctx, &dtoInput)
  if stats.IsError() {
    return nil, stats.ToGRPCError()
  }
  return &emptypb.Empty{}, nil
}

func (u *UserHandler) ResetPassword(ctx context.Context, request *proto.ResetUserPasswordRequest) (*emptypb.Empty, error) {
  dtoInput := mapper.ToDTOUserResetPasswordInput(request)

  err := util.ValidateStruct(ctx, &dtoInput)
  if err != nil {
    return nil, err
  }

  stats := u.userService.ResetPassword(ctx, &dtoInput)
  if stats.IsError() {
    return nil, stats.ToGRPCError()
  }
  return &emptypb.Empty{}, nil
}

func (u *UserHandler) FindUserByIds(ctx context.Context, request *proto.FindUsersByIdsRequest) (*proto.FindUserByIdsResponse, error) {
  ids, ierr := util.CastSliceErrs(request.Ids, func(from *string) (types.Id, error) {
    return types.IdFromString(*from)
  })

  if ierr != nil {
    return nil, sharedErr.GrpcFieldIndexedErrors("ids", ierr)
  }

  users, stats := u.userService.FindByIds(ctx, ids)
  return &proto.FindUserByIdsResponse{
    Users: util.CastSlice(users, mapper.ToProtoUser),
  }, stats.ToGRPCError()
}

func (u *UserHandler) FindUserByEmails(ctx context.Context, request *proto.FindUserByEmailsRequest) (*proto.FindUserByEmailsResponse, error) {
  emails, ierr := util.CastSliceErrs(request.Emails, func(from *string) (types.Email, error) {
    return types.EmailFromString(*from)
  })

  if ierr != nil {
    return nil, sharedErr.GrpcFieldIndexedErrors("emails", ierr)
  }

  users, stats := u.userService.FindByEmails(ctx, emails)
  return &proto.FindUserByEmailsResponse{
    Users: util.CastSlice(users, mapper.ToProtoUser),
  }, stats.ToGRPCError()
}

func (u *UserHandler) BannedUser(ctx context.Context, request *proto.BannedUserRequest) (*emptypb.Empty, error) {
  // TODO: Get user id from access token claims from ctx
  dtoInput := mapper.ToDTOUserBannedInput(request)

  err := util.ValidateStruct(ctx, &dtoInput)
  if err != nil {
    return nil, err
  }

  stats := u.userService.BannedUser(ctx, &dtoInput)
  if stats.IsError() {
    return nil, stats.ToGRPCError()
  }
  return &emptypb.Empty{}, nil
}

func (u *UserHandler) DeleteUser(ctx context.Context, request *proto.DeleteUserRequest) (*emptypb.Empty, error) {
  // TODO: Get user id from access token claims from ctx
  id, err := types.IdFromString(request.Ids)
  if err != nil {
    if errors.Is(err, types.ErrMalformedUUID) {
      return nil, sharedErr.GrpcFieldErrors(&errdetails.BadRequest_FieldViolation{
        Field:       "id",
        Description: err.Error(),
      })
    }
    return nil, err
  }

  stats := u.userService.DeleteById(ctx, id)
  if stats.IsError() {
    return nil, stats.ToGRPCError()
  }
  return &emptypb.Empty{}, nil
}

func (u *UserHandler) Validate(ctx context.Context, request *proto.ValidateUserRequest) (*proto.ValidateUserResponse, error) {
  // TODO: Implement it
  return nil, nil
}
