// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"
	entity "nexa/services/authentication/internal/domain/entity"

	mock "github.com/stretchr/testify/mock"

	repo "nexa/shared/util/repo"

	types "nexa/shared/types"
)

// CredentialMock is an autogenerated mock type for the ICredential type
type CredentialMock struct {
	mock.Mock
}

type CredentialMock_Expecter struct {
	mock *mock.Mock
}

func (_m *CredentialMock) EXPECT() *CredentialMock_Expecter {
	return &CredentialMock_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: ctx, credential
func (_m *CredentialMock) Create(ctx context.Context, credential *entity.Credential) error {
	ret := _m.Called(ctx, credential)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.Credential) error); ok {
		r0 = rf(ctx, credential)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CredentialMock_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type CredentialMock_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - credential *entity.Credential
func (_e *CredentialMock_Expecter) Create(ctx interface{}, credential interface{}) *CredentialMock_Create_Call {
	return &CredentialMock_Create_Call{Call: _e.mock.On("Create", ctx, credential)}
}

func (_c *CredentialMock_Create_Call) Run(run func(ctx context.Context, credential *entity.Credential)) *CredentialMock_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*entity.Credential))
	})
	return _c
}

func (_c *CredentialMock_Create_Call) Return(_a0 error) *CredentialMock_Create_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *CredentialMock_Create_Call) RunAndReturn(run func(context.Context, *entity.Credential) error) *CredentialMock_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: ctx, credIds
func (_m *CredentialMock) Delete(ctx context.Context, credIds ...types.Id) error {
	_va := make([]interface{}, len(credIds))
	for _i := range credIds {
		_va[_i] = credIds[_i]
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
		r0 = rf(ctx, credIds...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CredentialMock_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type CredentialMock_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - credIds ...types.Id
func (_e *CredentialMock_Expecter) Delete(ctx interface{}, credIds ...interface{}) *CredentialMock_Delete_Call {
	return &CredentialMock_Delete_Call{Call: _e.mock.On("Delete",
		append([]interface{}{ctx}, credIds...)...)}
}

func (_c *CredentialMock_Delete_Call) Run(run func(ctx context.Context, credIds ...types.Id)) *CredentialMock_Delete_Call {
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

func (_c *CredentialMock_Delete_Call) Return(_a0 error) *CredentialMock_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *CredentialMock_Delete_Call) RunAndReturn(run func(context.Context, ...types.Id) error) *CredentialMock_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteByUserId provides a mock function with given fields: ctx, userId, credIds
func (_m *CredentialMock) DeleteByUserId(ctx context.Context, userId types.Id, credIds ...types.Id) error {
	_va := make([]interface{}, len(credIds))
	for _i := range credIds {
		_va[_i] = credIds[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, userId)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for DeleteByUserId")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, types.Id, ...types.Id) error); ok {
		r0 = rf(ctx, userId, credIds...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CredentialMock_DeleteByUserId_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteByUserId'
type CredentialMock_DeleteByUserId_Call struct {
	*mock.Call
}

// DeleteByUserId is a helper method to define mock.On call
//   - ctx context.Context
//   - userId types.Id
//   - credIds ...types.Id
func (_e *CredentialMock_Expecter) DeleteByUserId(ctx interface{}, userId interface{}, credIds ...interface{}) *CredentialMock_DeleteByUserId_Call {
	return &CredentialMock_DeleteByUserId_Call{Call: _e.mock.On("DeleteByUserId",
		append([]interface{}{ctx, userId}, credIds...)...)}
}

func (_c *CredentialMock_DeleteByUserId_Call) Run(run func(ctx context.Context, userId types.Id, credIds ...types.Id)) *CredentialMock_DeleteByUserId_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]types.Id, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(types.Id)
			}
		}
		run(args[0].(context.Context), args[1].(types.Id), variadicArgs...)
	})
	return _c
}

func (_c *CredentialMock_DeleteByUserId_Call) Return(_a0 error) *CredentialMock_DeleteByUserId_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *CredentialMock_DeleteByUserId_Call) RunAndReturn(run func(context.Context, types.Id, ...types.Id) error) *CredentialMock_DeleteByUserId_Call {
	_c.Call.Return(run)
	return _c
}

// Find provides a mock function with given fields: ctx, refreshTokenId
func (_m *CredentialMock) Find(ctx context.Context, refreshTokenId types.Id) (*entity.Credential, error) {
	ret := _m.Called(ctx, refreshTokenId)

	if len(ret) == 0 {
		panic("no return value specified for Find")
	}

	var r0 *entity.Credential
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.Id) (*entity.Credential, error)); ok {
		return rf(ctx, refreshTokenId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.Id) *entity.Credential); ok {
		r0 = rf(ctx, refreshTokenId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Credential)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.Id) error); ok {
		r1 = rf(ctx, refreshTokenId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CredentialMock_Find_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Find'
type CredentialMock_Find_Call struct {
	*mock.Call
}

// Find is a helper method to define mock.On call
//   - ctx context.Context
//   - refreshTokenId types.Id
func (_e *CredentialMock_Expecter) Find(ctx interface{}, refreshTokenId interface{}) *CredentialMock_Find_Call {
	return &CredentialMock_Find_Call{Call: _e.mock.On("Find", ctx, refreshTokenId)}
}

func (_c *CredentialMock_Find_Call) Run(run func(ctx context.Context, refreshTokenId types.Id)) *CredentialMock_Find_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(types.Id))
	})
	return _c
}

func (_c *CredentialMock_Find_Call) Return(_a0 *entity.Credential, _a1 error) *CredentialMock_Find_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CredentialMock_Find_Call) RunAndReturn(run func(context.Context, types.Id) (*entity.Credential, error)) *CredentialMock_Find_Call {
	_c.Call.Return(run)
	return _c
}

// FindByUserId provides a mock function with given fields: ctx, userId
func (_m *CredentialMock) FindByUserId(ctx context.Context, userId types.Id) ([]entity.Credential, error) {
	ret := _m.Called(ctx, userId)

	if len(ret) == 0 {
		panic("no return value specified for FindByUserId")
	}

	var r0 []entity.Credential
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.Id) ([]entity.Credential, error)); ok {
		return rf(ctx, userId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.Id) []entity.Credential); ok {
		r0 = rf(ctx, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.Credential)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.Id) error); ok {
		r1 = rf(ctx, userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CredentialMock_FindByUserId_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindByUserId'
type CredentialMock_FindByUserId_Call struct {
	*mock.Call
}

// FindByUserId is a helper method to define mock.On call
//   - ctx context.Context
//   - userId types.Id
func (_e *CredentialMock_Expecter) FindByUserId(ctx interface{}, userId interface{}) *CredentialMock_FindByUserId_Call {
	return &CredentialMock_FindByUserId_Call{Call: _e.mock.On("FindByUserId", ctx, userId)}
}

func (_c *CredentialMock_FindByUserId_Call) Run(run func(ctx context.Context, userId types.Id)) *CredentialMock_FindByUserId_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(types.Id))
	})
	return _c
}

func (_c *CredentialMock_FindByUserId_Call) Return(_a0 []entity.Credential, _a1 error) *CredentialMock_FindByUserId_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CredentialMock_FindByUserId_Call) RunAndReturn(run func(context.Context, types.Id) ([]entity.Credential, error)) *CredentialMock_FindByUserId_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: ctx, parameter
func (_m *CredentialMock) Get(ctx context.Context, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Credential], error) {
	ret := _m.Called(ctx, parameter)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 repo.PaginatedResult[entity.Credential]
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, repo.QueryParameter) (repo.PaginatedResult[entity.Credential], error)); ok {
		return rf(ctx, parameter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, repo.QueryParameter) repo.PaginatedResult[entity.Credential]); ok {
		r0 = rf(ctx, parameter)
	} else {
		r0 = ret.Get(0).(repo.PaginatedResult[entity.Credential])
	}

	if rf, ok := ret.Get(1).(func(context.Context, repo.QueryParameter) error); ok {
		r1 = rf(ctx, parameter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CredentialMock_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type CredentialMock_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - parameter repo.QueryParameter
func (_e *CredentialMock_Expecter) Get(ctx interface{}, parameter interface{}) *CredentialMock_Get_Call {
	return &CredentialMock_Get_Call{Call: _e.mock.On("Get", ctx, parameter)}
}

func (_c *CredentialMock_Get_Call) Run(run func(ctx context.Context, parameter repo.QueryParameter)) *CredentialMock_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(repo.QueryParameter))
	})
	return _c
}

func (_c *CredentialMock_Get_Call) Return(_a0 repo.PaginatedResult[entity.Credential], _a1 error) *CredentialMock_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CredentialMock_Get_Call) RunAndReturn(run func(context.Context, repo.QueryParameter) (repo.PaginatedResult[entity.Credential], error)) *CredentialMock_Get_Call {
	_c.Call.Return(run)
	return _c
}

// Patch provides a mock function with given fields: ctx, credential
func (_m *CredentialMock) Patch(ctx context.Context, credential *entity.Credential) error {
	ret := _m.Called(ctx, credential)

	if len(ret) == 0 {
		panic("no return value specified for Patch")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.Credential) error); ok {
		r0 = rf(ctx, credential)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CredentialMock_Patch_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Patch'
type CredentialMock_Patch_Call struct {
	*mock.Call
}

// Patch is a helper method to define mock.On call
//   - ctx context.Context
//   - credential *entity.Credential
func (_e *CredentialMock_Expecter) Patch(ctx interface{}, credential interface{}) *CredentialMock_Patch_Call {
	return &CredentialMock_Patch_Call{Call: _e.mock.On("Patch", ctx, credential)}
}

func (_c *CredentialMock_Patch_Call) Run(run func(ctx context.Context, credential *entity.Credential)) *CredentialMock_Patch_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*entity.Credential))
	})
	return _c
}

func (_c *CredentialMock_Patch_Call) Return(_a0 error) *CredentialMock_Patch_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *CredentialMock_Patch_Call) RunAndReturn(run func(context.Context, *entity.Credential) error) *CredentialMock_Patch_Call {
	_c.Call.Return(run)
	return _c
}

// NewCredentialMock creates a new instance of CredentialMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCredentialMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *CredentialMock {
	mock := &CredentialMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}