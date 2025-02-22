// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

type Service_Expecter struct {
	mock *mock.Mock
}

func (_m *Service) EXPECT() *Service_Expecter {
	return &Service_Expecter{mock: &_m.Mock}
}

// Process provides a mock function with given fields:
func (_m *Service) Process() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Process")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Service_Process_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Process'
type Service_Process_Call struct {
	*mock.Call
}

// Process is a helper method to define mock.On call
func (_e *Service_Expecter) Process() *Service_Process_Call {
	return &Service_Process_Call{Call: _e.mock.On("Process")}
}

func (_c *Service_Process_Call) Run(run func()) *Service_Process_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Service_Process_Call) Return(_a0 string) *Service_Process_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Service_Process_Call) RunAndReturn(run func() string) *Service_Process_Call {
	_c.Call.Return(run)
	return _c
}

// NewService creates a new instance of Service. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewService(t interface {
	mock.TestingT
	Cleanup(func())
}) *Service {
	mock := &Service{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
