// Code generated by mockery v2.35.4. DO NOT EDIT.

package mock

import (
	generator "github.com/nkien0204/lets-go/internal/domain/entity/generator"
	mock "github.com/stretchr/testify/mock"
)

// GeneratorUsecase is an autogenerated mock type for the GeneratorUsecase type
type GeneratorUsecase struct {
	mock.Mock
}

// Generate provides a mock function with given fields: _a0
func (_m *GeneratorUsecase) Generate(_a0 generator.OnlineGeneratorInputEntity) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(generator.OnlineGeneratorInputEntity) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewGeneratorUsecase creates a new instance of GeneratorUsecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewGeneratorUsecase(t interface {
	mock.TestingT
	Cleanup(func())
}) *GeneratorUsecase {
	mock := &GeneratorUsecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}