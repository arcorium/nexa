// Code generated by mockery v2.43.2. DO NOT EDIT.

package repository

import (
	context "context"
	entity "nexa/services/mailer/internal/domain/entity"

	mock "github.com/stretchr/testify/mock"

	repo "nexa/shared/util/repo"

	types "nexa/shared/types"
)

// MailMock is an autogenerated mock type for the IMail type
type MailMock struct {
	mock.Mock
}

type MailMock_Expecter struct {
	mock *mock.Mock
}

func (_m *MailMock) EXPECT() *MailMock_Expecter {
	return &MailMock_Expecter{mock: &_m.Mock}
}

// AppendMultipleTags provides a mock function with given fields: ctx, mailTags
func (_m *MailMock) AppendMultipleTags(ctx context.Context, mailTags ...types.Pair[types.Id, []types.Id]) error {
	_va := make([]interface{}, len(mailTags))
	for _i := range mailTags {
		_va[_i] = mailTags[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for AppendMultipleTags")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ...types.Pair[types.Id, []types.Id]) error); ok {
		r0 = rf(ctx, mailTags...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MailMock_AppendMultipleTags_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AppendMultipleTags'
type MailMock_AppendMultipleTags_Call struct {
	*mock.Call
}

// AppendMultipleTags is a helper method to define mock.On call
//   - ctx context.Context
//   - mailTags ...types.Pair[types.Id,[]types.Id]
func (_e *MailMock_Expecter) AppendMultipleTags(ctx interface{}, mailTags ...interface{}) *MailMock_AppendMultipleTags_Call {
	return &MailMock_AppendMultipleTags_Call{Call: _e.mock.On("AppendMultipleTags",
		append([]interface{}{ctx}, mailTags...)...)}
}

func (_c *MailMock_AppendMultipleTags_Call) Run(run func(ctx context.Context, mailTags ...types.Pair[types.Id, []types.Id])) *MailMock_AppendMultipleTags_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]types.Pair[types.Id, []types.Id], len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(types.Pair[types.Id, []types.Id])
			}
		}
		run(args[0].(context.Context), variadicArgs...)
	})
	return _c
}

func (_c *MailMock_AppendMultipleTags_Call) Return(_a0 error) *MailMock_AppendMultipleTags_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MailMock_AppendMultipleTags_Call) RunAndReturn(run func(context.Context, ...types.Pair[types.Id, []types.Id]) error) *MailMock_AppendMultipleTags_Call {
	_c.Call.Return(run)
	return _c
}

// AppendTags provides a mock function with given fields: ctx, mailId, tagIds
func (_m *MailMock) AppendTags(ctx context.Context, mailId types.Id, tagIds []types.Id) error {
	ret := _m.Called(ctx, mailId, tagIds)

	if len(ret) == 0 {
		panic("no return value specified for AppendTags")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, types.Id, []types.Id) error); ok {
		r0 = rf(ctx, mailId, tagIds)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MailMock_AppendTags_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AppendTags'
type MailMock_AppendTags_Call struct {
	*mock.Call
}

// AppendTags is a helper method to define mock.On call
//   - ctx context.Context
//   - mailId types.Id
//   - tagIds []types.Id
func (_e *MailMock_Expecter) AppendTags(ctx interface{}, mailId interface{}, tagIds interface{}) *MailMock_AppendTags_Call {
	return &MailMock_AppendTags_Call{Call: _e.mock.On("AppendTags", ctx, mailId, tagIds)}
}

func (_c *MailMock_AppendTags_Call) Run(run func(ctx context.Context, mailId types.Id, tagIds []types.Id)) *MailMock_AppendTags_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(types.Id), args[2].([]types.Id))
	})
	return _c
}

func (_c *MailMock_AppendTags_Call) Return(_a0 error) *MailMock_AppendTags_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MailMock_AppendTags_Call) RunAndReturn(run func(context.Context, types.Id, []types.Id) error) *MailMock_AppendTags_Call {
	_c.Call.Return(run)
	return _c
}

// Create provides a mock function with given fields: ctx, mails
func (_m *MailMock) Create(ctx context.Context, mails ...entity.Mail) error {
	_va := make([]interface{}, len(mails))
	for _i := range mails {
		_va[_i] = mails[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ...entity.Mail) error); ok {
		r0 = rf(ctx, mails...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MailMock_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type MailMock_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - mails ...entity.Mail
func (_e *MailMock_Expecter) Create(ctx interface{}, mails ...interface{}) *MailMock_Create_Call {
	return &MailMock_Create_Call{Call: _e.mock.On("Create",
		append([]interface{}{ctx}, mails...)...)}
}

func (_c *MailMock_Create_Call) Run(run func(ctx context.Context, mails ...entity.Mail)) *MailMock_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]entity.Mail, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(entity.Mail)
			}
		}
		run(args[0].(context.Context), variadicArgs...)
	})
	return _c
}

func (_c *MailMock_Create_Call) Return(_a0 error) *MailMock_Create_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MailMock_Create_Call) RunAndReturn(run func(context.Context, ...entity.Mail) error) *MailMock_Create_Call {
	_c.Call.Return(run)
	return _c
}

// FindAll provides a mock function with given fields: ctx, query
func (_m *MailMock) FindAll(ctx context.Context, query repo.QueryParameter) (repo.PaginatedResult[entity.Mail], error) {
	ret := _m.Called(ctx, query)

	if len(ret) == 0 {
		panic("no return value specified for FindAll")
	}

	var r0 repo.PaginatedResult[entity.Mail]
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, repo.QueryParameter) (repo.PaginatedResult[entity.Mail], error)); ok {
		return rf(ctx, query)
	}
	if rf, ok := ret.Get(0).(func(context.Context, repo.QueryParameter) repo.PaginatedResult[entity.Mail]); ok {
		r0 = rf(ctx, query)
	} else {
		r0 = ret.Get(0).(repo.PaginatedResult[entity.Mail])
	}

	if rf, ok := ret.Get(1).(func(context.Context, repo.QueryParameter) error); ok {
		r1 = rf(ctx, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MailMock_FindAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindAll'
type MailMock_FindAll_Call struct {
	*mock.Call
}

// FindAll is a helper method to define mock.On call
//   - ctx context.Context
//   - query repo.QueryParameter
func (_e *MailMock_Expecter) FindAll(ctx interface{}, query interface{}) *MailMock_FindAll_Call {
	return &MailMock_FindAll_Call{Call: _e.mock.On("FindAll", ctx, query)}
}

func (_c *MailMock_FindAll_Call) Run(run func(ctx context.Context, query repo.QueryParameter)) *MailMock_FindAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(repo.QueryParameter))
	})
	return _c
}

func (_c *MailMock_FindAll_Call) Return(_a0 repo.PaginatedResult[entity.Mail], _a1 error) *MailMock_FindAll_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MailMock_FindAll_Call) RunAndReturn(run func(context.Context, repo.QueryParameter) (repo.PaginatedResult[entity.Mail], error)) *MailMock_FindAll_Call {
	_c.Call.Return(run)
	return _c
}

// FindByIds provides a mock function with given fields: ctx, ids
func (_m *MailMock) FindByIds(ctx context.Context, ids ...types.Id) ([]entity.Mail, error) {
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

	var r0 []entity.Mail
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ...types.Id) ([]entity.Mail, error)); ok {
		return rf(ctx, ids...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ...types.Id) []entity.Mail); ok {
		r0 = rf(ctx, ids...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.Mail)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, ...types.Id) error); ok {
		r1 = rf(ctx, ids...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MailMock_FindByIds_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindByIds'
type MailMock_FindByIds_Call struct {
	*mock.Call
}

// FindByIds is a helper method to define mock.On call
//   - ctx context.Context
//   - ids ...types.Id
func (_e *MailMock_Expecter) FindByIds(ctx interface{}, ids ...interface{}) *MailMock_FindByIds_Call {
	return &MailMock_FindByIds_Call{Call: _e.mock.On("FindByIds",
		append([]interface{}{ctx}, ids...)...)}
}

func (_c *MailMock_FindByIds_Call) Run(run func(ctx context.Context, ids ...types.Id)) *MailMock_FindByIds_Call {
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

func (_c *MailMock_FindByIds_Call) Return(_a0 []entity.Mail, _a1 error) *MailMock_FindByIds_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MailMock_FindByIds_Call) RunAndReturn(run func(context.Context, ...types.Id) ([]entity.Mail, error)) *MailMock_FindByIds_Call {
	_c.Call.Return(run)
	return _c
}

// FindByTag provides a mock function with given fields: ctx, tag
func (_m *MailMock) FindByTag(ctx context.Context, tag types.Id) ([]entity.Mail, error) {
	ret := _m.Called(ctx, tag)

	if len(ret) == 0 {
		panic("no return value specified for FindByTag")
	}

	var r0 []entity.Mail
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.Id) ([]entity.Mail, error)); ok {
		return rf(ctx, tag)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.Id) []entity.Mail); ok {
		r0 = rf(ctx, tag)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.Mail)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.Id) error); ok {
		r1 = rf(ctx, tag)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MailMock_FindByTag_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindByTag'
type MailMock_FindByTag_Call struct {
	*mock.Call
}

// FindByTag is a helper method to define mock.On call
//   - ctx context.Context
//   - tag types.Id
func (_e *MailMock_Expecter) FindByTag(ctx interface{}, tag interface{}) *MailMock_FindByTag_Call {
	return &MailMock_FindByTag_Call{Call: _e.mock.On("FindByTag", ctx, tag)}
}

func (_c *MailMock_FindByTag_Call) Run(run func(ctx context.Context, tag types.Id)) *MailMock_FindByTag_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(types.Id))
	})
	return _c
}

func (_c *MailMock_FindByTag_Call) Return(_a0 []entity.Mail, _a1 error) *MailMock_FindByTag_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MailMock_FindByTag_Call) RunAndReturn(run func(context.Context, types.Id) ([]entity.Mail, error)) *MailMock_FindByTag_Call {
	_c.Call.Return(run)
	return _c
}

// Patch provides a mock function with given fields: ctx, mail
func (_m *MailMock) Patch(ctx context.Context, mail *entity.Mail) error {
	ret := _m.Called(ctx, mail)

	if len(ret) == 0 {
		panic("no return value specified for Patch")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.Mail) error); ok {
		r0 = rf(ctx, mail)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MailMock_Patch_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Patch'
type MailMock_Patch_Call struct {
	*mock.Call
}

// Patch is a helper method to define mock.On call
//   - ctx context.Context
//   - mail *entity.Mail
func (_e *MailMock_Expecter) Patch(ctx interface{}, mail interface{}) *MailMock_Patch_Call {
	return &MailMock_Patch_Call{Call: _e.mock.On("Patch", ctx, mail)}
}

func (_c *MailMock_Patch_Call) Run(run func(ctx context.Context, mail *entity.Mail)) *MailMock_Patch_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*entity.Mail))
	})
	return _c
}

func (_c *MailMock_Patch_Call) Return(_a0 error) *MailMock_Patch_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MailMock_Patch_Call) RunAndReturn(run func(context.Context, *entity.Mail) error) *MailMock_Patch_Call {
	_c.Call.Return(run)
	return _c
}

// Remove provides a mock function with given fields: ctx, id
func (_m *MailMock) Remove(ctx context.Context, id types.Id) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for Remove")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, types.Id) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MailMock_Remove_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Remove'
type MailMock_Remove_Call struct {
	*mock.Call
}

// Remove is a helper method to define mock.On call
//   - ctx context.Context
//   - id types.Id
func (_e *MailMock_Expecter) Remove(ctx interface{}, id interface{}) *MailMock_Remove_Call {
	return &MailMock_Remove_Call{Call: _e.mock.On("Remove", ctx, id)}
}

func (_c *MailMock_Remove_Call) Run(run func(ctx context.Context, id types.Id)) *MailMock_Remove_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(types.Id))
	})
	return _c
}

func (_c *MailMock_Remove_Call) Return(_a0 error) *MailMock_Remove_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MailMock_Remove_Call) RunAndReturn(run func(context.Context, types.Id) error) *MailMock_Remove_Call {
	_c.Call.Return(run)
	return _c
}

// RemoveTags provides a mock function with given fields: ctx, mailId, tagIds
func (_m *MailMock) RemoveTags(ctx context.Context, mailId types.Id, tagIds []types.Id) error {
	ret := _m.Called(ctx, mailId, tagIds)

	if len(ret) == 0 {
		panic("no return value specified for RemoveTags")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, types.Id, []types.Id) error); ok {
		r0 = rf(ctx, mailId, tagIds)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MailMock_RemoveTags_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RemoveTags'
type MailMock_RemoveTags_Call struct {
	*mock.Call
}

// RemoveTags is a helper method to define mock.On call
//   - ctx context.Context
//   - mailId types.Id
//   - tagIds []types.Id
func (_e *MailMock_Expecter) RemoveTags(ctx interface{}, mailId interface{}, tagIds interface{}) *MailMock_RemoveTags_Call {
	return &MailMock_RemoveTags_Call{Call: _e.mock.On("RemoveTags", ctx, mailId, tagIds)}
}

func (_c *MailMock_RemoveTags_Call) Run(run func(ctx context.Context, mailId types.Id, tagIds []types.Id)) *MailMock_RemoveTags_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(types.Id), args[2].([]types.Id))
	})
	return _c
}

func (_c *MailMock_RemoveTags_Call) Return(_a0 error) *MailMock_RemoveTags_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MailMock_RemoveTags_Call) RunAndReturn(run func(context.Context, types.Id, []types.Id) error) *MailMock_RemoveTags_Call {
	_c.Call.Return(run)
	return _c
}

// NewMailMock creates a new instance of MailMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMailMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *MailMock {
	mock := &MailMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}