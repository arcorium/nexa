package service

import (
  "context"
  "go.opentelemetry.io/otel/trace"
  userUow "nexa/services/user/internal/app/uow"
  "nexa/services/user/internal/domain/dto"
  "nexa/services/user/internal/domain/entity"
  "nexa/services/user/internal/domain/external"
  "nexa/services/user/internal/domain/mapper"
  "nexa/services/user/internal/domain/service"
  spanUtil "nexa/shared/span"
  "nexa/shared/status"
  "nexa/shared/types"
  "nexa/shared/uow"
  sharedUtil "nexa/shared/util"
)

func NewUser(mailClient external.IMailerClient, work uow.IUnitOfWork[userUow.UserStorage], tracer trace.Tracer) service.IUser {
  return &userService{
    unit:       work,
    mailClient: mailClient,
    tracer:     tracer,
  }
}

type userService struct {
  unit   uow.IUnitOfWork[userUow.UserStorage]
  tracer trace.Tracer

  mailClient external.IMailerClient
  authClient external.IAuthenticationClient
}

func (u userService) Create(ctx context.Context, input *dto.UserCreateDTO) status.Object {
  ctx, span := u.tracer.Start(ctx, "UserService.Create")
  defer span.End()

  user, profile, stat := mapper.MapUserCreateDTO(input)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return stat
  }

  err := u.unit.DoTx(ctx, func(ctx context.Context, storage userUow.UserStorage) error {
    err := storage.User().Create(ctx, &user)
    if err != nil {
      return err
    }

    err = storage.Profile().Create(ctx, &profile)
    return err
  })
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }

  // Create verification token
  tokenGenerationDTO := dto.NewEmailVerificationToken()
  token, err := u.authClient.GenerateToken(ctx, &tokenGenerationDTO)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrExternal(err)
  }

  // Send email verification
  mailDTO := dto.SendEmailVerificationDTO{
    Recipient: user.Email.Underlying(),
    Token:     token.Token,
  }
  err = u.mailClient.SendEmailVerification(ctx, &mailDTO)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrExternal(err)
  }

  return status.Created()
}

func (u userService) Update(ctx context.Context, input *dto.UserUpdateDTO) status.Object {
  ctx, span := u.tracer.Start(ctx, "UserService.Update")
  defer span.End()

  user, stat := mapper.MapUserUpdateDTO(input)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return stat
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
  // Get user
  user, stat := mapper.MapUserUpdatePasswordDTO(input)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return stat
  }

  repo := u.unit.Repositories()

  users, err := repo.User().FindByIds(ctx, user.Id)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }

  // Validate last password
  err = users[0].Password.Equal(input.LastPassword)
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

func (u userService) ResetPassword(ctx context.Context, input *dto.UserResetPasswordDTO) status.Object {
  ctx, span := u.tracer.Start(ctx, "UserService.ResetPassword")
  defer span.End()

  user, stat := mapper.MapUserResetPasswordDTO(input)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return stat
  }

  repo := u.unit.Repositories()
  err := repo.User().Patch(ctx, &user)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }
  return status.Updated()
}

func (u userService) UpdateVerified(ctx context.Context, id types.Id) status.Object {
  ctx, span := u.tracer.Start(ctx, "UserService.UpdateVerified")
  defer span.End()

  user := entity.User{Id: id, IsVerified: true}

  repo := u.unit.Repositories()
  err := repo.User().Patch(ctx, &user)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }
  return status.Updated()
}

func (u userService) FindByEmails(ctx context.Context, emails []types.Email) ([]dto.UserResponseDTO, status.Object) {
  ctx, span := u.tracer.Start(ctx, "UserService.FindByEmails")
  defer span.End()

  repo := u.unit.Repositories()
  users, err := repo.User().FindByEmails(ctx, emails...)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, status.FromRepository(err, status.NullCode)
  }
  responses := sharedUtil.CastSliceP(users, mapper.ToUserResponse)
  return responses, status.Success()
}

func (u userService) FindByIds(ctx context.Context, ids []types.Id) ([]dto.UserResponseDTO, status.Object) {
  ctx, span := u.tracer.Start(ctx, "UserService.FindByIds")
  defer span.End()

  repo := u.unit.Repositories()
  users, err := repo.User().FindByIds(ctx, ids...)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, status.FromRepository(err, status.NullCode)
  }
  responses := sharedUtil.CastSliceP(users, mapper.ToUserResponse)
  return responses, status.Success()
}

func (u userService) DeleteById(ctx context.Context, id types.Id) status.Object {
  ctx, span := u.tracer.Start(ctx, "UserService.DeleteById")
  defer span.End()

  repo := u.unit.Repositories()
  err := repo.User().Delete(ctx, []types.Id{id}...)
  // TODO: Communicate to authentication service to delete refresh token for this id
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }
  return status.Deleted()
}

func (u userService) BannedUser(ctx context.Context, input *dto.UserBannedDTO) status.Object {
  ctx, span := u.tracer.Start(ctx, "UserService.BannedUser")
  defer span.End()

  // TODO: Check role and permission
  user := mapper.MapUserBannedDTO(input)
  repo := u.unit.Repositories()
  err := repo.User().Patch(ctx, &user)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }
  return status.Success()
}
