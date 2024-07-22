package service

import (
  "context"
  "database/sql"
  sharedDto "github.com/arcorium/nexa/shared/dto"
  sharedErr "github.com/arcorium/nexa/shared/errors"
  sharedJwt "github.com/arcorium/nexa/shared/jwt"
  "github.com/arcorium/nexa/shared/status"
  "github.com/arcorium/nexa/shared/types"
  "github.com/arcorium/nexa/shared/uow"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  authUtil "github.com/arcorium/nexa/shared/util/auth"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "go.opentelemetry.io/otel/trace"
  "nexa/services/authentication/constant"
  userUow "nexa/services/authentication/internal/app/uow"
  "nexa/services/authentication/internal/domain/dto"
  "nexa/services/authentication/internal/domain/entity"
  "nexa/services/authentication/internal/domain/external"
  "nexa/services/authentication/internal/domain/mapper"
  "nexa/services/authentication/internal/domain/repository"
  "nexa/services/authentication/internal/domain/service"
  "nexa/services/authentication/util"
)

func NewUser(credRepo repository.ICredential, unit uow.IUnitOfWork[userUow.UserStorage], config UserConfig) service.IUser {
  return &userService{
    unit:     unit,
    credRepo: credRepo,
    tracer:   util.GetTracer(),
    config:   config,
  }
}

type UserConfig struct {
  RoleClient    external.IRoleClient
  MailClient    external.IMailerClient
  TokenClient   external.ITokenClient
  StorageClient external.IFileStorageClient
}

type userService struct {
  unit     uow.IUnitOfWork[userUow.UserStorage]
  credRepo repository.ICredential
  tracer   trace.Tracer

  config UserConfig
}

func (u userService) checkPermission(ctx context.Context, targetId types.Id, permission string) error {
  // Validate permission
  claims, _ := sharedJwt.GetUserClaimsFromCtx(ctx)
  if !targetId.EqWithString(claims.UserId) {
    // Need permission to update other users
    if !authUtil.ContainsPermission(claims.Roles, permission) {
      return sharedErr.ErrUnauthorizedPermission
    }
  }
  return nil
}

func (u userService) Create(ctx context.Context, createDto *dto.UserCreateDTO) (types.Id, status.Object) {
  ctx, span := u.tracer.Start(ctx, "UserService.Create")
  defer span.End()

  user, profile, err := createDto.ToDomain()
  if err != nil {
    spanUtil.RecordError(err, span)
    return types.NullId(), status.ErrBadRequest(err)
  }

  err = u.unit.DoTx(ctx, func(ctx context.Context, storage userUow.UserStorage) error {
    err := storage.User().Create(ctx, &user)
    if err != nil {
      return err
    }

    err = storage.Profile().Create(ctx, &profile)
    return err
  })
  if err != nil {
    spanUtil.RecordError(err, span)
    return types.NullId(), status.FromRepository(err, status.NullCode)
  }

  return user.Id, status.Created()
}

func (u userService) Update(ctx context.Context, updateDto *dto.UserUpdateDTO) status.Object {
  ctx, span := u.tracer.Start(ctx, "UserService.Update")
  defer span.End()

  // Validate permission
  if err := u.checkPermission(ctx, updateDto.Id, constant.AUTHN_PERMISSIONS[constant.AUTHN_UPDATE_USER_ARB]); err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrUnAuthorized(err)
  }

  user, profile := updateDto.ToDomain()

  err := u.unit.DoTx(ctx, func(ctx context.Context, storage userUow.UserStorage) error {
    ctx, span := u.tracer.Start(ctx, "UOW.Update")
    defer span.End()

    err := storage.User().Patch(ctx, &user)
    if err != nil {
      return err
    }

    err = storage.Profile().Patch(ctx, &profile)
    return err
  })

  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }

  return status.Updated()
}

func (u userService) UpdateAvatar(ctx context.Context, updateDto *dto.UpdateUserAvatarDTO) status.Object {
  ctx, span := u.tracer.Start(ctx, "UserService.UpdateAvatar")
  defer span.End()

  // Check if user already has photo
  repos := u.unit.Repositories()
  profiles, err := repos.Profile().FindByIds(ctx, updateDto.UserId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }

  // Upload new avatar
  fileId, filePath, err := u.config.StorageClient.UploadProfileImage(ctx, &dto.UploadImageDTO{
    Filename: updateDto.Filename,
    Data:     updateDto.Bytes,
  })
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrExternal(err)
  }

  // Update profiles data
  profile := entity.PatchedProfile{
    Id:       updateDto.UserId,
    PhotoId:  types.SomeNullable(fileId),
    PhotoURL: types.SomeNullable(filePath),
  }

  err = repos.Profile().Patch(ctx, &profile)
  if err != nil {
    // Delete new avatar when error happens
    spanUtil.RecordError(err, span)
    extErr := u.config.StorageClient.DeleteProfileImage(ctx, fileId)
    if extErr != nil {
      spanUtil.RecordError(extErr, span)
      return status.ErrExternal(extErr)
    }
    return status.FromRepository(err, status.NullCode)
  }

  // Delete last avatar
  if profiles[0].HasAvatar() {
    err = u.config.StorageClient.DeleteProfileImage(ctx, profiles[0].PhotoId)
    if err != nil {
      spanUtil.RecordError(err, span)
      return status.ErrExternal(err)
    }
  }

  return status.Updated()
}

func (u userService) UpdatePassword(ctx context.Context, updateDto *dto.UserUpdatePasswordDTO) status.Object {
  ctx, span := u.tracer.Start(ctx, "UserService.UpdatePassword")
  defer span.End()

  // Mapping and validate
  user, err := updateDto.ToDomain()
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrInternal(err)
  }

  // Validate permission
  if err := u.checkPermission(ctx, user.Id, constant.AUTHN_PERMISSIONS[constant.AUTHN_UPDATE_USER_ARB]); err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrUnAuthorized(err)
  }

  repos := u.unit.Repositories()
  users, err := repos.User().FindByIds(ctx, user.Id)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }

  // Validate last password
  err = users[0].ValidatePassword(updateDto.LastPassword)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrInternal(err)
  }

  // Update
  err = repos.User().Patch(ctx, &user)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }
  return status.Updated()
}

func (u userService) BannedUser(ctx context.Context, bannedDto *dto.UserBannedDTO) status.Object {
  ctx, span := u.tracer.Start(ctx, "UserService.BannedUser")
  defer span.End()

  user := bannedDto.ToDomain()

  repos := u.unit.Repositories()
  err := repos.User().Patch(ctx, &user)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }
  return status.Success()
}

func (u userService) GetAll(ctx context.Context, pagedDto sharedDto.PagedElementDTO) (sharedDto.PagedElementResult[dto.UserResponseDTO], status.Object) {
  ctx, span := u.tracer.Start(ctx, "UserService.FindByEmails")
  defer span.End()

  repos := u.unit.Repositories()

  result, err := repos.User().Get(ctx, pagedDto.ToQueryParam())
  if err != nil {
    spanUtil.RecordError(err, span)
    return sharedDto.PagedElementResult[dto.UserResponseDTO]{}, status.FromRepository(err, status.NullCode)
  }

  responseDtos := sharedUtil.CastSliceP(result.Data, mapper.ToUserResponse)

  response := sharedDto.NewPagedElementResult2(responseDtos, &pagedDto, result.Total)
  return response, status.Success()
}

func (u userService) FindByIds(ctx context.Context, ids ...types.Id) ([]dto.UserResponseDTO, status.Object) {
  ctx, span := u.tracer.Start(ctx, "UserService.FindByIds")
  defer span.End()

  // Find by profiles
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
  if err := u.checkPermission(ctx, userId, constant.AUTHN_PERMISSIONS[constant.AUTHN_DELETE_USER_ARB]); err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrUnAuthorized(sharedErr.ErrUnauthorized)
  }

  // Delete from repo
  var stat = status.Deleted()
  _ = u.unit.DoTx(ctx, func(ctx context.Context, storage userUow.UserStorage) error {
    err := storage.User().Delete(ctx, userId)
    if err != nil {
      spanUtil.RecordError(err, span)
      stat = status.FromRepository(err, status.NullCode)
      return err
    }

    // Delete all user credentials
    err = u.credRepo.DeleteByUserId(ctx, userId)
    if err != nil && err != sql.ErrNoRows {
      span.RecordError(err)
      stat = status.FromRepository(err, status.NullCode)
      return err
    }

    err = u.config.RoleClient.RemoveUserRoles(ctx, userId)
    if err != nil {
      span.RecordError(err)
      stat = status.ErrExternal(err)
      return err
    }

    return nil
  })

  return stat
}

func (u userService) VerifyEmail(ctx context.Context, token string) status.Object {
  ctx, span := u.tracer.Start(ctx, "UserService.VerifyEmail")
  defer span.End()

  // Verify the token
  userId, err := u.config.TokenClient.Verify(ctx, &dto.TokenVerificationDTO{
    Token:   token,
    Purpose: dto.TokenUsageEmailVerification,
  })
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrExternal(err)
  }

  // Update is_verified field
  user := entity.PatchedUser{
    Id:         userId,
    IsVerified: types.SomeNullable(true),
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

  userClaims, _ := sharedJwt.GetUserClaimsFromCtx(ctx)
  userId, err := types.IdFromString(userClaims.UserId)
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

  // Generate token from auth client
  tokenReqDto := dto.TokenGenerationDTO{
    UserId: userId,
    Usage:  dto.TokenUsageEmailVerification,
  }
  token, err := u.config.TokenClient.Generate(ctx, &tokenReqDto)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.TokenResponseDTO{}, status.ErrExternal(err)
  }

  // Send token to email
  err = u.config.MailClient.SendEmailVerification(ctx, &dto.SendEmailVerificationDTO{
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
  tokenReq := dto.TokenGenerationDTO{
    UserId: users[0].Id,
    Usage:  dto.TokenUsageResetPassword,
  }
  token, err := u.config.TokenClient.Generate(ctx, &tokenReq)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.TokenResponseDTO{}, status.ErrExternal(err)
  }

  // Send to email
  err = u.config.MailClient.SendForgotPassword(ctx, &dto.SendForgotPasswordDTO{
    Recipient: users[0].Email,
    Token:     token.Token,
  })
  if err != nil {
    spanUtil.RecordError(err, span)
  }
  return token, status.Success()
}

func (u userService) ResetPasswordWithToken(ctx context.Context, resetDTO *dto.ResetPasswordWithTokenDTO) status.Object {
  ctx, span := u.tracer.Start(ctx, "UserService.ResetPasswordWithToken")
  defer span.End()

  password, err := resetDTO.NewPassword.Hash()
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrInternal(err)
  }

  // Verify token
  userId, err := u.config.TokenClient.Verify(ctx, &dto.TokenVerificationDTO{
    Token:   resetDTO.Token,
    Purpose: dto.TokenUsageResetPassword,
  })

  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrBadRequest(err)
  }

  user := entity.PatchedUser{
    Id:       userId,
    Password: password,
  }

  repos := u.unit.Repositories()
  err = repos.User().Patch(ctx, &user)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }

  if resetDTO.LogoutAll {
    err = u.credRepo.DeleteByUserId(ctx, user.Id)
    if err != nil {
      span.RecordError(err)
    }
  }

  return status.Updated()
}

func (u userService) ResetPassword(ctx context.Context, resetDTO *dto.ResetUserPasswordDTO) status.Object {
  ctx, span := u.tracer.Start(ctx, "UserService.ResetPassword")
  defer span.End()

  password, err := resetDTO.NewPassword.Hash()
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrInternal(err)
  }

  // Update user data
  user := entity.PatchedUser{
    Id:       resetDTO.UserId,
    Password: password,
  }

  repos := u.unit.Repositories()
  err = repos.User().Patch(ctx, &user)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }

  if resetDTO.LogoutAll {
    err = u.credRepo.DeleteByUserId(ctx, user.Id)
    if err != nil {
      span.RecordError(err)
    }
  }

  return status.Updated()
}
