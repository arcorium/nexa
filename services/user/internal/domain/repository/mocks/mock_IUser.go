// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"
	entity "nexa/services/user/internal/domain/entity"

	mock "github.com/stretchr/testify/mock"

	repo "nexa/shared/util/repo"

	types "nexa/shared/types"
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

// Create provides a mock function with given fields: ctx, user
func (_m *UserMock) Create(ctx context.Context, user *entity.User) error {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.User) error); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UserMock_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type UserMock_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - user *entity.User
func (_e *UserMock_Expecter) Create(ctx interface{}, user interface{}) *UserMock_Create_Call {
	return &UserMock_Create_Call{Call: _e.mock.On("Create", ctx, user)}
}

func (_c *UserMock_Create_Call) Run(run func(ctx context.Context, user *entity.User)) *UserMock_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*entity.User))
	})
	return _c
}

func (_c *UserMock_Create_Call) Return(_a0 error) *UserMock_Create_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserMock_Create_Call) RunAndReturn(run func(context.Context, *entity.User) error) *UserMock_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: ctx, ids
func (_m *UserMock) Delete(ctx context.Context, ids ...types.Id) error {
	_va := make([]interface{}, len(ids))
	for _i := range ids {
		_va[_i] = ids[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ...types.Id) error); ok {
		r0 = rf(ctx, ids...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UserMock_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type UserMock_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - ids ...types.Id
func (_e *UserMock_Expecter) Delete(ctx interface{}, ids ...interface{}) *UserMock_Delete_Call {
	return &UserMock_Delete_Call{Call: _e.mock.On("Delete",
		append([]interface{}{ctx}, ids...)...)}
}

func (_c *UserMock_Delete_Call) Run(run func(ctx context.Context, ids ...types.Id)) *UserMock_Delete_Call {
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

func (_c *UserMock_Delete_Call) Return(_a0 error) *UserMock_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserMock_Delete_Call) RunAndReturn(run func(context.Context, ...types.Id) error) *UserMock_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// FindByEmails provides a mock function with given fields: ctx, emails
func (_m *UserMock) FindByEmails(ctx context.Context, emails ...types.Email) ([]entity.User, error) {
	_va := make([]interface{}, len(emails))
	for _i := range emails {
		_va[_i] = emails[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for FindByEmails")
	}

	var r0 []entity.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ...types.Email) ([]entity.User, error)); ok {
		return rf(ctx, emails...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ...types.Email) []entity.User); ok {
		r0 = rf(ctx, emails...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, ...types.Email) error); ok {
		r1 = rf(ctx, emails...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserMock_FindByEmails_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindByEmails'
type UserMock_FindByEmails_Call struct {
	*mock.Call
}

// FindByEmails is a helper method to define mock.On call
//   - ctx context.Context
//   - emails ...types.Email
func (_e *UserMock_Expecter) FindByEmails(ctx interface{}, emails ...interface{}) *UserMock_FindByEmails_Call {
	return &UserMock_FindByEmails_Call{Call: _e.mock.On("FindByEmails",
		append([]interface{}{ctx}, emails...)...)}
}

func (_c *UserMock_FindByEmails_Call) Run(run func(ctx context.Context, emails ...types.Email)) *UserMock_FindByEmails_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]types.Email, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(types.Email)
			}
		}
		run(args[0].(context.Context), variadicArgs...)
	})
	return _c
}

func (_c *UserMock_FindByEmails_Call) Return(_a0 []entity.User, _a1 error) *UserMock_FindByEmails_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserMock_FindByEmails_Call) RunAndReturn(run func(context.Context, ...types.Email) ([]entity.User, error)) *UserMock_FindByEmails_Call {
	_c.Call.Return(run)
	return _c
}

// FindByIds provides a mock function with given fields: ctx, userIds
func (_m *UserMock) FindByIds(ctx context.Context, userIds ...types.Id) ([]entity.User, error) {
	_va := make([]interface{}, len(userIds))
	for _i := range userIds {
		_va[_i] = userIds[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for FindByIds")
	}

	var r0 []entity.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ...types.Id) ([]entity.User, error)); ok {
		return rf(ctx, userIds...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ...types.Id) []entity.User); ok {
		r0 = rf(ctx, userIds...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, ...types.Id) error); ok {
		r1 = rf(ctx, userIds...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserMock_FindByIds_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindByIds'
type UserMock_FindByIds_Call struct {
	*mock.Call
}

// FindByIds is a helper method to define mock.On call
//   - ctx context.Context
//   - userIds ...types.Id
func (_e *UserMock_Expecter) FindByIds(ctx interface{}, userIds ...interface{}) *UserMock_FindByIds_Call {
	return &UserMock_FindByIds_Call{Call: _e.mock.On("FindByIds",
		append([]interface{}{ctx}, userIds...)...)}
}

func (_c *UserMock_FindByIds_Call) Run(run func(ctx context.Context, userIds ...types.Id)) *UserMock_FindByIds_Call {
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

func (_c *UserMock_FindByIds_Call) Return(_a0 []entity.User, _a1 error) *UserMock_FindByIds_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserMock_FindByIds_Call) RunAndReturn(run func(context.Context, ...types.Id) ([]entity.User, error)) *UserMock_FindByIds_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: ctx, query
func (_m *UserMock) Get(ctx context.Context, query repo.QueryParameter) (repo.PaginatedResult[entity.User], error) {
	ret := _m.Called(ctx, query)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 repo.PaginatedResult[entity.User]
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, repo.QueryParameter) (repo.PaginatedResult[entity.User], error)); ok {
		return rf(ctx, query)
	}
	if rf, ok := ret.Get(0).(func(context.Context, repo.QueryParameter) repo.PaginatedResult[entity.User]); ok {
		r0 = rf(ctx, query)
	} else {
		r0 = ret.Get(0).(repo.PaginatedResult[entity.User])
	}

	if rf, ok := ret.Get(1).(func(context.Context, repo.QueryParameter) error); ok {
		r1 = rf(ctx, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserMock_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type UserMock_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - query repo.QueryParameter
func (_e *UserMock_Expecter) Get(ctx interface{}, query interface{}) *UserMock_Get_Call {
	return &UserMock_Get_Call{Call: _e.mock.On("Get", ctx, query)}
}

func (_c *UserMock_Get_Call) Run(run func(ctx context.Context, query repo.QueryParameter)) *UserMock_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(repo.QueryParameter))
	})
	return _c
}

func (_c *UserMock_Get_Call) Return(_a0 repo.PaginatedResult[entity.User], _a1 error) *UserMock_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserMock_Get_Call) RunAndReturn(run func(context.Context, repo.QueryParameter) (repo.PaginatedResult[entity.User], error)) *UserMock_Get_Call {
	_c.Call.Return(run)
	return _c
}

// Patch provides a mock function with given fields: ctx, user
func (_m *UserMock) Patch(ctx context.Context, user *entity.PatchedUser) error {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for Patch")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.PatchedUser) error); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UserMock_Patch_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Patch'
type UserMock_Patch_Call struct {
	*mock.Call
}

// Patch is a helper method to define mock.On call
//   - ctx context.Context
//   - user *entity.PatchedUser
func (_e *UserMock_Expecter) Patch(ctx interface{}, user interface{}) *UserMock_Patch_Call {
	return &UserMock_Patch_Call{Call: _e.mock.On("Patch", ctx, user)}
}

func (_c *UserMock_Patch_Call) Run(run func(ctx context.Context, user *entity.PatchedUser)) *UserMock_Patch_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*entity.PatchedUser))
	})
	return _c
}

func (_c *UserMock_Patch_Call) Return(_a0 error) *UserMock_Patch_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserMock_Patch_Call) RunAndReturn(run func(context.Context, *entity.PatchedUser) error) *UserMock_Patch_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: ctx, user
func (_m *UserMock) Update(ctx context.Context, user *entity.User) error {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.User) error); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UserMock_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type UserMock_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - user *entity.User
func (_e *UserMock_Expecter) Update(ctx interface{}, user interface{}) *UserMock_Update_Call {
	return &UserMock_Update_Call{Call: _e.mock.On("Update", ctx, user)}
}

func (_c *UserMock_Update_Call) Run(run func(ctx context.Context, user *entity.User)) *UserMock_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*entity.User))
	})
	return _c
}

func (_c *UserMock_Update_Call) Return(_a0 error) *UserMock_Update_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UserMock_Update_Call) RunAndReturn(run func(context.Context, *entity.User) error) *UserMock_Update_Call {
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