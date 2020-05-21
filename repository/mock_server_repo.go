// Code generated by mockery v1.0.0. DO NOT EDIT.

package repository

import (
	mock "github.com/stretchr/testify/mock"
	model "github.com/rahulbharuka/dns-record-manager/model"
)

// MockServerRepo is an autogenerated mock type for the ServerRepo type
type MockServerRepo struct {
	mock.Mock
}

// AddToRotation provides a mock function with given fields: id
func (_m *MockServerRepo) AddToRotation(id uint64) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint64) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindByIPs provides a mock function with given fields: ips
func (_m *MockServerRepo) FindByIPs(ips []string) ([]*model.Server, error) {
	ret := _m.Called(ips)

	var r0 []*model.Server
	if rf, ok := ret.Get(0).(func([]string) []*model.Server); ok {
		r0 = rf(ips)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Server)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]string) error); ok {
		r1 = rf(ips)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListAll provides a mock function with given fields:
func (_m *MockServerRepo) ListAll() ([]*model.Server, error) {
	ret := _m.Called()

	var r0 []*model.Server
	if rf, ok := ret.Get(0).(func() []*model.Server); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Server)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListAllWithClusterInfo provides a mock function with given fields:
func (_m *MockServerRepo) ListAllWithClusterInfo() ([]*model.Server, error) {
	ret := _m.Called()

	var r0 []*model.Server
	if rf, ok := ret.Get(0).(func() []*model.Server); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Server)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LoadByID provides a mock function with given fields: id
func (_m *MockServerRepo) LoadByID(id uint64) (*model.Server, error) {
	ret := _m.Called(id)

	var r0 *model.Server
	if rf, ok := ret.Get(0).(func(uint64) *model.Server); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Server)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint64) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveFromRotation provides a mock function with given fields: id
func (_m *MockServerRepo) RemoveFromRotation(id uint64) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint64) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}