// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"
	dto "nexa/services/authentication/internal/domain/dto"

	mock "github.com/stretchr/testify/mock"

	shareddto "github.com/arcorium/nexa/shared/dto"

	status "github.com/arcorium/nexa/shared/status"

	types "github.com/arcorium/nexa/shared/types"
)

// UserMock is an autogenerated mock type for the IUser type
type UserMock struct {
	mock.Mock
}

type UserMock_Expecter struct {
	mock *mock.Mock
}

func (_m *UserMock) EXPECT() *UserMock_Expecter {
	return &UserMock_Expecter{mock: &_m.Mock}
}

// BannedUser provides a mock function with given fields: ctx, input
func (_m *UserMock) BannedUser(ctx context.Context, input *dto.UserBannedDTO) status.Object {
	ret := _m.Called(ctx, input)

	if len(ret) == 0 {
		panic("no return value specified for BannedUser")
	}

	var r0 status.Object
	if rf, ok := ret.Get(0).(func(context.Context, *dto.UserBannedDTO) status.Object); ok {
		r0 = rf(ctx, input)
	} else {
		r0 = ret.Get(0).(status.Object)
	}

	return r0
}

// UserMock_BannedUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'BannedUser'
type UserMock_BannedUser_Call struct {
	*mock.Call
}

// BannedUser is a helper method to define mock.On call
//   - ctx context.Context
//   - input *dto.UserBannedDTO
func (_e *UserMock_Expecter) BannedUser(ctx interface{}, input interface{}) *UserMock_BannedUser_Call {
	return &UserMock_BannedUser_Call{Call: _e.mock.On("BannedUser", ctx, input)}
}

func (_c *UserMock_BannedUser_Call) Run(run func(ctx context.Context, input *dto.UserBannedDTO)) *UserMock_BannedUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*dto.UserBannedDTO))
	})
	return _c
}

func (_c *UserMock_BannedUser_Call) Return(_a0 status.Object) *UserMock_BannedUser_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserMock_BannedUser_Call) RunAndReturn(run func(context.Context, *dto.UserBannedDTO) status.Object) *UserMock_BannedUser_Call {
	_c.Call.Return(run)
	return _c
}

// Create provides a mock function with given fields: ctx, input
func (_m *UserMock) Create(ctx context.Context, input *dto.UserCreateDTO) (types.Id, status.Object) {
	ret := _m.Called(ctx, input)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 types.Id
	var r1 status.Object
	if rf, ok := ret.Get(0).(func(context.Context, *dto.UserCreateDTO) (types.Id, status.Object)); ok {
		return rf(ctx, input)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *dto.UserCreateDTO) types.Id); ok {
		r0 = rf(ctx, input)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Id)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *dto.UserCreateDTO) status.Object); ok {
		r1 = rf(ctx, input)
	} else {
		r1 = ret.Get(1).(status.Object)
	}

	return r0, r1
}

// UserMock_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type UserMock_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - input *dto.UserCreateDTO
func (_e *UserMock_Expecter) Create(ctx interface{}, input interface{}) *UserMock_Create_Call {
	return &UserMock_Create_Call{Call: _e.mock.On("Create", ctx, input)}
}

func (_c *UserMock_Create_Call) Run(run func(ctx context.Context, input *dto.UserCreateDTO)) *UserMock_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*dto.UserCreateDTO))
	})
	return _c
}

func (_c *UserMock_Create_Call) Return(_a0 types.Id, _a1 status.Object) *UserMock_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserMock_Create_Call) RunAndReturn(run func(context.Context, *dto.UserCreateDTO) (types.Id, status.Object)) *UserMock_Create_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteById provides a mock function with given fields: ctx, id
func (_m *UserMock) DeleteById(ctx context.Context, id types.Id) status.Object {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for DeleteById")
	}

	var r0 status.Object
	if rf, ok := ret.Get(0).(func(context.Context, types.Id) status.Object); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(status.Object)
	}

	return r0
}

// UserMock_DeleteById_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteById'
type UserMock_DeleteById_Call struct {
	*mock.Call
}

// DeleteById is a helper method to define mock.On call
//   - ctx context.Context
//   - id types.Id
func (_e *UserMock_Expecter) DeleteById(ctx interface{}, id interface{}) *UserMock_DeleteById_Call {
	return &UserMock_DeleteById_Call{Call: _e.mock.On("DeleteById", ctx, id)}
}

func (_c *UserMock_DeleteById_Call) Run(run func(ctx context.Context, id types.Id)) *UserMock_DeleteById_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(types.Id))
	})
	return _c
}

func (_c *UserMock_DeleteById_Call) Return(_a0 status.Object) *UserMock_DeleteById_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserMock_DeleteById_Call) RunAndReturn(run func(context.Context, types.Id) status.Object) *UserMock_DeleteById_Call {
	_c.Call.Return(run)
	return _c
}

// EmailVerificationRequest provides a mock function with given fields: ctx
func (_m *UserMock) EmailVerificationRequest(ctx context.Context) (dto.TokenResponseDTO, status.Object) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for EmailVerificationRequest")
	}

	var r0 dto.TokenResponseDTO
	var r1 status.Object
	if rf, ok := ret.Get(0).(func(context.Context) (dto.TokenResponseDTO, status.Object)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) dto.TokenResponseDTO); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(dto.TokenResponseDTO)
	}

	if rf, ok := ret.Get(1).(func(context.Context) status.Object); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Get(1).(status.Object)
	}

	return r0, r1
}

// UserMock_EmailVerificationRequest_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EmailVerificationRequest'
type UserMock_EmailVerificationRequest_Call struct {
	*mock.Call
}

// EmailVerificationRequest is a helper method to define mock.On call
//   - ctx context.Context
func (_e *UserMock_Expecter) EmailVerificationRequest(ctx interface{}) *UserMock_EmailVerificationRequest_Call {
	return &UserMock_EmailVerificationRequest_Call{Call: _e.mock.On("EmailVerificationRequest", ctx)}
}

func (_c *UserMock_EmailVerificationRequest_Call) Run(run func(ctx context.Context)) *UserMock_EmailVerificationRequest_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *UserMock_EmailVerificationRequest_Call) Return(_a0 dto.TokenResponseDTO, _a1 status.Object) *UserMock_EmailVerificationRequest_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserMock_EmailVerificationRequest_Call) RunAndReturn(run func(context.Context) (dto.TokenResponseDTO, status.Object)) *UserMock_EmailVerificationRequest_Call {
	_c.Call.Return(run)
	return _c
}

// FindByIds provides a mock function with given fields: ctx, ids
func (_m *UserMock) FindByIds(ctx context.Context, ids ...types.Id) ([]dto.UserResponseDTO, status.Object) {
	_va := make([]interface{}, len(ids))
	for _i := range ids {
		_va[_i] = ids[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for FindByIds")
	}

	var r0 []dto.UserResponseDTO
	var r1 status.Object
	if rf, ok := ret.Get(0).(func(context.Context, ...types.Id) ([]dto.UserResponseDTO, status.Object)); ok {
		return rf(ctx, ids...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ...types.Id) []dto.UserResponseDTO); ok {
		r0 = rf(ctx, ids...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]dto.UserResponseDTO)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, ...types.Id) status.Object); ok {
		r1 = rf(ctx, ids...)
	} else {
		r1 = ret.Get(1).(status.Object)
	}

	return r0, r1
}

// UserMock_FindByIds_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindByIds'
type UserMock_FindByIds_Call struct {
	*mock.Call
}

// FindByIds is a helper method to define mock.On call
//   - ctx context.Context
//   - ids ...types.Id
func (_e *UserMock_Expecter) FindByIds(ctx interface{}, ids ...interface{}) *UserMock_FindByIds_Call {
	return &UserMock_FindByIds_Call{Call: _e.mock.On("FindByIds",
		append([]interface{}{ctx}, ids...)...)}
}

func (_c *UserMock_FindByIds_Call) Run(run func(ctx context.Context, ids ...types.Id)) *UserMock_FindByIds_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]types.Id, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(types.Id)
			}
		}
		run(args[0].(context.Context), variadicArgs...)
	})
	return _c
}

func (_c *UserMock_FindByIds_Call) Return(_a0 []dto.UserResponseDTO, _a1 status.Object) *UserMock_FindByIds_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserMock_FindByIds_Call) RunAndReturn(run func(context.Context, ...types.Id) ([]dto.UserResponseDTO, status.Object)) *UserMock_FindByIds_Call {
	_c.Call.Return(run)
	return _c
}

// ForgotPassword provides a mock function with given fields: ctx, email
func (_m *UserMock) ForgotPassword(ctx context.Context, email types.Email) (dto.TokenResponseDTO, status.Object) {
	ret := _m.Called(ctx, email)

	if len(ret) == 0 {
		panic("no return value specified for ForgotPassword")
	}

	var r0 dto.TokenResponseDTO
	var r1 status.Object
	if rf, ok := ret.Get(0).(func(context.Context, types.Email) (dto.TokenResponseDTO, status.Object)); ok {
		return rf(ctx, email)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.Email) dto.TokenResponseDTO); ok {
		r0 = rf(ctx, email)
	} else {
		r0 = ret.Get(0).(dto.TokenResponseDTO)
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.Email) status.Object); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Get(1).(status.Object)
	}

	return r0, r1
}

// UserMock_ForgotPassword_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ForgotPassword'
type UserMock_ForgotPassword_Call struct {
	*mock.Call
}

// ForgotPassword is a helper method to define mock.On call
//   - ctx context.Context
//   - email types.Email
func (_e *UserMock_Expecter) ForgotPassword(ctx interface{}, email interface{}) *UserMock_ForgotPassword_Call {
	return &UserMock_ForgotPassword_Call{Call: _e.mock.On("ForgotPassword", ctx, email)}
}

func (_c *UserMock_ForgotPassword_Call) Run(run func(ctx context.Context, email types.Email)) *UserMock_ForgotPassword_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(types.Email))
	})
	return _c
}

func (_c *UserMock_ForgotPassword_Call) Return(_a0 dto.TokenResponseDTO, _a1 status.Object) *UserMock_ForgotPassword_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserMock_ForgotPassword_Call) RunAndReturn(run func(context.Context, types.Email) (dto.TokenResponseDTO, status.Object)) *UserMock_ForgotPassword_Call {
	_c.Call.Return(run)
	return _c
}

// GetAll provides a mock function with given fields: ctx, pagedDto
func (_m *UserMock) GetAll(ctx context.Context, pagedDto shareddto.PagedElementDTO) (shareddto.PagedElementResult[dto.UserResponseDTO], status.Object) {
	ret := _m.Called(ctx, pagedDto)

	if len(ret) == 0 {
		panic("no return value specified for GetAll")
	}

	var r0 shareddto.PagedElementResult[dto.UserResponseDTO]
	var r1 status.Object
	if rf, ok := ret.Get(0).(func(context.Context, shareddto.PagedElementDTO) (shareddto.PagedElementResult[dto.UserResponseDTO], status.Object)); ok {
		return rf(ctx, pagedDto)
	}
	if rf, ok := ret.Get(0).(func(context.Context, shareddto.PagedElementDTO) shareddto.PagedElementResult[dto.UserResponseDTO]); ok {
		r0 = rf(ctx, pagedDto)
	} else {
		r0 = ret.Get(0).(shareddto.PagedElementResult[dto.UserResponseDTO])
	}

	if rf, ok := ret.Get(1).(func(context.Context, shareddto.PagedElementDTO) status.Object); ok {
		r1 = rf(ctx, pagedDto)
	} else {
		r1 = ret.Get(1).(status.Object)
	}

	return r0, r1
}

// UserMock_GetAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAll'
type UserMock_GetAll_Call struct {
	*mock.Call
}

// GetAll is a helper method to define mock.On call
//   - ctx context.Context
//   - pagedDto shareddto.PagedElementDTO
func (_e *UserMock_Expecter) GetAll(ctx interface{}, pagedDto interface{}) *UserMock_GetAll_Call {
	return &UserMock_GetAll_Call{Call: _e.mock.On("GetAll", ctx, pagedDto)}
}

func (_c *UserMock_GetAll_Call) Run(run func(ctx context.Context, pagedDto shareddto.PagedElementDTO)) *UserMock_GetAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(shareddto.PagedElementDTO))
	})
	return _c
}

func (_c *UserMock_GetAll_Call) Return(_a0 shareddto.PagedElementResult[dto.UserResponseDTO], _a1 status.Object) *UserMock_GetAll_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserMock_GetAll_Call) RunAndReturn(run func(context.Context, shareddto.PagedElementDTO) (shareddto.PagedElementResult[dto.UserResponseDTO], status.Object)) *UserMock_GetAll_Call {
	_c.Call.Return(run)
	return _c
}

// ResetPassword provides a mock function with given fields: ctx, input
func (_m *UserMock) ResetPassword(ctx context.Context, input *dto.ResetUserPasswordDTO) status.Object {
	ret := _m.Called(ctx, input)

	if len(ret) == 0 {
		panic("no return value specified for ResetPassword")
	}

	var r0 status.Object
	if rf, ok := ret.Get(0).(func(context.Context, *dto.ResetUserPasswordDTO) status.Object); ok {
		r0 = rf(ctx, input)
	} else {
		r0 = ret.Get(0).(status.Object)
	}

	return r0
}

// UserMock_ResetPassword_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ResetPassword'
type UserMock_ResetPassword_Call struct {
	*mock.Call
}

// ResetPassword is a helper method to define mock.On call
//   - ctx context.Context
//   - input *dto.ResetUserPasswordDTO
func (_e *UserMock_Expecter) ResetPassword(ctx interface{}, input interface{}) *UserMock_ResetPassword_Call {
	return &UserMock_ResetPassword_Call{Call: _e.mock.On("ResetPassword", ctx, input)}
}

func (_c *UserMock_ResetPassword_Call) Run(run func(ctx context.Context, input *dto.ResetUserPasswordDTO)) *UserMock_ResetPassword_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*dto.ResetUserPasswordDTO))
	})
	return _c
}

func (_c *UserMock_ResetPassword_Call) Return(_a0 status.Object) *UserMock_ResetPassword_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserMock_ResetPassword_Call) RunAndReturn(run func(context.Context, *dto.ResetUserPasswordDTO) status.Object) *UserMock_ResetPassword_Call {
	_c.Call.Return(run)
	return _c
}

// ResetPasswordWithToken provides a mock function with given fields: ctx, resetDTO
func (_m *UserMock) ResetPasswordWithToken(ctx context.Context, resetDTO *dto.ResetPasswordWithTokenDTO) status.Object {
	ret := _m.Called(ctx, resetDTO)

	if len(ret) == 0 {
		panic("no return value specified for ResetPasswordWithToken")
	}

	var r0 status.Object
	if rf, ok := ret.Get(0).(func(context.Context, *dto.ResetPasswordWithTokenDTO) status.Object); ok {
		r0 = rf(ctx, resetDTO)
	} else {
		r0 = ret.Get(0).(status.Object)
	}

	return r0
}

// UserMock_ResetPasswordWithToken_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ResetPasswordWithToken'
type UserMock_ResetPasswordWithToken_Call struct {
	*mock.Call
}

// ResetPasswordWithToken is a helper method to define mock.On call
//   - ctx context.Context
//   - resetDTO *dto.ResetPasswordWithTokenDTO
func (_e *UserMock_Expecter) ResetPasswordWithToken(ctx interface{}, resetDTO interface{}) *UserMock_ResetPasswordWithToken_Call {
	return &UserMock_ResetPasswordWithToken_Call{Call: _e.mock.On("ResetPasswordWithToken", ctx, resetDTO)}
}

func (_c *UserMock_ResetPasswordWithToken_Call) Run(run func(ctx context.Context, resetDTO *dto.ResetPasswordWithTokenDTO)) *UserMock_ResetPasswordWithToken_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*dto.ResetPasswordWithTokenDTO))
	})
	return _c
}

func (_c *UserMock_ResetPasswordWithToken_Call) Return(_a0 status.Object) *UserMock_ResetPasswordWithToken_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserMock_ResetPasswordWithToken_Call) RunAndReturn(run func(context.Context, *dto.ResetPasswordWithTokenDTO) status.Object) *UserMock_ResetPasswordWithToken_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: ctx, input
func (_m *UserMock) Update(ctx context.Context, input *dto.UserUpdateDTO) status.Object {
	ret := _m.Called(ctx, input)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 status.Object
	if rf, ok := ret.Get(0).(func(context.Context, *dto.UserUpdateDTO) status.Object); ok {
		r0 = rf(ctx, input)
	} else {
		r0 = ret.Get(0).(status.Object)
	}

	return r0
}

// UserMock_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type UserMock_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - input *dto.UserUpdateDTO
func (_e *UserMock_Expecter) Update(ctx interface{}, input interface{}) *UserMock_Update_Call {
	return &UserMock_Update_Call{Call: _e.mock.On("Update", ctx, input)}
}

func (_c *UserMock_Update_Call) Run(run func(ctx context.Context, input *dto.UserUpdateDTO)) *UserMock_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*dto.UserUpdateDTO))
	})
	return _c
}

func (_c *UserMock_Update_Call) Return(_a0 status.Object) *UserMock_Update_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserMock_Update_Call) RunAndReturn(run func(context.Context, *dto.UserUpdateDTO) status.Object) *UserMock_Update_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateAvatar provides a mock function with given fields: ctx, input
func (_m *UserMock) UpdateAvatar(ctx context.Context, input *dto.UpdateUserAvatarDTO) status.Object {
	ret := _m.Called(ctx, input)

	if len(ret) == 0 {
		panic("no return value specified for UpdateAvatar")
	}

	var r0 status.Object
	if rf, ok := ret.Get(0).(func(context.Context, *dto.UpdateUserAvatarDTO) status.Object); ok {
		r0 = rf(ctx, input)
	} else {
		r0 = ret.Get(0).(status.Object)
	}

	return r0
}

// UserMock_UpdateAvatar_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateAvatar'
type UserMock_UpdateAvatar_Call struct {
	*mock.Call
}

// UpdateAvatar is a helper method to define mock.On call
//   - ctx context.Context
//   - input *dto.UpdateUserAvatarDTO
func (_e *UserMock_Expecter) UpdateAvatar(ctx interface{}, input interface{}) *UserMock_UpdateAvatar_Call {
	return &UserMock_UpdateAvatar_Call{Call: _e.mock.On("UpdateAvatar", ctx, input)}
}

func (_c *UserMock_UpdateAvatar_Call) Run(run func(ctx context.Context, input *dto.UpdateUserAvatarDTO)) *UserMock_UpdateAvatar_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*dto.UpdateUserAvatarDTO))
	})
	return _c
}

func (_c *UserMock_UpdateAvatar_Call) Return(_a0 status.Object) *UserMock_UpdateAvatar_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserMock_UpdateAvatar_Call) RunAndReturn(run func(context.Context, *dto.UpdateUserAvatarDTO) status.Object) *UserMock_UpdateAvatar_Call {
	_c.Call.Return(run)
	return _c
}

// UpdatePassword provides a mock function with given fields: ctx, input
func (_m *UserMock) UpdatePassword(ctx context.Context, input *dto.UserUpdatePasswordDTO) status.Object {
	ret := _m.Called(ctx, input)

	if len(ret) == 0 {
		panic("no return value specified for UpdatePassword")
	}

	var r0 status.Object
	if rf, ok := ret.Get(0).(func(context.Context, *dto.UserUpdatePasswordDTO) status.Object); ok {
		r0 = rf(ctx, input)
	} else {
		r0 = ret.Get(0).(status.Object)
	}

	return r0
}

// UserMock_UpdatePassword_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdatePassword'
type UserMock_UpdatePassword_Call struct {
	*mock.Call
}

// UpdatePassword is a helper method to define mock.On call
//   - ctx context.Context
//   - input *dto.UserUpdatePasswordDTO
func (_e *UserMock_Expecter) UpdatePassword(ctx interface{}, input interface{}) *UserMock_UpdatePassword_Call {
	return &UserMock_UpdatePassword_Call{Call: _e.mock.On("UpdatePassword", ctx, input)}
}

func (_c *UserMock_UpdatePassword_Call) Run(run func(ctx context.Context, input *dto.UserUpdatePasswordDTO)) *UserMock_UpdatePassword_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*dto.UserUpdatePasswordDTO))
	})
	return _c
}

func (_c *UserMock_UpdatePassword_Call) Return(_a0 status.Object) *UserMock_UpdatePassword_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserMock_UpdatePassword_Call) RunAndReturn(run func(context.Context, *dto.UserUpdatePasswordDTO) status.Object) *UserMock_UpdatePassword_Call {
	_c.Call.Return(run)
	return _c
}

// VerifyEmail provides a mock function with given fields: ctx, token
func (_m *UserMock) VerifyEmail(ctx context.Context, token string) status.Object {
	ret := _m.Called(ctx, token)

	if len(ret) == 0 {
		panic("no return value specified for VerifyEmail")
	}

	var r0 status.Object
	if rf, ok := ret.Get(0).(func(context.Context, string) status.Object); ok {
		r0 = rf(ctx, token)
	} else {
		r0 = ret.Get(0).(status.Object)
	}

	return r0
}

// UserMock_VerifyEmail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'VerifyEmail'
type UserMock_VerifyEmail_Call struct {
	*mock.Call
}

// VerifyEmail is a helper method to define mock.On call
//   - ctx context.Context
//   - token string
func (_e *UserMock_Expecter) VerifyEmail(ctx interface{}, token interface{}) *UserMock_VerifyEmail_Call {
	return &UserMock_VerifyEmail_Call{Call: _e.mock.On("VerifyEmail", ctx, token)}
}

func (_c *UserMock_VerifyEmail_Call) Run(run func(ctx context.Context, token string)) *UserMock_VerifyEmail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *UserMock_VerifyEmail_Call) Return(_a0 status.Object) *UserMock_VerifyEmail_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserMock_VerifyEmail_Call) RunAndReturn(run func(context.Context, string) status.Object) *UserMock_VerifyEmail_Call {
	_c.Call.Return(run)
	return _c
}

// NewUserMock creates a new instance of UserMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserMock {
	mock := &UserMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
