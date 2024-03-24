package service

import (
	"context"
	userUow "nexa/services/user/internal/app/uow"
	"nexa/services/user/internal/domain/service"
	"nexa/services/user/shared/domain/dto"
	"nexa/services/user/shared/domain/entity"
	"nexa/services/user/shared/domain/mapper"
	"nexa/shared/status"
	"nexa/shared/types"
	"nexa/shared/uow"
	"nexa/shared/util"
)

func NewUser(work uow.IUnitOfWork[userUow.UserStorage]) service.IUser {
	return &userService{
		unit: work,
	}
}

type userService struct {
	unit uow.IUnitOfWork[userUow.UserStorage]
}

func (u userService) Create(ctx context.Context, input *dto.UserCreateInput) status.Object {
	user, profile := mapper.MapUserCreateInput(input)

	err := u.unit.DoTx(ctx, func(ctx context.Context, storage userUow.UserStorage) error {
		err := storage.User().Create(ctx, &user)
		if err != nil {
			return err
		}

		err = storage.Profile().Create(ctx, &profile)
		return err
	})

	if err != nil {
		return status.FromRepository(err, status.NullCode)
	}
	return status.Created()
}

func (u userService) Update(ctx context.Context, input *dto.UserUpdateInput) status.Object {
	user := mapper.MapUserUpdateInput(input)
	err := u.unit.Repositories().User().Patch(ctx, &user)
	if err != nil {
		return status.FromRepository(err, status.NullCode)
	}
	return status.Updated()
}

func (u userService) UpdatePassword(ctx context.Context, input *dto.UserUpdatePasswordInput) status.Object {
	// Get user
	user := mapper.MapUserUpdatePasswordInput(input)
	users, err := u.unit.Repositories().User().FindByIds(ctx, user.Id)
	if err != nil {
		return status.FromRepository(err, status.NullCode)
	}

	// Validate last password
	err = users[0].Password.Equal(input.LastPassword)
	if err != nil {
		return status.Internal(err)
	}

	// Hash new password
	user.Password, err = user.Password.Hash()
	if err != nil {
		return status.Internal(err)
	}

	// Update
	err = u.unit.Repositories().User().Patch(ctx, &user)
	if err != nil {
		return status.FromRepository(err, status.NullCode)
	}
	return status.Updated()
}

func (u userService) ResetPassword(ctx context.Context, input *dto.UserResetPasswordInput) status.Object {
	user := mapper.MapUserResetPasswordInput(input)
	err := u.unit.Repositories().User().Patch(ctx, &user)
	if err != nil {
		return status.FromRepository(err, status.NullCode)
	}
	return status.Updated()
}

func (u userService) UpdateVerified(ctx context.Context, id types.Id) status.Object {
	user := entity.User{Id: id}
	err := u.unit.Repositories().User().Patch(ctx, &user)
	if err != nil {
		return status.FromRepository(err, status.NullCode)
	}
	return status.Updated()
}

func (u userService) FindByEmails(ctx context.Context, emails []types.Email) ([]dto.UserResponse, status.Object) {
	users, err := u.unit.Repositories().User().FindByEmails(ctx, emails...)
	if err != nil {
		return nil, status.FromRepository(err, status.NullCode)
	}
	responses := util.CastSlice(users, mapper.ToUserResponse)
	return responses, status.Success()
}

func (u userService) FindByIds(ctx context.Context, ids []types.Id) ([]dto.UserResponse, status.Object) {
	users, err := u.unit.Repositories().User().FindByIds(ctx, ids...)
	if err != nil {
		return nil, status.FromRepository(err, status.NullCode)
	}
	responses := util.CastSlice(users, mapper.ToUserResponse)
	return responses, status.Success()
}

func (u userService) DeleteById(ctx context.Context, id types.Id) status.Object {
	err := u.unit.Repositories().User().Delete(ctx, id)
	// TODO: Communicate to auth service to delete refresh token for this id
	if err != nil {
		return status.FromRepository(err, status.NullCode)
	}
	return status.Deleted()
}

func (u userService) BannedUser(ctx context.Context, input *dto.UserBannedInput) status.Object {
	// TODO: Check role and permission
	user := mapper.MapUserBannedInput(input)
	err := u.unit.Repositories().User().Patch(ctx, &user)
	if err != nil {
		return status.FromRepository(err, status.NullCode)
	}
	return status.Success()
}
