// Code generated by mockery v1.0.0. DO NOT EDIT.

package repository

import mock "github.com/stretchr/testify/mock"

// MockClusterRepo is an autogenerated mock type for the ClusterRepo type
type MockClusterRepo struct {
	mock.Mock
}

// FindByIDs provides a mock function with given fields: ids
func (_m *MockClusterRepo) FindByIDs(ids []uint64) ([]*Cluster, error) {
	ret := _m.Called(ids)

	var r0 []*Cluster
	if rf, ok := ret.Get(0).(func([]uint64) []*Cluster); ok {
		r0 = rf(ids)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*Cluster)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]uint64) error); ok {
		r1 = rf(ids)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindBySubdomains provides a mock function with given fields: subdomains
func (_m *MockClusterRepo) FindBySubdomains(subdomains []string) ([]*Cluster, error) {
	ret := _m.Called(subdomains)

	var r0 []*Cluster
	if rf, ok := ret.Get(0).(func([]string) []*Cluster); ok {
		r0 = rf(subdomains)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*Cluster)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]string) error); ok {
		r1 = rf(subdomains)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListAll provides a mock function with given fields:
func (_m *MockClusterRepo) ListAll() ([]*Cluster, error) {
	ret := _m.Called()

	var r0 []*Cluster
	if rf, ok := ret.Get(0).(func() []*Cluster); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*Cluster)
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
func (_m *MockClusterRepo) LoadByID(id uint64) (*Cluster, error) {
	ret := _m.Called(id)

	var r0 *Cluster
	if rf, ok := ret.Get(0).(func(uint64) *Cluster); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Cluster)
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