// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"
	dto "nexa/services/authentication/internal/domain/dto"

	mock "github.com/stretchr/testify/mock"

	status "nexa/shared/status"

	types "nexa/shared/types"
)

// TokenMock is an autogenerated mock type for the IToken type
type TokenMock struct {
	mock.Mock
}

type TokenMock_Expecter struct {
	mock *mock.Mock
}

func (_m *TokenMock) EXPECT() *TokenMock_Expecter {
	return &TokenMock_Expecter{mock: &_m.Mock}
}

// Request provides a mock function with given fields: ctx, _a1
func (_m *TokenMock) Request(ctx context.Context, _a1 *dto.TokenCreateDTO) (dto.TokenResponseDTO, status.Object) {
	ret := _m.Called(ctx, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Request")
	}

	var r0 dto.TokenResponseDTO
	var r1 status.Object
	if rf, ok := ret.Get(0).(func(context.Context, *dto.TokenCreateDTO) (dto.TokenResponseDTO, status.Object)); ok {
		return rf(ctx, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *dto.TokenCreateDTO) dto.TokenResponseDTO); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Get(0).(dto.TokenResponseDTO)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *dto.TokenCreateDTO) status.Object); ok {
		r1 = rf(ctx, _a1)
	} else {
		r1 = ret.Get(1).(status.Object)
	}

	return r0, r1
}

// TokenMock_Request_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Request'
type TokenMock_Request_Call struct {
	*mock.Call
}

// Request is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 *dto.TokenCreateDTO
func (_e *TokenMock_Expecter) Request(ctx interface{}, _a1 interface{}) *TokenMock_Request_Call {
	return &TokenMock_Request_Call{Call: _e.mock.On("Request", ctx, _a1)}
}

func (_c *TokenMock_Request_Call) Run(run func(ctx context.Context, _a1 *dto.TokenCreateDTO)) *TokenMock_Request_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*dto.TokenCreateDTO))
	})
	return _c
}

func (_c *TokenMock_Request_Call) Return(_a0 dto.TokenResponseDTO, _a1 status.Object) *TokenMock_Request_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *TokenMock_Request_Call) RunAndReturn(run func(context.Context, *dto.TokenCreateDTO) (dto.TokenResponseDTO, status.Object)) *TokenMock_Request_Call {
	_c.Call.Return(run)
	return _c
}

// Verify provides a mock function with given fields: ctx, _a1
func (_m *TokenMock) Verify(ctx context.Context, _a1 *dto.TokenVerifyDTO) (types.Id, status.Object) {
	ret := _m.Called(ctx, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Verify")
	}

	var r0 types.Id
	var r1 status.Object
	if rf, ok := ret.Get(0).(func(context.Context, *dto.TokenVerifyDTO) (types.Id, status.Object)); ok {
		return rf(ctx, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *dto.TokenVerifyDTO) types.Id); ok {
		r0 = rf(ctx, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Id)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *dto.TokenVerifyDTO) status.Object); ok {
		r1 = rf(ctx, _a1)
	} else {
		r1 = ret.Get(1).(status.Object)
	}

	return r0, r1
}

// TokenMock_Verify_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Verify'
type TokenMock_Verify_Call struct {
	*mock.Call
}

// Verify is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 *dto.TokenVerifyDTO
func (_e *TokenMock_Expecter) Verify(ctx interface{}, _a1 interface{}) *TokenMock_Verify_Call {
	return &TokenMock_Verify_Call{Call: _e.mock.On("Verify", ctx, _a1)}
}

func (_c *TokenMock_Verify_Call) Run(run func(ctx context.Context, _a1 *dto.TokenVerifyDTO)) *TokenMock_Verify_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*dto.TokenVerifyDTO))
	})
	return _c
}

func (_c *TokenMock_Verify_Call) Return(_a0 types.Id, _a1 status.Object) *TokenMock_Verify_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *TokenMock_Verify_Call) RunAndReturn(run func(context.Context, *dto.TokenVerifyDTO) (types.Id, status.Object)) *TokenMock_Verify_Call {
	_c.Call.Return(run)
	return _c
}

// NewTokenMock creates a new instance of TokenMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTokenMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *TokenMock {
	mock := &TokenMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
