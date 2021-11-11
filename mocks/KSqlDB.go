// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"
	http "net/http"

	io "io"

	ksqldb "github.com/thmeitz/ksqldb-go"

	mock "github.com/stretchr/testify/mock"
)

// KSqlDB is an autogenerated mock type for the KSqlDB type
type KSqlDB struct {
	mock.Mock
}

// GetServerInfo provides a mock function with given fields: _a0
func (_m *KSqlDB) GetServerInfo(_a0 http.Client) (*ksqldb.KsqlServerInfo, error) {
	ret := _m.Called(_a0)

	var r0 *ksqldb.KsqlServerInfo
	if rf, ok := ret.Get(0).(func(http.Client) *ksqldb.KsqlServerInfo); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ksqldb.KsqlServerInfo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(http.Client) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Healthcheck provides a mock function with given fields: _a0, _a1
func (_m *KSqlDB) Healthcheck(_a0 http.Client, _a1 string) (*ksqldb.ServerStatusResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *ksqldb.ServerStatusResponse
	if rf, ok := ret.Get(0).(func(http.Client, string) *ksqldb.ServerStatusResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ksqldb.ServerStatusResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(http.Client, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewKsqlRequest provides a mock function with given fields: _a0, _a1
func (_m *KSqlDB) NewKsqlRequest(_a0 http.Client, _a1 io.Reader) (*http.Request, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *http.Request
	if rf, ok := ret.Get(0).(func(http.Client, io.Reader) *http.Request); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*http.Request)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(http.Client, io.Reader) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewQueryStreamRequest provides a mock function with given fields: _a0, _a1, _a2
func (_m *KSqlDB) NewQueryStreamRequest(_a0 http.Client, _a1 context.Context, _a2 io.Reader) (*http.Request, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 *http.Request
	if rf, ok := ret.Get(0).(func(http.Client, context.Context, io.Reader) *http.Request); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*http.Request)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(http.Client, context.Context, io.Reader) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}