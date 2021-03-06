// Code generated by mockery v2.9.4. DO NOT EDIT.

package ksqldb

import (
	"fmt"
	io "io"

	mock "github.com/stretchr/testify/mock"
)

// BodyReader is an autogenerated mock type for the BodyReader type
type BodyReader struct {
	mock.Mock
}

// Execute provides a mock function with given fields: _a0
func (_m *BodyReader) Execute(_a0 io.Reader) ([]byte, error) {
	fmt.Println("Bodyreader called")
	ret := _m.Called(_a0)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(io.Reader) []byte); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(io.Reader) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
