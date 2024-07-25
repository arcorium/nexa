package handler

import (
  "context"
  authNv1 "github.com/arcorium/nexa/proto/gen/go/authentication/v1"
  "github.com/arcorium/nexa/proto/gen/go/common"
  sharedDto "github.com/arcorium/nexa/shared/dto"
  sharedErr "github.com/arcorium/nexa/shared/errors"
  "github.com/arcorium/nexa/shared/grpc/interceptor/authz"
  "github.com/arcorium/nexa/shared/jwt"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  "google.golang.org/protobuf/types/known/emptypb"
  "io"
  "nexa/services/authentication/internal/api/grpc/mapper"
  "nexa/services/authentication/internal/domain/dto"
  "nexa/services/authentication/internal/domain/service"
  "nexa/services/authentication/util"
)

func NewUser(user service.IUser) UserHandler {
  return UserHandler{
    svc:    user,
    tracer: util.GetTracer(),
  }
}

type UserHandler struct {
  authNv1.UnimplementedUserServiceServer

  svc    service.IUser
  tracer trace.Tracer
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

  id, stats := u.svc.Create(ctx, &dtoInput)
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

  claims := types.Must(jwt.GetUserClaimsFromCtx(ctx))
  updateDTO, err := mapper.ToUserUpdateDTO(claims, request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  stat := u.svc.Update(ctx, &updateDTO)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (u *UserHandler) UpdateAvatar(server authNv1.UserService_UpdateAvatarServer) error {
  ctx := authz.GetWrappedContext(server)
  ctx, span := u.tracer.Start(ctx, "UserHandler.UpdateAvatar")
  defer span.End()

  var filename string
  var bytes []byte
  for {
    recv, err := server.Recv()
    if err != nil {
      if err == io.EOF {
        break
      }
      spanUtil.RecordError(err, span)
      return err
    }
    filename = recv.Filename
    bytes = append(bytes, recv.Chunk...)
  }

  claims := types.Must(jwt.GetUserClaimsFromCtx(ctx))
  userId, err := types.IdFromString(claims.UserId)
  if err != nil {
    spanUtil.RecordError(err, span)
    err = sharedErr.GrpcFieldErrors2(sharedErr.NewFieldError("user_id", err))
    return err
  }

  // Mapping and Validation
  dtoInput := dto.UpdateUserAvatarDTO{
    UserId:   userId,
    Filename: filename,
    Bytes:    bytes,
  }

  err = sharedUtil.ValidateStructCtx(ctx, &dtoInput)
  if err != nil {
    spanUtil.RecordError(err, span)
    return err
  }

  stat := u.svc.UpdateAvatar(ctx, &dtoInput)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return server.SendAndClose(nil)
  }

  return server.SendAndClose(nil)
}

func (u *UserHandler) UpdatePassword(ctx context.Context, request *authNv1.UpdateUserPasswordRequest) (*emptypb.Empty, error) {
  ctx, span := u.tracer.Start(ctx, "UserHandler.UpdatePassword")
  defer span.End()

  claims := types.Must(jwt.GetUserClaimsFromCtx(ctx))
  dtoInput, err := mapper.ToUserUpdatePasswordDTO(claims, request)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  stat := u.svc.UpdatePassword(ctx, &dtoInput)
  return nil, stat.ToGRPCErrorWithSpan(span)
}

func (u *UserHandler) Find(ctx context.Context, input *common.PagedElementInput) (*authNv1.FindUsersResponse, error) {
  ctx, span := u.tracer.Start(ctx, "UserHandler.Find")
  defer span.End()

  pagedDto := sharedDto.PagedElementDTO{
    Element: input.Element,
    Page:    input.Page,
  }

  result, stat := u.svc.GetAll(ctx, pagedDto)
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

  // Nil means to get the user itself
  if len(request.Ids) == 0 {
    // Get id from claims
    claims := types.Must(jwt.GetUserClaimsFromCtx(ctx))
    request.Ids = append(request.Ids, claims.UserId)
  }
  userIds, ierr := sharedUtil.CastSliceErrs(request.Ids, types.IdFromString)
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    err := sharedErr.GrpcFieldIndexedErrors("ids", ierr)
    return nil, err
  }

  users, stat := u.svc.FindByIds(ctx, userIds...)
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

  stats := u.svc.BannedUser(ctx, &dtoInput)
  return nil, stats.ToGRPCErrorWithSpan(span)
}

func (u *UserHandler) DeleteAvatar(ctx context.Context, request *authNv1.DeleteProfileAvatarRequest) (*emptypb.Empty, error) {
  ctx, span := u.tracer.Start(ctx, "UserHandler.DeleteAvatar")
  defer span.End()

  claims := types.Must(jwt.GetUserClaimsFromCtx(ctx))
  id := types.NewNullable(request.UserId)
  userId, err := types.IdFromString(id.ValueOr(claims.UserId))
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, sharedErr.NewFieldError("user_id", err).ToGrpcError()
  }

  stat := u.svc.DeleteAvatar(ctx, userId)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  return nil, nil
}

func (u *UserHandler) Delete(ctx context.Context, request *authNv1.DeleteUserRequest) (*emptypb.Empty, error) {
  ctx, span := u.tracer.Start(ctx, "UserHandler.Delete")
  defer span.End()

  // Validation
  claims := types.Must(jwt.GetUserClaimsFromCtx(ctx))
  id := types.NewNullable(request.UserId)
  userId, err := types.IdFromString(id.ValueOr(claims.UserId))
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, sharedErr.NewFieldError("ids", err).ToGrpcError()
  }

  stats := u.svc.DeleteById(ctx, userId)
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

  tokenResp, stat := u.svc.ForgotPassword(ctx, recipientEmail)
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

  stat := u.svc.ResetPasswordWithToken(ctx, &resetDTO)
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

  stats := u.svc.ResetPassword(ctx, &resetDTO)
  return nil, stats.ToGRPCErrorWithSpan(span)
}

func (u *UserHandler) VerifyEmail(ctx context.Context, request *authNv1.VerifyUserEmailRequest) (*emptypb.Empty, error) {
  ctx, span := u.tracer.Start(ctx, "UserHandler.VerifyEmail")
  defer span.End()

  stat := u.svc.VerifyEmail(ctx, request.Token)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }
  return nil, nil
}

func (u *UserHandler) EmailVerificationRequest(ctx context.Context, request *authNv1.UserEmailVerificationRequest) (*authNv1.VerifyUserEmailResponse, error) {
  ctx, span := u.tracer.Start(ctx, "UserHandler.EmailVerificationRequest")
  defer span.End()

  // Request
  result, stat := u.svc.EmailVerificationRequest(ctx)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return nil, stat.ToGRPCError()
  }

  resp := &authNv1.VerifyUserEmailResponse{
    Token: mapper.ToProtoTokenResponse(&result),
  }
  return resp, nil
}
