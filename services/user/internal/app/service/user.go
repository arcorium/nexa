package service

import (
  "context"
  "go.opentelemetry.io/otel/trace"
  "nexa/services/user/constant"
  userUow "nexa/services/user/internal/app/uow"
  "nexa/services/user/internal/domain/dto"
  domain "nexa/services/user/internal/domain/entity"
  "nexa/services/user/internal/domain/external"
  "nexa/services/user/internal/domain/mapper"
  "nexa/services/user/internal/domain/service"
  "nexa/services/user/util"
  sharedDto "nexa/shared/dto"
  sharedErr "nexa/shared/errors"
  sharedJwt "nexa/shared/jwt"
  "nexa/shared/status"
  "nexa/shared/types"
  "nexa/shared/uow"
  sharedUtil "nexa/shared/util"
  authUtil "nexa/shared/util/auth"
  spanUtil "nexa/shared/util/span"
)

func NewUser(mailClient external.IMailerClient, authClient external.IAuthenticationClient, work uow.IUnitOfWork[userUow.UserStorage]) service.IUser {
  return &userService{
    unit:       work,
    mailClient: mailClient,
    authClient: authClient,
    tracer:     util.GetTracer(),
  }
}

type userService struct {
  unit   uow.IUnitOfWork[userUow.UserStorage]
  tracer trace.Tracer

  mailClient external.IMailerClient
  authClient external.IAuthenticationClient
}

func (u userService) checkPermission(ctx context.Context, targetId types.Id, permissions ...string) error {
  // Validate permission
  claims, _ := sharedJwt.GetClaimsFromCtx(ctx)
  if !targetId.EqWithString(claims.UserId) {
    // Need permission to update other users
    if !authUtil.ContainsPermissions(claims.Roles, permissions...) {
      return sharedErr.ErrUnauthorizedPermission
    }
  }
  return nil
}

func (u userService) Create(ctx context.Context, input *dto.UserCreateDTO) (string, status.Object) {
  ctx, span := u.tracer.Start(ctx, "UserService.Create")
  defer span.End()

  user, profile, err := input.ToDomain()
  if err != nil {
    spanUtil.RecordError(err, span)
    return "", status.ErrBadRequest(err)
  }

  err = u.unit.DoTx(ctx, func(ctx context.Context, storage userUow.UserStorage) error {
    err := storage.User().Create(ctx, user)
    if err != nil {
      return err
    }

    err = storage.Profile().Create(ctx, profile)
    return err
  })
  if err != nil {
    spanUtil.RecordError(err, span)
    return "", status.FromRepository(err, status.NullCode)
  }

  return user.Id.String(), status.Created()
}

func (u userService) Update(ctx context.Context, input *dto.UserUpdateDTO) status.Object {
  ctx, span := u.tracer.Start(ctx, "UserService.Update")
  defer span.End()

  user := input.ToDomain()

  // Validate permission
  if err := u.checkPermission(ctx, user.Id, constant.USER_PERMISSIONS[constant.USER_UPDATE_OTHER]); err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrUnAuthorized(sharedErr.ErrUnauthorized)
  }

  repo := u.unit.Repositories()
  err := repo.User().Patch(ctx, &user)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }
  return status.Updated()
}

func (u userService) UpdatePassword(ctx context.Context, input *dto.UserUpdatePasswordDTO) status.Object {
  ctx, span := u.tracer.Start(ctx, "UserService.UpdatePassword")
  defer span.End()

  // Mapping and validate
  user, err := input.ToDomain()
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrInternal(err)
  }

  // Validate permission
  if err := u.checkPermission(ctx, user.Id, constant.USER_PERMISSIONS[constant.USER_UPDATE_OTHER]); err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrUnAuthorized(sharedErr.ErrUnauthorized)
  }

  repo := u.unit.Repositories()

  users, err := repo.User().FindByIds(ctx, user.Id)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }

  // Validate last password
  err = users[0].Password.Equal(input.LastPassword.String())
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrInternal(err)
  }

  // Update
  err = repo.User().Patch(ctx, &user)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }
  return status.Updated()
}

func (u userService) BannedUser(ctx context.Context, input *dto.UserBannedDTO) status.Object {
  ctx, span := u.tracer.Start(ctx, "UserService.BannedUser")
  defer span.End()

  user := input.ToDomain()

  repos := u.unit.Repositories()
  err := repos.User().Patch(ctx, &user)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }
  return status.Success()
}

func (u userService) FindAll(ctx context.Context, pagedDto sharedDto.PagedElementDTO) (sharedDto.PagedElementResult[dto.UserResponseDTO], status.Object) {
  ctx, span := u.tracer.Start(ctx, "UserService.FindByEmails")
  defer span.End()

  repos := u.unit.Repositories()

  result, err := repos.User().Get(ctx, pagedDto.ToQueryParam())
  if err != nil {
    spanUtil.RecordError(err, span)
    return sharedDto.PagedElementResult[dto.UserResponseDTO]{}, status.FromRepository(err, status.NullCode)
  }

  responseDtos := sharedUtil.CastSliceP(result.Data, func(user *domain.User) dto.UserResponseDTO {
    return mapper.ToUserResponse(user)
  })

  response := sharedDto.NewPagedElementResult2(responseDtos, &pagedDto, result.Total)
  return response, status.Success()
}

func (u userService) FindByIds(ctx context.Context, ids ...types.Id) ([]dto.UserResponseDTO, status.Object) {
  ctx, span := u.tracer.Start(ctx, "UserService.FindByIds")
  defer span.End()

  repos := u.unit.Repositories()
  users, err := repos.User().FindByIds(ctx, ids...)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, status.FromRepository(err, status.NullCode)
  }

  responses := sharedUtil.CastSliceP(users, mapper.ToUserResponse)
  return responses, status.Success()
}

func (u userService) DeleteById(ctx context.Context, userId types.Id) status.Object {
  ctx, span := u.tracer.Start(ctx, "UserService.DeleteById")
  defer span.End()

  // Check target id with claims id and the permissions needed
  if err := u.checkPermission(ctx, userId, constant.USER_PERMISSIONS[constant.USER_DELETE_OTHER]); err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrUnAuthorized(sharedErr.ErrUnauthorized)
  }

  // Delete from repo
  repos := u.unit.Repositories()
  err := repos.User().Delete(ctx, userId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }

  // Delete all user credentials
  err = u.authClient.DeleteCredentials(ctx, userId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrExternal(err)
  }

  return status.Deleted()
}

func (u userService) Validate(ctx context.Context, email types.Email, password types.Password) (dto.UserResponseDTO, status.Object) {
  ctx, span := u.tracer.Start(ctx, "UserService.Validate")
  defer span.End()

  // Find user by email
  repos := u.unit.Repositories()
  user, err := repos.User().FindByEmails(ctx, email)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.UserResponseDTO{}, status.FromRepository(err, status.NullCode)
  }

  // Validate the password
  err = user[0].ValidatePassword(password.String())
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.UserResponseDTO{}, status.ErrBadRequest(err)
  }

  // return the user response
  response := mapper.ToUserResponse(&user[0])
  return response, status.Success()
}

func (u userService) VerifyEmail(ctx context.Context, token string) status.Object {
  ctx, span := u.tracer.Start(ctx, "UserService.VerifyEmail")
  defer span.End()

  // Verify the token
  userId, err := u.authClient.VerifyToken(ctx, &dto.TokenVerificationDTO{
    Token:   token,
    Purpose: dto.EmailVerificationToken,
  })
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrExternal(err)
  }

  // Update is_verified field
  verified := true
  user := domain.User{
    Id:         userId,
    IsVerified: &verified,
  }

  repos := u.unit.Repositories()
  err = repos.User().Patch(ctx, &user)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }

  return status.Updated()
}

func (u userService) EmailVerificationRequest(ctx context.Context) (dto.TokenResponseDTO, status.Object) {
  ctx, span := u.tracer.Start(ctx, "UserService.EmailVerificationRequest")
  defer span.End()

  userClaims, _ := sharedJwt.GetClaimsFromCtx(ctx)
  userId, err := types.IdFromString(userClaims.UserId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.TokenResponseDTO{}, status.ErrBadRequest(err)
  }

  // Generate token from auth client
  tokenReqDto := dto.NewEmailVerificationToken(userId)
  token, err := u.authClient.GenerateToken(ctx, &tokenReqDto)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.TokenResponseDTO{}, status.ErrBadRequest(err)
  }

  // Get user details
  repos := u.unit.Repositories()
  users, err := repos.User().FindByIds(ctx, userId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.TokenResponseDTO{}, status.FromRepository(err, status.NullCode)
  }

  // Send token to email
  err = u.mailClient.SendEmailVerification(ctx, &dto.SendEmailVerificationDTO{
    Recipient: users[0].Email,
    Token:     token.Token,
  })
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.TokenResponseDTO{}, status.ErrExternal(err)
  }

  return token, status.Success()
}

func (u userService) ForgotPassword(ctx context.Context, email types.Email) (dto.TokenResponseDTO, status.Object) {
  ctx, span := u.tracer.Start(ctx, "UserService.ForgotPassword")
  defer span.End()

  // Get user details
  repos := u.unit.Repositories()
  users, err := repos.User().FindByEmails(ctx, email)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.TokenResponseDTO{}, status.FromRepository(err, status.NullCode)
  }

  // Create token
  tokenReq := dto.NewForgotPasswordToken(users[0].Id)
  token, err := u.authClient.GenerateToken(ctx, &tokenReq)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.TokenResponseDTO{}, status.ErrExternal(err)
  }

  // Send to email
  err = u.mailClient.SendForgotPassword(ctx, &dto.SendForgotPasswordDTO{
    Recipient: users[0].Email,
    Token:     token.Token,
  })
  if err != nil {
    spanUtil.RecordError(err, span)
  }
  return token, status.Success()
}

func (u userService) ResetPassword(ctx context.Context, input *dto.UserResetPasswordDTO) status.Object {
  ctx, span := u.tracer.Start(ctx, "UserService.ResetPassword")
  defer span.End()

  password, err := input.NewPassword.Hash()
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrInternal(err)
  }

  var userId types.Id
  // Check if the request doesn't have token, so the authorized user can do reset user password without token
  if !input.Token.HasValue() {
    // Check claims
    userClaims, err := sharedJwt.GetClaimsFromCtx(ctx)
    if err != nil {
      spanUtil.RecordError(err, span)
      return status.ErrUnAuthenticated(err)
    }

    // Check permissions
    if !authUtil.ContainsPermissions(userClaims.Roles, constant.USER_PERMISSIONS[constant.USER_UPDATE_OTHER]) {
      spanUtil.RecordError(sharedErr.ErrUnauthorizedPermission, span)
      return status.ErrUnAuthorized(sharedErr.ErrUnauthorizedPermission)
    }
  } else {
    // Using token
    // Verify
    userId, err = u.authClient.VerifyToken(ctx, &dto.TokenVerificationDTO{
      Token:   input.Token.RawValue(),
      Purpose: dto.ForgotPasswordToken,
    })
    if err != nil {
      spanUtil.RecordError(err, span)
      return status.ErrExternal(err)
    }
  }

  // Update user data
  user := domain.User{
    Id:       userId,
    Password: password,
  }

  repos := u.unit.Repositories()
  err = repos.User().Patch(ctx, &user)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }

  // Logout all the credentials
  if input.LogoutAll {
    err = u.authClient.DeleteCredentials(ctx, user.Id)
    if err != nil {
      spanUtil.RecordError(err, span)
      return status.ErrExternal(err)
    }
  }
  // TODO: Send email that the password has changed

  return status.Updated()
}
