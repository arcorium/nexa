package handler

import (
  "context"
  authNv1 "github.com/arcorium/nexa/proto/gen/go/authentication/v1"
  "github.com/arcorium/nexa/proto/gen/go/common"
  "github.com/arcorium/nexa/shared/dto"
  sharedErr "github.com/arcorium/nexa/shared/errors"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  "google.golang.org/protobuf/types/known/emptypb"
  "nexa/services/authentication/internal/api/grpc/mapper"
  "nexa/services/authentication/internal/domain/service"
  "nexa/services/authentication/util"
)

func NewUser(user service.IUser) UserHandler {
  return UserHandler{
    userService: user,
    tracer:      util.GetTracer(),
  }
}

type UserHandler struct {
  authNv1.UnimplementedUserServiceServer

  userService service.IUser
  tracer      trace.Tracer
}

func (u *UserHandler) Register(server *grpc.Server) {
  authNv1.RegisterUserServiceServer(server, u)
}

func (u *UserHandler) Create(ctx context.Context, request *authNv1.CreateUserRequest) (*authNv1.CreateUserResponse, error) {
  ctx, span := u.tracer.Start(ctx, "UserHandler.Create")
  defer span.End()

  dtoInput, err := mapper.ToUserCreateDTO(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  id, stats := u.userService.Create(ctx, &dtoInput)
  if stats.IsError() {
    spanUtil.RecordError(stats.Error, span)
    return nil, stats.ToGRPCError()
  }

  resp := &authNv1.CreateUserResponse{
    Id: id.String(),
  }
  return resp, nil
}

func (u *UserHandler) Update(ctx context.Context, request *authNv1.UpdateUserRequest) (*emptypb.Empty, error) {
  ctx, span := u.tracer.Start(ctx, "UserHandler.Update")
  defer span.End()

  dtoInput, err := mapper.ToUserUpdateDTO(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  stat := u.userService.Update(ctx, &dtoInput)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (u *UserHandler) UpdatePassword(ctx context.Context, request *authNv1.UpdateUserPasswordRequest) (*emptypb.Empty, error) {
  ctx, span := u.tracer.Start(ctx, "UserHandler.UpdatePassword")
  defer span.End()

  dtoInput, err := mapper.ToUserUpdatePasswordDTO(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  stat := u.userService.UpdatePassword(ctx, &dtoInput)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (u *UserHandler) Find(ctx context.Context, input *common.PagedElementInput) (*authNv1.FindUsersResponse, error) {
  ctx, span := u.tracer.Start(ctx, "UserHandler.Find")
  defer span.End()

  pagedDto := dto.PagedElementDTO{
    Element: input.Element,
    Page:    input.Page,
  }

  result, stat := u.userService.GetAll(ctx, pagedDto)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  resp := &authNv1.FindUsersResponse{
    Details: &common.PagedElementOutput{
      Element:       result.Element,
      Page:          result.Page,
      TotalElements: result.TotalElements,
      TotalPages:    result.TotalPages,
    },
    Users: sharedUtil.CastSliceP(result.Data, mapper.ToProtoUser),
  }
  return resp, nil
}

func (u *UserHandler) FindByIds(ctx context.Context, request *authNv1.FindUsersByIdsRequest) (*authNv1.FindUserByIdsResponse, error) {
  ctx, span := u.tracer.Start(ctx, "UserHandler.FindByIds")
  defer span.End()

  userIds, ierr := sharedUtil.CastSliceErrs(request.Ids, types.IdFromString)
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    err := sharedErr.GrpcFieldIndexedErrors("ids", ierr)
    return nil, err
  }

  users, stat := u.userService.FindByIds(ctx, userIds...)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  resp := &authNv1.FindUserByIdsResponse{
    Users: sharedUtil.CastSliceP(users, mapper.ToProtoUser),
  }
  return resp, nil
}

func (u *UserHandler) Banned(ctx context.Context, request *authNv1.BannedUserRequest) (*emptypb.Empty, error) {
  ctx, span := u.tracer.Start(ctx, "UserHandler.Banned")
  defer span.End()

  // Map and validation
  dtoInput, err := mapper.ToDTOUserBannedInput(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  stats := u.userService.BannedUser(ctx, &dtoInput)
  return nil, stats.ToGRPCErrorWithSpan(span)
}

func (u *UserHandler) Delete(ctx context.Context, request *authNv1.DeleteUserRequest) (*emptypb.Empty, error) {
  ctx, span := u.tracer.Start(ctx, "UserHandler.Delete")
  defer span.End()

  // Validation
  userId, err := types.IdFromString(request.Ids)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, sharedErr.NewFieldError("ids", err).ToGrpcError()
  }

  stats := u.userService.DeleteById(ctx, userId)
  return nil, stats.ToGRPCErrorWithSpan(span)
}

func (u *UserHandler) ForgotPassword(ctx context.Context, request *authNv1.ForgotUserPasswordRequest) (*authNv1.ForgotUserPasswordResponse, error) {
  ctx, span := u.tracer.Start(ctx, "UserHandler.ForgotPassword")
  defer span.End()

  recipientEmail, err := types.EmailFromString(request.Email)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, sharedErr.NewFieldError("email", err).ToGrpcError()
  }

  tokenResp, stat := u.userService.ForgotPassword(ctx, recipientEmail)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  return &authNv1.ForgotUserPasswordResponse{
    Token: mapper.ToProtoTokenResponse(&tokenResp),
  }, nil
}

func (u *UserHandler) ResetPasswordByToken(ctx context.Context, request *authNv1.ResetPasswordByTokenRequest) (*emptypb.Empty, error) {
  ctx, span := u.tracer.Start(ctx, "UserHandler.ForgotPassword")
  defer span.End()

  resetDTO, err := mapper.ToResetPasswordByTokenDTO(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  stat := u.userService.ResetPasswordWithToken(ctx, &resetDTO)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (u *UserHandler) ResetPassword(ctx context.Context, request *authNv1.ResetUserPasswordRequest) (*emptypb.Empty, error) {
  ctx, span := u.tracer.Start(ctx, "UserHandler.ResetPassword")
  defer span.End()

  resetDTO, err := mapper.ToResetUserPasswordDTO(request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  stats := u.userService.ResetPassword(ctx, &resetDTO)
  return nil, stats.ToGRPCErrorWithSpan(span)
}

func (u *UserHandler) VerifyEmail(ctx context.Context, request *authNv1.VerifyUserEmailRequest) (*authNv1.VerifyUserEmailResponse, error) {
  ctx, span := u.tracer.Start(ctx, "UserHandler.VerifyEmail")
  defer span.End()

  token := types.NewNullable(request.Token)
  if token.HasValue() {
    // Verify
    stat := u.userService.VerifyEmail(ctx, token.RawValue())
    if stat.IsError() {
      spanUtil.RecordError(stat.Error, span)
      return nil, stat.ToGRPCError()
    }
    return nil, nil
  }

  // Request
  tokenResp, stat := u.userService.EmailVerificationRequest(ctx)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  return &authNv1.VerifyUserEmailResponse{
    Token: mapper.ToProtoTokenResponse(&tokenResp),
  }, nil
}
