// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"
	entity "nexa/services/authorization/internal/domain/entity"

	mock "github.com/stretchr/testify/mock"

	repo "nexa/shared/util/repo"

	types "nexa/shared/types"
)

// RoleMock is an autogenerated mock type for the IRole type
type RoleMock struct {
	mock.Mock
}

type RoleMock_Expecter struct {
	mock *mock.Mock
}

func (_m *RoleMock) EXPECT() *RoleMock_Expecter {
	return &RoleMock_Expecter{mock: &_m.Mock}
}

// AddPermissions provides a mock function with given fields: ctx, roleId, permissionIds
func (_m *RoleMock) AddPermissions(ctx context.Context, roleId types.Id, permissionIds ...types.Id) error {
	_va := make([]interface{}, len(permissionIds))
	for _i := range permissionIds {
		_va[_i] = permissionIds[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, roleId)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for AddPermissions")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, types.Id, ...types.Id) error); ok {
		r0 = rf(ctx, roleId, permissionIds...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RoleMock_AddPermissions_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddPermissions'
type RoleMock_AddPermissions_Call struct {
	*mock.Call
}

// AddPermissions is a helper method to define mock.On call
//   - ctx context.Context
//   - roleId types.Id
//   - permissionIds ...types.Id
func (_e *RoleMock_Expecter) AddPermissions(ctx interface{}, roleId interface{}, permissionIds ...interface{}) *RoleMock_AddPermissions_Call {
	return &RoleMock_AddPermissions_Call{Call: _e.mock.On("AddPermissions",
		append([]interface{}{ctx, roleId}, permissionIds...)...)}
}

func (_c *RoleMock_AddPermissions_Call) Run(run func(ctx context.Context, roleId types.Id, permissionIds ...types.Id)) *RoleMock_AddPermissions_Call {
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

func (_c *RoleMock_AddPermissions_Call) Return(_a0 error) *RoleMock_AddPermissions_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *RoleMock_AddPermissions_Call) RunAndReturn(run func(context.Context, types.Id, ...types.Id) error) *RoleMock_AddPermissions_Call {
	_c.Call.Return(run)
	return _c
}

// AddUser provides a mock function with given fields: ctx, userId, roleIds
func (_m *RoleMock) AddUser(ctx context.Context, userId types.Id, roleIds ...types.Id) error {
	_va := make([]interface{}, len(roleIds))
	for _i := range roleIds {
		_va[_i] = roleIds[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, userId)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for AddUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, types.Id, ...types.Id) error); ok {
		r0 = rf(ctx, userId, roleIds...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RoleMock_AddUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddUser'
type RoleMock_AddUser_Call struct {
	*mock.Call
}

// AddUser is a helper method to define mock.On call
//   - ctx context.Context
//   - userId types.Id
//   - roleIds ...types.Id
func (_e *RoleMock_Expecter) AddUser(ctx interface{}, userId interface{}, roleIds ...interface{}) *RoleMock_AddUser_Call {
	return &RoleMock_AddUser_Call{Call: _e.mock.On("AddUser",
		append([]interface{}{ctx, userId}, roleIds...)...)}
}

func (_c *RoleMock_AddUser_Call) Run(run func(ctx context.Context, userId types.Id, roleIds ...types.Id)) *RoleMock_AddUser_Call {
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

func (_c *RoleMock_AddUser_Call) Return(_a0 error) *RoleMock_AddUser_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *RoleMock_AddUser_Call) RunAndReturn(run func(context.Context, types.Id, ...types.Id) error) *RoleMock_AddUser_Call {
	_c.Call.Return(run)
	return _c
}

// Create provides a mock function with given fields: ctx, role
func (_m *RoleMock) Create(ctx context.Context, role *entity.Role) error {
	ret := _m.Called(ctx, role)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.Role) error); ok {
		r0 = rf(ctx, role)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RoleMock_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type RoleMock_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - role *entity.Role
func (_e *RoleMock_Expecter) Create(ctx interface{}, role interface{}) *RoleMock_Create_Call {
	return &RoleMock_Create_Call{Call: _e.mock.On("Create", ctx, role)}
}

func (_c *RoleMock_Create_Call) Run(run func(ctx context.Context, role *entity.Role)) *RoleMock_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*entity.Role))
	})
	return _c
}

func (_c *RoleMock_Create_Call) Return(_a0 error) *RoleMock_Create_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *RoleMock_Create_Call) RunAndReturn(run func(context.Context, *entity.Role) error) *RoleMock_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: ctx, id
func (_m *RoleMock) Delete(ctx context.Context, id types.Id) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, types.Id) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RoleMock_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type RoleMock_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - id types.Id
func (_e *RoleMock_Expecter) Delete(ctx interface{}, id interface{}) *RoleMock_Delete_Call {
	return &RoleMock_Delete_Call{Call: _e.mock.On("Delete", ctx, id)}
}

func (_c *RoleMock_Delete_Call) Run(run func(ctx context.Context, id types.Id)) *RoleMock_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(types.Id))
	})
	return _c
}

func (_c *RoleMock_Delete_Call) Return(_a0 error) *RoleMock_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *RoleMock_Delete_Call) RunAndReturn(run func(context.Context, types.Id) error) *RoleMock_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// FindByIds provides a mock function with given fields: ctx, ids
func (_m *RoleMock) FindByIds(ctx context.Context, ids ...types.Id) ([]entity.Role, error) {
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

	var r0 []entity.Role
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ...types.Id) ([]entity.Role, error)); ok {
		return rf(ctx, ids...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ...types.Id) []entity.Role); ok {
		r0 = rf(ctx, ids...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.Role)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, ...types.Id) error); ok {
		r1 = rf(ctx, ids...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RoleMock_FindByIds_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindByIds'
type RoleMock_FindByIds_Call struct {
	*mock.Call
}

// FindByIds is a helper method to define mock.On call
//   - ctx context.Context
//   - ids ...types.Id
func (_e *RoleMock_Expecter) FindByIds(ctx interface{}, ids ...interface{}) *RoleMock_FindByIds_Call {
	return &RoleMock_FindByIds_Call{Call: _e.mock.On("FindByIds",
		append([]interface{}{ctx}, ids...)...)}
}

func (_c *RoleMock_FindByIds_Call) Run(run func(ctx context.Context, ids ...types.Id)) *RoleMock_FindByIds_Call {
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

func (_c *RoleMock_FindByIds_Call) Return(_a0 []entity.Role, _a1 error) *RoleMock_FindByIds_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *RoleMock_FindByIds_Call) RunAndReturn(run func(context.Context, ...types.Id) ([]entity.Role, error)) *RoleMock_FindByIds_Call {
	_c.Call.Return(run)
	return _c
}

// FindByName provides a mock function with given fields: ctx, name
func (_m *RoleMock) FindByName(ctx context.Context, name string) (entity.Role, error) {
	ret := _m.Called(ctx, name)

	if len(ret) == 0 {
		panic("no return value specified for FindByName")
	}

	var r0 entity.Role
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (entity.Role, error)); ok {
		return rf(ctx, name)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) entity.Role); ok {
		r0 = rf(ctx, name)
	} else {
		r0 = ret.Get(0).(entity.Role)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RoleMock_FindByName_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindByName'
type RoleMock_FindByName_Call struct {
	*mock.Call
}

// FindByName is a helper method to define mock.On call
//   - ctx context.Context
//   - name string
func (_e *RoleMock_Expecter) FindByName(ctx interface{}, name interface{}) *RoleMock_FindByName_Call {
	return &RoleMock_FindByName_Call{Call: _e.mock.On("FindByName", ctx, name)}
}

func (_c *RoleMock_FindByName_Call) Run(run func(ctx context.Context, name string)) *RoleMock_FindByName_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *RoleMock_FindByName_Call) Return(_a0 entity.Role, _a1 error) *RoleMock_FindByName_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *RoleMock_FindByName_Call) RunAndReturn(run func(context.Context, string) (entity.Role, error)) *RoleMock_FindByName_Call {
	_c.Call.Return(run)
	return _c
}

// FindByUserId provides a mock function with given fields: ctx, userId
func (_m *RoleMock) FindByUserId(ctx context.Context, userId types.Id) ([]entity.Role, error) {
	ret := _m.Called(ctx, userId)

	if len(ret) == 0 {
		panic("no return value specified for FindByUserId")
	}

	var r0 []entity.Role
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.Id) ([]entity.Role, error)); ok {
		return rf(ctx, userId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.Id) []entity.Role); ok {
		r0 = rf(ctx, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.Role)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.Id) error); ok {
		r1 = rf(ctx, userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RoleMock_FindByUserId_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindByUserId'
type RoleMock_FindByUserId_Call struct {
	*mock.Call
}

// FindByUserId is a helper method to define mock.On call
//   - ctx context.Context
//   - userId types.Id
func (_e *RoleMock_Expecter) FindByUserId(ctx interface{}, userId interface{}) *RoleMock_FindByUserId_Call {
	return &RoleMock_FindByUserId_Call{Call: _e.mock.On("FindByUserId", ctx, userId)}
}

func (_c *RoleMock_FindByUserId_Call) Run(run func(ctx context.Context, userId types.Id)) *RoleMock_FindByUserId_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(types.Id))
	})
	return _c
}

func (_c *RoleMock_FindByUserId_Call) Return(_a0 []entity.Role, _a1 error) *RoleMock_FindByUserId_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *RoleMock_FindByUserId_Call) RunAndReturn(run func(context.Context, types.Id) ([]entity.Role, error)) *RoleMock_FindByUserId_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: ctx, parameter
func (_m *RoleMock) Get(ctx context.Context, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Role], error) {
	ret := _m.Called(ctx, parameter)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 repo.PaginatedResult[entity.Role]
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, repo.QueryParameter) (repo.PaginatedResult[entity.Role], error)); ok {
		return rf(ctx, parameter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, repo.QueryParameter) repo.PaginatedResult[entity.Role]); ok {
		r0 = rf(ctx, parameter)
	} else {
		r0 = ret.Get(0).(repo.PaginatedResult[entity.Role])
	}

	if rf, ok := ret.Get(1).(func(context.Context, repo.QueryParameter) error); ok {
		r1 = rf(ctx, parameter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RoleMock_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type RoleMock_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - parameter repo.QueryParameter
func (_e *RoleMock_Expecter) Get(ctx interface{}, parameter interface{}) *RoleMock_Get_Call {
	return &RoleMock_Get_Call{Call: _e.mock.On("Get", ctx, parameter)}
}

func (_c *RoleMock_Get_Call) Run(run func(ctx context.Context, parameter repo.QueryParameter)) *RoleMock_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(repo.QueryParameter))
	})
	return _c
}

func (_c *RoleMock_Get_Call) Return(_a0 repo.PaginatedResult[entity.Role], _a1 error) *RoleMock_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *RoleMock_Get_Call) RunAndReturn(run func(context.Context, repo.QueryParameter) (repo.PaginatedResult[entity.Role], error)) *RoleMock_Get_Call {
	_c.Call.Return(run)
	return _c
}

// Patch provides a mock function with given fields: ctx, role
func (_m *RoleMock) Patch(ctx context.Context, role *entity.PatchedRole) error {
	ret := _m.Called(ctx, role)

	if len(ret) == 0 {
		panic("no return value specified for Patch")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.PatchedRole) error); ok {
		r0 = rf(ctx, role)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RoleMock_Patch_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Patch'
type RoleMock_Patch_Call struct {
	*mock.Call
}

// Patch is a helper method to define mock.On call
//   - ctx context.Context
//   - role *entity.PatchedRole
func (_e *RoleMock_Expecter) Patch(ctx interface{}, role interface{}) *RoleMock_Patch_Call {
	return &RoleMock_Patch_Call{Call: _e.mock.On("Patch", ctx, role)}
}

func (_c *RoleMock_Patch_Call) Run(run func(ctx context.Context, role *entity.PatchedRole)) *RoleMock_Patch_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*entity.PatchedRole))
	})
	return _c
}

func (_c *RoleMock_Patch_Call) Return(_a0 error) *RoleMock_Patch_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *RoleMock_Patch_Call) RunAndReturn(run func(context.Context, *entity.PatchedRole) error) *RoleMock_Patch_Call {
	_c.Call.Return(run)
	return _c
}

// RemovePermissions provides a mock function with given fields: ctx, roleId, permissionIds
func (_m *RoleMock) RemovePermissions(ctx context.Context, roleId types.Id, permissionIds ...types.Id) error {
	_va := make([]interface{}, len(permissionIds))
	for _i := range permissionIds {
		_va[_i] = permissionIds[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, roleId)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for RemovePermissions")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, types.Id, ...types.Id) error); ok {
		r0 = rf(ctx, roleId, permissionIds...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RoleMock_RemovePermissions_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RemovePermissions'
type RoleMock_RemovePermissions_Call struct {
	*mock.Call
}

// RemovePermissions is a helper method to define mock.On call
//   - ctx context.Context
//   - roleId types.Id
//   - permissionIds ...types.Id
func (_e *RoleMock_Expecter) RemovePermissions(ctx interface{}, roleId interface{}, permissionIds ...interface{}) *RoleMock_RemovePermissions_Call {
	return &RoleMock_RemovePermissions_Call{Call: _e.mock.On("RemovePermissions",
		append([]interface{}{ctx, roleId}, permissionIds...)...)}
}

func (_c *RoleMock_RemovePermissions_Call) Run(run func(ctx context.Context, roleId types.Id, permissionIds ...types.Id)) *RoleMock_RemovePermissions_Call {
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

func (_c *RoleMock_RemovePermissions_Call) Return(_a0 error) *RoleMock_RemovePermissions_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *RoleMock_RemovePermissions_Call) RunAndReturn(run func(context.Context, types.Id, ...types.Id) error) *RoleMock_RemovePermissions_Call {
	_c.Call.Return(run)
	return _c
}

// RemoveUser provides a mock function with given fields: ctx, userId, roleIds
func (_m *RoleMock) RemoveUser(ctx context.Context, userId types.Id, roleIds ...types.Id) error {
	_va := make([]interface{}, len(roleIds))
	for _i := range roleIds {
		_va[_i] = roleIds[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, userId)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for RemoveUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, types.Id, ...types.Id) error); ok {
		r0 = rf(ctx, userId, roleIds...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RoleMock_RemoveUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RemoveUser'
type RoleMock_RemoveUser_Call struct {
	*mock.Call
}

// RemoveUser is a helper method to define mock.On call
//   - ctx context.Context
//   - userId types.Id
//   - roleIds ...types.Id
func (_e *RoleMock_Expecter) RemoveUser(ctx interface{}, userId interface{}, roleIds ...interface{}) *RoleMock_RemoveUser_Call {
	return &RoleMock_RemoveUser_Call{Call: _e.mock.On("RemoveUser",
		append([]interface{}{ctx, userId}, roleIds...)...)}
}

func (_c *RoleMock_RemoveUser_Call) Run(run func(ctx context.Context, userId types.Id, roleIds ...types.Id)) *RoleMock_RemoveUser_Call {
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

func (_c *RoleMock_RemoveUser_Call) Return(_a0 error) *RoleMock_RemoveUser_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *RoleMock_RemoveUser_Call) RunAndReturn(run func(context.Context, types.Id, ...types.Id) error) *RoleMock_RemoveUser_Call {
	_c.Call.Return(run)
	return _c
}

// NewRoleMock creates a new instance of RoleMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRoleMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *RoleMock {
	mock := &RoleMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
